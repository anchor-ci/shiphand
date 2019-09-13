package app

import (
	"gopkg.in/yaml.v2"
	"shiphand/app/job"

	"errors"
)

type Definition struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Metadata job.JobMetadata
	Jobs     []job.Job `json:"jobs"`
}

func PayloadFromJson(file string) (*Payload, error) {
	payload := &Payload{}

	jobs, err := CreateJobs(file)

	if err != nil {
		return payload, err
	}

	payload.Jobs = jobs

	return payload, nil
}

func NewPayload(payload map[string]interface{}) (Payload, error) {
	instance := Payload{}

	// Constructs the instructions for the job
	if val, ok := payload["instructions"]; ok {
		transformedVal := val.(map[string]interface{})
		jobs, err := jobsFromStrInter(transformedVal)

		if err != nil {
			return instance, err
		}

		instance.Jobs = jobs
	} else {
		return instance, errors.New("No instructions defined")
	}

	return instance, nil
}

func jobsFromStrInter(payload map[string]interface{}) ([]job.Job, error) {
	jobs := []job.Job{}

	for k, v := range payload {
		job, err := job.NewJob(k, v)

		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func CreateJobs(file string) ([]job.Job, error) {
	var vals map[string]interface{}
	jobs := []job.Job{}
	err := yaml.Unmarshal([]byte(file), &vals)

	if err != nil {
		return jobs, err
	}

	for name, currentJob := range vals {
		jerb, err := job.NewJob(name, currentJob)

		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, jerb)
	}

	return jobs, nil
}

func (p *Payload) Run() error {
	for _, job := range p.Jobs {
		_, err := job.Run(p.Metadata)
		if err != nil {
			return err
		}
	}

	return nil
}
