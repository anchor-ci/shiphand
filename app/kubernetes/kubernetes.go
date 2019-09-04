package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"log"
	"os"
)

func getInClusterClientSet() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Panicf(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Panicf(err.Error())
	}

	return clientset
}

func getOutOfClusterClientSet() *kubernetes.Clientset {
	path := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", path)

	if err != nil {
		log.Panicf("ERR CREATING KUBECONFIG %v\n", err)
	}

	clientset, kubeErr := kubernetes.NewForConfig(config)

	if kubeErr != nil {
		log.Panicf("Kubernetes connecting failure")
	}

	return clientset
}

func GetKubernetesClient(inCluster bool, path string) *kubernetes.Clientset {
	env := os.Getenv("ENV")

	if env == "local" {
		return getOutOfClusterClientSet()
	} else {
		return getInClusterClientSet()
	}

	return nil
}

func GetRestConfig() *rest.Config {
	config, _ := rest.InClusterConfig()
	return config
}
