package job

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"shiphand/app/stage"
	"time"
)

type Job struct {
	Name   string        `json:"name"`
	Stages []stage.Stage `json:"stages"`
}

func (j *Job) Run(metadata JobMetadata) (bool, error) {
	for _, stage := range j.Stages {
		// TODO: Remove this name v
		err := stage.Run("test-job", metadata.Id, metadata.HistoryId)

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
func (j *Job) InterStage(s *stage.Stage) (bool, error) {
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
	transformedVal := payload.(map[interface{}]interface{})

	for k, v := range transformedVal {
		stage, err := stage.NewStage(k.(string), v)

		if err != nil {
			return instance, err
		}

		instance.Stages = append(instance.Stages, stage)
	}

	instance.Name = name

	return instance, nil
}

func (j *Job) DebugRun() error {
	for _, currentStage := range j.Stages {
		log.Printf("> Running stage: %s\n", currentStage.Name)

		// Create a new time based seed
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Attach a random # to prevent pod naming collisions
		name := fmt.Sprintf("debug-run-%d", random.Intn(10000000))
		err := currentStage.DebugRun(name)

		if err != nil {
			log.Printf("> Stage [%s] failed with reason: %s", currentStage.Name, err)
			return err
		}
	}

	return nil
}
