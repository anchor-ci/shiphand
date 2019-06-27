package main

import (
  "fmt"
  "os"
  //"path/filepath"
  "strconv"
  "github.com/go-redis/redis"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/kubernetes"
  apiv1 "k8s.io/api/core/v1"
)


const REDIS_URL string = "0.0.0.0"
const REDIS_PORT int = 6379
const JOB_KEY string = "job:v1:*"

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
  client := redis.NewClient(&redis.Options{
    Addr: REDIS_URL + ":" + strconv.Itoa(REDIS_PORT),
  })

  _, err := client.Ping().Result()

  if err != nil {
    panic("Error connecting to job's database")
  }

  for /* ever */ {
    // TODO: Make this activated on pubsub, this is going to hammer the server
    jobs, err := client.Keys(JOB_KEY).Result()

    if err == nil {
      for _, job := range jobs {
        jid, err := client.Get(job).Result()

        // If errors grabbing key, just skip iteration
        if err != nil {
          fmt.Printf("Unable to grab job payload for %s\n", job)
          continue
        }

        go startJob(job, jid)

        _, delErr := client.Del(job).Result()

        if delErr != nil {
          fmt.Printf("Error removing job %s: %v\n", job, err)
        } else {
          fmt.Printf("Finished starting job %s\n", job)
        }
      }
    }
  }
}

func startJob(key string, payload string) {
  fmt.Printf("Starting job: %s\n", key)
  // TODO: Move kube stuff to method
  kubeconfig := "./config"

  config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
  if err != nil {
    fmt.Printf("ERR CREATING KUBECONFIG %v\n", err)
  }

  clientset, kubeErr := kubernetes.NewForConfig(config)
  if kubeErr != nil {
    fmt.Printf("Kubernetes connecting failure")
  }

  jobClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
  fmt.Println(jobClient)
}
