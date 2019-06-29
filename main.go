package main

import (
	"fmt"
	"os"
	//"path/filepath"
	"github.com/go-redis/redis"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
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

				fmt.Printf("Starting job: %s\n", job)
				go startJob(job, jid)

				_, delErr := client.Del(job).Result()

				if delErr != nil {
					fmt.Printf("Error removing job %s: %v\n", job, err)
				}
			}
		}
	}
}

func startJob(key string, payload string) {
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
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-job",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []apiv1.Container{
						{
							Name:  "demo",
							Image: "myimage",
						},
					},
				},
			},
		},
	}

	result, jobErr := jobClient.Create(job)
	if jobErr == nil {
		fmt.Printf("Job started: %v\n", result)
	} else {
		fmt.Printf("Error starting job: %v", jobErr)
	}
}
