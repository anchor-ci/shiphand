package main
import (
	"errors"
	"fmt"
	"log"
    "encoding/json"
    "net/http"
    "bytes"
    "os"
)

var JOB_URL string = os.Getenv("JOB_URL")

type Stage struct {
    Metadata     JobMetadata
	Name         string
	Instructions []string
	Image        string
	Complete     bool
	Success      bool
}

func (s *Stage) Run(metadata JobMetadata) error {
    s.Metadata = metadata
	var retErr error = nil

	// Create an anchor ci managed pod
	podId := fmt.Sprintf("%s-%s", s.Metadata.Id, s.Name)
	pod, err := NewControlledPod(podId, s.Image)

	if err != nil {
		retErr = err
	}

	// Wait for pod to start before sending instructions
	pod.WaitForStart()

    updateErr := s.UpdateJobState("RUNNING")

    if updateErr != nil {
      s.Success = false
      s.Complete = true
      return updateErr
    }

	log.Printf("Controlled pod %s is ready to take commands\n", pod.Id)

	// Iterate through instructions and send to pod for execution
	for index, instruction := range s.Instructions {
		// Send series of instructions to pod
        log.Printf("Running command %s", instruction)
		history, execErr := pod.RunCommand(instruction)

		if execErr != nil || history.Failed {
			s.Success = false
		}

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Instructions)-1 {
            s.Complete = true
			s.Success = true
		}

        reportErr := s.ReportStatus(history)

        // Kill job, can't connect to job server
        if reportErr != nil {
          s.Success = false
          s.Complete = true
          log.Printf("Error updating history: %+v\n", reportErr)
          return errors.New("Stage run error")
        }
	}

    pod.CleanupPod()
	return retErr
}

func (s *Stage) UpdateJobState(state string) error {
  //Sends a PUT request to the job API that updates the current state of the job

  url := JOB_URL + "/jobs/" + s.Metadata.Id
  payload := []byte(fmt.Sprintf(`{"state":"%s"}`, state))

  req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer(payload))

  if reqErr != nil {
    return reqErr
  }

  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    return err
  }

  defer resp.Body.Close()

  if resp.StatusCode != 204 {
    return errors.New("There was an error updating the job information")
  }

  return nil
}

func (s *Stage) ReportStatus(history History) error {
  data, err := json.Marshal(history)
  payload := []byte(fmt.Sprintf(`{"history":[%s]}`, data))

  if err != nil {
    return errors.New("Couldn't connect to jobs API")
  }

  log.Printf("Updating history: %+v", history)

  url := JOB_URL + "/histories/" + s.Metadata.HistoryId
  req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer(payload))

  if reqErr != nil {
    return reqErr
  }

  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    return err
  }

  defer resp.Body.Close()

  if resp.StatusCode != 204 {
    return errors.New("There was an error updating the job information")
  }

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
