package main

import (
  "errors"
)

type Stage struct {
  Name string
  Instructions []string
  Image string
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
