package k8s_utils

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubeconfigPath(kubeconfigFlag string) string {
	if kubeconfigFlag != "" {
		return kubeconfigFlag
	}
	return filepath.Join(homedir.HomeDir(), ".kube", "config")
}

func NewKubeClient(kubeconfigFlag string) (*kubernetes.Clientset, error) {
	kubeconfigPath := getKubeconfigPath(kubeconfigFlag)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
