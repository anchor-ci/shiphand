package stage

import (
    "shiphand/app/autobuild"
    "shiphand/app/manager"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var JOB_URL string = os.Getenv("JOB_URL")

type Stage struct {
	Image        string
	Instructions []string
	Name         string
	Complete     bool
	Success      bool
}

type StageConfig struct {
	Clone     bool
	AutoBuild autobuild.AutoBuildConfig
	Script    []string
	Image     string
}

func (s *Stage) Run(name string,
                    jobId string,
                    historyId string) error {

	var retErr error = nil

	// Create an anchor ci managed pod
	pod, err := manager.NewControlledPod(name, s.Image)

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
	for index, instruction := range s.Instructions {
		// Send series of instructions to pod
		log.Printf("Running command %s", instruction)
		report, execErr := pod.RunCommand(instruction)

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Instructions)-1 {
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

func getBaseStage() Stage {
	instance := Stage{}

	instance.Image = "debian:stable-slim"
	instance.Complete = false
	instance.Success = false

	return instance
}

func NewStage(name string, payload interface{}) (Stage, error) {
	instance := getBaseStage()
	transformedVal := payload.(map[interface{}]interface{})

	if script, ok := transformedVal["script"].([]interface{}); ok {
		instructions, err := getInstructions(script)

		if err != nil {
			return instance, err
		}

		instance.Instructions = instructions
	} else {
		return instance, errors.New("Couldn't get instructions")
	}

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
