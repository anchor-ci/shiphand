package stage

import (
	"errors"
	"log"
	"shiphand/app/manager"
)

func (s *Stage) DebugRun(name string) error {
	log.Printf(">> Pod [%s] being created..\n", name)

	// Create an anchor ci managed pod
	pod, err := manager.NewControlledPod(name, s.Image)

	log.Printf(">> Pod [%s] created!\n", name)

	if err != nil {
		return err
	}

	err = pod.WaitForStart()

	if err != nil {
		return err
	}

	// Iterate through instructions and send to pod for execution
	for index, instruction := range s.Script {

		log.Printf(">> Running instruction {%s}", instruction)

		report, execErr := pod.RunCommand(instruction)

		// Means we hit the end of all instructions, can be marked as success
		if index == len(s.Script)-1 {
			s.Complete = true
		}

		if s.Complete {
			s.Success = execErr == nil && !report.Failed
		}
	}

	if s.Success {
		log.Printf(">> Stage [%s] passed\n", name)
	} else {
		return errors.New("Stage failed")
	}

	pod.CleanupPod()

	log.Printf(">> Cleaning up pod [%s]\n", name)

	return nil
}
