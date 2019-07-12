package main

import (
	"errors"
	"fmt"
	"log"
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
		_, execErr := pod.RunCommand(instruction)

		if execErr != nil {
			log.Printf("Exec error: %+v\n", execErr)
            s.Success = false
            break
		}

        if index == len(s.Instructions) - 1 {
            s.Success = true
        }
	}

    s.Complete = true

	return retErr
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
