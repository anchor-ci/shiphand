package main

import (
	"errors"
	"fmt"
	"log"
    "encoding/json"
    "net/http"
)

type Stage struct {
	Name         string
	Instructions []string
	Image        string
	Complete     bool
	Success      bool
}

func (s *Stage) Run(metadata JobMetadata) error {
	var retErr error = nil

	// Create an anchor ci managed pod
	podId := fmt.Sprintf("%s-%s", metadata.Id, s.Name)
	pod, err := NewControlledPod(podId, s.Image)

	if err != nil {
		retErr = err
	}

	log.Printf("Created controlled pod: %s\n", pod.Id)

	// Wait for pod to start before sending instructions
	pod.WaitForStart()

	log.Printf("Controlled pod %s is ready to take commands\n", pod.Id)

	// Iterate through instructions and send to pod for execution
	for index, instruction := range s.Instructions {
		// Send series of instructions to pod
		history, execErr := pod.RunCommand(instruction)

		if execErr != nil || history.Failed {
			s.Success = false
			break
		}

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Instructions)-1 {
			s.Success = true
		}

        reportErr := s.ReportStatus(history)

        // Kill job, can't connect to job server
        if reportErr != nil {
          s.Success = false
          s.Complete = true
        }
	}

	s.Complete = true

	return retErr
}

func (s *Stage) ReportStatus(history History) error {
  data, err := json.Marshal(history)
  if err != nil {
    return errors.New("Couldn't connect to jobs API")
  }

  resp, httpErr := http.Put("http://172.18.0.5:8080/")
  return nil
}

func getBaseStage() Stage {
	instance := Stage{}

	instance.Image = "debian:stable-slim"
	instance.Complete = false
	instance.Success = false

	return instance
}

func NewStage(name string, payload interface{}) (Stage, error) {
	instance := getBaseStage()
	transformedVal := payload.(map[string]interface{})

	if script, ok := transformedVal["script"].([]interface{}); ok {
		for _, v := range script {
			if instruction, ok := v.(string); ok {
				instance.Instructions = append(instance.Instructions, instruction)
			} else {
				return instance, errors.New("Couldn't get instructions")
			}
		}
	} else {
		return instance, errors.New("Couldn't get instructions")
	}

	instance.Name = name

	return instance, nil
}
