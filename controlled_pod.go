package main

import (
   apiv1 "k8s.io/api/core/v1"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControlledPod struct {
  Pod *apiv1.Pod
}

func NewControlledPod(name string, image string) (ControlledPod, error) {
  instance := ControlledPod{}

  // TODO: Uncouple the kube client stuff from here
  clientset := getKubernetesClient("./config")

  req := &apiv1.Pod{
      TypeMeta: metav1.TypeMeta{
          Kind:       "Pod",
          APIVersion: "v1",
      },
      ObjectMeta: metav1.ObjectMeta{
          Name: name,
      },
      Spec: apiv1.PodSpec{
          Containers: []apiv1.Container{
              {
                  Name:  "job",
                  Image: image,
              },
          },
      },
  }

  pod, err := clientset.CoreV1().Pods("default").Create(req)

  instance.Pod = pod

  return instance, err
}
