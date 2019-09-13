package app

import (
	"github.com/buger/jsonparser"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"

	"encoding/json"
	"log"
	"os"
	"shiphand/app/job"
)

func AppMain(c *cli.Context) {
	// Check for debug mode
	if c.String("run") != "" {
		DebugMode(c.String("run"))
		os.Exit(0)
	}

	client := redis.NewClient(&redis.Options{
		Addr: c.String("redis-host") + ":" + c.String("redis-port"),
	})

	_, err := client.Ping().Result()

	if err != nil {
		log.Fatal("Error connecting to Redis")
	}

	for /* ever */ {
		// TODO: Make this activated on pubsub, this is going to hammer the server
		jobs, err := client.Keys(c.String("key")).Result()

		if err == nil {
			for _, job := range jobs {
				jid, err := client.Get(job).Result()

				// If errors grabbing key, just skip iteration
				if err != nil {
					log.Printf("Unable to grab job payload for %s\n", job)
					continue
				}

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
	metadata := job.JobMetadata{}

	if err := json.Unmarshal([]byte(payload), &metadata); err != nil {
		log.Printf("Couldn't unmarshal payload into metadata %+v\n", err)
	}

	if err := json.Unmarshal([]byte(payload), &f); err != nil {
		log.Printf("Couldn't unmarshal payload %+v\n", err)
	}

	if historyId, jsonErr := jsonparser.GetString([]byte(payload), "history", "id"); jsonErr == nil {
		metadata.HistoryId = historyId
	} else {
		log.Printf("Couldn't get history ID from: %s\n", payload)
	}

    log.Println(payload)
}
