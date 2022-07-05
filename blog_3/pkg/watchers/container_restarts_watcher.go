package watchers

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

const (
	DEFAULT_RESYNC = time.Hour * 24
)

type ContainerRestartEvent struct {
	Namespace     string
	PodName       string
	ContainerName string
	RestartCount  int32
	ExitCode      int32
	Reason        string
	Message       string
}

type PodLoggingController struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

// NewPodLoggingController creates a PodLoggingController
func NewPodLoggingController(client *kubernetes.Clientset) *PodLoggingController {
	factory := informers.NewSharedInformerFactory(client, DEFAULT_RESYNC)

	podInformer := factory.Core().V1().Pods()

	c := &PodLoggingController{
		informerFactory: factory,
		podInformer:     podInformer,
	}
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: podUpdate,
		},
	)

	return c
}

func (c *PodLoggingController) Run(stopCh chan struct{}) error {
	c.informerFactory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}

	return nil
}

func podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)

	// on informer sync we can get an update with the same pod ResourceVersion
	if oldPod.ResourceVersion == newPod.ResourceVersion {
		return
	}

	for _, oldContainer := range oldPod.Status.ContainerStatuses {
		for _, newContainer := range newPod.Status.ContainerStatuses {
			if oldContainer.Name != newContainer.Name {
				continue
			}

			if oldContainer.RestartCount != newContainer.RestartCount && newContainer.LastTerminationState.Terminated != nil {
				restartEvent := &ContainerRestartEvent{
					Namespace:     newPod.Namespace,
					PodName:       newPod.Name,
					ContainerName: newContainer.Name,
					RestartCount:  newContainer.RestartCount,
					ExitCode:      newContainer.LastTerminationState.Terminated.ExitCode,
					Reason:        newContainer.LastTerminationState.Terminated.Reason,
					Message:       newContainer.LastTerminationState.Terminated.Message,
				}

				klog.Infof("container_restart_event: %+v", restartEvent)
			}
		}
	}
}
