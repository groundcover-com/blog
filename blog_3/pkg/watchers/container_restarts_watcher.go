package watchers

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

type PodWatcher struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
	eventCh         chan *ContainerRestartEvent
}

func NewPodWatcher(client *kubernetes.Clientset, eventCh chan *ContainerRestartEvent) *PodWatcher {
	factory := informers.NewSharedInformerFactory(client, DEFAULT_RESYNC)

	podInformer := factory.Core().V1().Pods()

	c := &PodWatcher{
		informerFactory: factory,
		podInformer:     podInformer,
		eventCh:         eventCh,
	}
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: c.podUpdate,
		},
	)

	return c
}

func (c *PodWatcher) Run(stopCh chan struct{}) error {
	c.informerFactory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}

	return nil
}

func (c *PodWatcher) podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)

	// on informer sync we can get an update with the same pod ResourceVersion, we ignore these
	if oldPod.ResourceVersion == newPod.ResourceVersion {
		return
	}

	for _, oldContainer := range oldPod.Status.ContainerStatuses {
		for _, newContainer := range newPod.Status.ContainerStatuses {
			if oldContainer.Name != newContainer.Name {
				continue
			}

			// this means the container has restarted, and has termination details
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

				c.eventCh <- restartEvent
			}
		}
	}
}
