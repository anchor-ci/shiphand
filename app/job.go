package app

import (
	"errors"
)

type Job struct {
	Name   string
	Stages []Stage
}

func (j *Job) Run(metadata JobMetadata) (bool, error) {
	for _, stage := range j.Stages {
		err := stage.Run(metadata)
		if err != nil {
			return false, err
		}

		pass, interErr := j.InterStage(&stage)

		if interErr != nil {
			return false, interErr
		}

		if !pass {
			return pass, nil
		}
	}

	return true, nil
}

// Ran in between each stage, verifies the stage ran ok
// and sees if it needs to report any failure. Returns
// a bool indicating if the job should continue, error
// indicating an error within the function
func (j *Job) InterStage(s *Stage) (bool, error) {
	if !s.Complete {
		return false, errors.New("Stage %s didn't complete")
	}

	if s.Complete && !s.Success {
		return false, nil
	}

	return s.Complete && s.Success, nil
}

func NewJob(name string, payload interface{}) (Job, error) {
	instance := Job{}
	transformedVal := payload.(map[string]interface{})

	for k, v := range transformedVal {
		stage, err := NewStage(k, v)

		if err != nil {
			return instance, err
		}

		instance.Stages = append(instance.Stages, stage)
	}

	instance.Name = name

	return instance, nil
}
