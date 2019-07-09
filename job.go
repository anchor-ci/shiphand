package main

type Job struct {
  Name string
  Stages []Stage
}

func (j *Job) Run(metadata JobMetadata) error {
  for _, stage := range j.Stages {
    err := stage.Run(metadata)
    if err != nil {
      return err
    }
  }

  return nil
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
