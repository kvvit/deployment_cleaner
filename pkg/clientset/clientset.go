package clientset

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetClientset() *kubernetes.Clientset {
	configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	var config *rest.Config
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("failed to build local config: %s\n", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			log.Fatalf("failed to build in-cluster config: %s\n", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create endpoints watcher: %s\n", err)
	}
	return clientset
}
