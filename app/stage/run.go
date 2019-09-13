package stage

import (
	"shiphand/app/manager"

	"errors"
	"log"
)

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
	for index, instruction := range s.Script {
		// Send series of instructions to pod
		log.Printf("Running command %s", instruction)
		report, execErr := pod.RunCommand(instruction)

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Script)-1 {
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
