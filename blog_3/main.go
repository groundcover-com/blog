package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"groundcover.com/pkg/k8s_utils"
	"groundcover.com/pkg/watchers"
	watcher "groundcover.com/pkg/watchers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path")
}

func main() {
	flag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	client, err := k8s_utils.NewKubeClient(kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	eventCh := make(chan *watchers.ContainerRestartEvent)
	stopCh := make(chan struct{})

	listen(client, eventCh, stopCh)
	go processRestartEvents(eventCh)

	waitForInterrupt()
	close(stopCh)
	close(eventCh)
}

func waitForInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch
	klog.Errorf("interrupt invoked. signal %s, closing...", sig.String())
}

func listen(client *kubernetes.Clientset, eventCh chan *watcher.ContainerRestartEvent, stopCh chan struct{}) {
	watcher := watcher.NewPodWatcher(client, eventCh)

	err := watcher.Run(stopCh)
	if err != nil {
		klog.Fatal(err)
	}
}

func processRestartEvents(eventCh chan *watcher.ContainerRestartEvent) {
	for event := range eventCh {
		klog.Infof("event: %+v", event)
	}
}
