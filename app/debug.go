package app

import (
	"io/ioutil"
	"log"
)

func DebugMode(file string) {
	log.Printf("Running debug mode for: %s\n", file)

	contents, fileErr := ioutil.ReadFile(file)

	if fileErr != nil {
		log.Fatalf("Couldn't create payload from:\n %s", file)
	}

	payload, err := PayloadFromJson(string(contents))

	log.Printf("Payload:\n %+v", payload)

	if err != nil {
		log.Fatalf("Couldn't create payload from:\n %s\nCause:%v\n", file, err)
	}

	Run(&payload)
}

func Run(payload *Payload) error {
	for _, job := range payload.Jobs {
		_, err := job.DebugRun()

		if err != nil {
			return err
		}
	}

	return nil
}
