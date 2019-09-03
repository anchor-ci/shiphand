package app

// Struct representing a job request from redis
type JobRequest struct {
	JobId        string `json:"id"`
	State        string `json:"state"`
	RepositoryId string `json:"repository_id"`
}
