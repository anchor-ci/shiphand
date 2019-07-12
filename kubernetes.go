package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func getKubernetesClient(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("ERR CREATING KUBECONFIG %v\n", err)
	}

	clientset, kubeErr := kubernetes.NewForConfig(config)
	if kubeErr != nil {
		log.Printf("Kubernetes connecting failure")
	}

	return clientset
}
