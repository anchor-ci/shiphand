package stage

import (
	"shiphand/app/manager"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

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
