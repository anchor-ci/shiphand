package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/rest"

	"log"
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

func getOutOfClusterClientSet(path string) *kubernetes.Clientset {
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

func getKubernetesClient(inCluster bool, path string) *kubernetes.Clientset {
  if !inCluster {
    return getOutOfClusterClientSet(path)
  } else {
    return getInClusterClientSet()
  }

  return nil
}

func getRestConfig() *rest.Config {
  config, _ := rest.InClusterConfig()
  return config
}
