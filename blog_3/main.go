package main

import (
	"flag"

	"groundcover.com/pkg/k8s_utils"
	watcher "groundcover.com/pkg/watchers"
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

	controller := watcher.NewPodLoggingController(client)
	stop := make(chan struct{})
	defer close(stop)

	err = controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}

	select {}
}
