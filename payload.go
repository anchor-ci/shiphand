package main

import (
  //"log"
  "errors"
)

type Stage struct {
  Name string
  Instructions []string
}

type Job struct {
  Name string
  Stages []Stage
}

type Payload struct {
  Jobs []Job
}

func NewStage(name string, payload interface{}) (Stage, error) {
  instance := Stage{}
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
