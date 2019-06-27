package main

import (
  "fmt"
  "strconv"
  "github.com/go-redis/redis"
)

const REDIS_URL string = "0.0.0.0"
const REDIS_PORT int = 6379
const JOB_KEY string = "job:v1:*"

func main() {
  client := redis.NewClient(&redis.Options{
    Addr: REDIS_URL + ":" + strconv.Itoa(REDIS_PORT),
  })

  _, err := client.Ping().Result()

  if err != nil {
    panic("Error connecting to job's database")
  }

  for /* ever */ {
    jobs, err := client.Keys(JOB_KEY).Result()

    if err != nil {
      fmt.Println("No jobs")
    } else {
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
          fmt.Printf("Finished processing job %s\n", job)
        }
      }
    }
  }
}

func startJob(key string, payload string) {
  fmt.Println(key)
  fmt.Println(payload)
}
