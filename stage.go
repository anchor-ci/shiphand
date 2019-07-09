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
  pod, err := NewControlledPod(fmt.Sprintf("%s-%s", metadata.Id, s.Name), s.Image)

  log.Printf("Created controlled pod: %+v\n", pod)

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
