package main

import (
  "log"
  "fmt"
  "errors"
)

type Stage struct {
  Name string
  Instructions []string
  Image string
}

func (s *Stage) Run(metadata JobMetadata) error {
  // Create an anchor ci managed pod
  podId := fmt.Sprintf("%s-%s", metadata.Id, s.Name)
  pod, err := NewControlledPod(podId, s.Image)

  log.Printf("Created controlled pod: %s\n", pod.Id)

  // Wait for pod to start before sending instructions
  pod.WaitForStart()

  log.Printf("Controlled pod %s is ready to take commands\n", pod.Id)

  // Iterate through instructions and send to pod for execution
  for _, instruction := range s.Instructions {
    // Send series of instructions to pod
    _, execErr := pod.RunCommand(instruction)

    if execErr != nil {
      log.Printf("Exec error: %+v\n", execErr)
    }
  }

  return err
}

func getBaseStage() Stage {
  instance := Stage{}

  instance.Image = "debian:stable-slim"

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
