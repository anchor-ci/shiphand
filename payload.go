package main

import (
  "log"
  "errors"
   batchv1 "k8s.io/api/batch/v1"
   apiv1 "k8s.io/api/core/v1"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Payload struct {
  Metadata JobMetadata
  Jobs []Job
}

// TODO: Figure out a better implementation of this
//func getValue(name string, payload interface{}) (string, error) {
//  if val, ok := payload.(map[string]interface{}); ok {
//    if extracted, ok := val.(string); ok {
//      return extracted, nil
//    } else {
//      return "", errors.New("Couldn't extract value")
//    }
//  } else {
//    return "", errors.New("Couldn't extract value")
//  }
//}

func NewPayload(payload map[string]interface{}) (Payload, error) {
  instance := Payload{}

  if val, ok := payload["instructions"]; ok {
    transformedVal := val.(map[string]interface{})

    for k, v := range transformedVal {
      job, err := NewJob(k, v)

      if err != nil {
        return instance, err
      }

      instance.Jobs = append(instance.Jobs, job)
    }
  } else {
    return instance, errors.New("No instructions defined")
  }

  return instance, nil
}

func (p *Payload) Run() error {
  // TODO: Uncouple the kube client stuff from here
  clientset := getKubernetesClient("./config")

  jobClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
  job := &batchv1.Job{
      ObjectMeta: metav1.ObjectMeta{
          Name:      "job-" + p.Metadata.Id,
          // TODO: Namespace should probably be keyed on the owner
          Namespace: "default",
      },
      Spec: batchv1.JobSpec{
          Template: apiv1.PodTemplateSpec{
              Spec: apiv1.PodSpec{
                  RestartPolicy: "Never",
                  Containers: []apiv1.Container{
                      {
                          Name:  "job",
                          Image: "debian:stable-slim",
                          Command: []string{},
                      },
                  },
              },
          },
      },
  }

  _, jobErr := jobClient.Create(job)
  if jobErr == nil {
  } else {
      log.Printf("Error starting job: %v", jobErr)
  }
  return nil
}
