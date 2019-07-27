package main

import (
	"errors"
)

type Payload struct {
	Metadata JobMetadata
	Jobs     []Job
}

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
	for _, job := range p.Jobs {
		_, err := job.Run(p.Metadata)
		if err != nil {
			return err
		}
	}

	return nil
}
