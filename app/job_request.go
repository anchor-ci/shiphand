package app

import (
  "shiphand/app/repository"
)

// Struct representing a job request from redis
type JobRequest struct {
	JobId        string `json:"id"`
	State        string `json:"state"`
    Repository   repository.Repository `json:"repository"`
}
