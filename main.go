package main

import (
    "log"
    "strings"
    "encoding/json"
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

// Nested struct representing an instruction set for a job
type InstructionSet struct {
  Id string `json:"id"`
  Instructions []string `json:"instructions"`
  JobId string `json:"job_id"`
}

// Struct representing a job request from redis
type JobRequest struct {
  JobId string `json:"id"`
  State string `json:"state"`
  RepositoryId string `json:"repository_id"`
  Instructions InstructionSet `json:"instruction_set"`
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
					log.Printf("Unable to grab job payload for %s\n", job)
					continue
				}

				log.Printf("Starting job: %s\n", job)
				go startJob(job, jid)

				_, delErr := client.Del(job).Result()

				if delErr != nil {
					log.Printf("Error removing job %s: %v\n", job, err)
				}
			}
		}
	}
}

func startJob(key string, payload string) {
	// TODO: Move kube stuff to method
	kubeconfig := "./config"

    // Services that queue up the jobs are written in Python and can possibly send a ' for a string instead of ", for now the fix lives here, in the future (TODO) it needs to move to the Python services.
    r := strings.NewReplacer("'", "\"")
    instructions := JobRequest{}
    bytes := []byte(r.Replace(payload))

    if err := json.Unmarshal(bytes, &instructions); err != nil {
      log.Printf("Couldn't unmarshal json %+v\n", err)
    }

    log.Printf("Result is: %+v\n", instructions)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("ERR CREATING KUBECONFIG %v\n", err)
	}

	clientset, kubeErr := kubernetes.NewForConfig(config)
	if kubeErr != nil {
		log.Printf("Kubernetes connecting failure")
	}

	jobClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-" + instructions.JobId,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []apiv1.Container{
						{
							Name:  "job",
                            Image: "debian:stable-slim",
                            Command: instructions.Instructions.Instructions,
						},
					},
				},
			},
		},
	}

	result, jobErr := jobClient.Create(job)
	if jobErr == nil {
		log.Printf("Job started: %v\n", result)
	} else {
		log.Printf("Error starting job: %v", jobErr)
	}
}
