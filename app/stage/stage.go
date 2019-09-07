package stage

import (
    "shiphand/app/manager"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var STAGE_DATABASE map[string]func(int) = GetStageDatabase()
var JOB_URL string = os.Getenv("JOB_URL")

func (s *Stage) Run(name string,
                    jobId string,
                    historyId string) error {

	var retErr error = nil

	// Create an anchor ci managed pod
	pod, err := manager.NewControlledPod(name, s.Config.Image)

	if err != nil {
		retErr = err
	}

	// Wait for pod to start before sending instructions
	pod.WaitForStart()

	updateErr := s.UpdateJobState(jobId, "RUNNING")

	if updateErr != nil {
		s.Success = false
		s.Complete = true
		return updateErr
	}

	log.Printf("Controlled pod %s is ready to take commands\n", pod.Id)

	// Iterate through instructions and send to pod for execution
	for index, instruction := range s.Config.Script {
		// Send series of instructions to pod
		log.Printf("Running command %s", instruction)
		report, execErr := pod.RunCommand(instruction)

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Config.Script)-1 {
			s.Complete = true
			s.Success = true
		}

		if execErr != nil || report.Failed {
			s.Success = false
		}

		reportErr := s.ReportStatus(historyId, report)

		// Kill job, can't connect to job server
		if reportErr != nil {
			s.Success = false
			s.Complete = true
			log.Printf("Error updating history: %+v\n", reportErr)
			return errors.New("Stage run error")
		}
	}

	if s.Success {
		s.UpdateJobState(jobId, "SUCCESS")
	} else {
		s.UpdateJobState(jobId, "FAILED")
	}

	pod.CleanupPod()
	return retErr
}

// TODO: Decouple metadata from stage, move it up to job level.
// Maybe use goroutine to communicate back to API?
func (s *Stage) UpdateJobState(id string, state string) error {
	//Sends a PUT request to the job API that updates the current state of the job

	url := JOB_URL + "/jobs/" + id
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

func (s *Stage) ReportStatus(id string, report *manager.Report) error {
	data, err := json.Marshal(report)
	payload := []byte(fmt.Sprintf(`{"history":[%s]}`, data))

	if err != nil {
		return errors.New("Couldn't connect to jobs API")
	}

	log.Printf("Updating report: %+v", report)

	url := JOB_URL + "/histories/" + id
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

func NewStage(name string, payload interface{}) (Stage, error) {
	instance := getBaseStage()

    for k, v := range payload.(map[interface{}]interface{}) {
      switch k := k.(type) {
        case string:
          switch k {
          case "auto-build":
            log.Println("Auto-buildin")
          }
          log.Printf("Key: %+v, value: %+v\n", k, v)
      }
    }

    panic("Bad!")

//	if script, ok := transformedVal["script"].([]interface{}); ok {
//		instructions, err := getInstructions(script)
//
//		if err != nil {
//			return instance, err
//		}
//
//		instance.Config.Script = instructions
//	} else {
//		return instance, errors.New("Couldn't get instructions")
//	}
//
	instance.Name = name

	return instance, nil
}

func getInstructions(instructions []interface{}) ([]string, error) {
	instances := []string{}

	for _, v := range instructions {
		if instruction, ok := v.(string); ok {
			instances = append(instances, instruction)
		} else {
			return instances, errors.New("Couldn't get instructions")
		}
	}

	return instances, nil
}
