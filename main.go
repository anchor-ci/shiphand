package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"strconv"
)

const REDIS_URL string = "0.0.0.0"
const REDIS_PORT int = 6379
const JOB_KEY string = "job:v1:*"

// Struct representing a job request from redis
type JobRequest struct {
	JobId        string `json:"id"`
	State        string `json:"state"`
	RepositoryId string `json:"repository_id"`
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
	var f interface{}
	metadata := JobMetadata{}

	if err := json.Unmarshal([]byte(payload), &metadata); err != nil {
		log.Printf("Couldn't unmarshal payload into metadata %+v\n", err)
	}

	if err := json.Unmarshal([]byte(payload), &f); err != nil {
		log.Printf("Couldn't unmarshal payload %+v\n", err)
	}

	instructionSet := f.(map[string]interface{})
	tSet := instructionSet["instruction_set"].(map[string]interface{})
	finalPayload, payloadErr := NewPayload(tSet)

	finalPayload.Metadata = metadata

	if payloadErr != nil {
		log.Printf("Failed to create payload for %s\n", key)
		return
	}

	log.Printf("Created payload: %d, starting.\n", finalPayload.Metadata.Id)

	err := finalPayload.Run()

	if err != nil {
		log.Printf("Error running job: %+v\n", err)
	}
}
