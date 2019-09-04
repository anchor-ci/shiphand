package app

import (
    shiphand_payload "shiphand/app/payload"

	"io/ioutil"
	"log"
)

func DebugMode(file string) {
	log.Printf("Running debug mode for: %s\n", file)

	contents, fileErr := ioutil.ReadFile(file)

	if fileErr != nil {
		log.Fatalf("Couldn't create payload from:\n %s", file)
	}

    payload, err := shiphand_payload.PayloadFromJson(string(contents))

	log.Printf("Payload:\n %+v", payload)

	if err != nil {
		log.Fatalf("Couldn't create payload from:\n %s\nCause:%v\n", file, err)
	}

    Run(payload)
}

func Run(payload *shiphand_payload.Payload) error {
	for _, job := range payload.Jobs {
        err := job.DebugRun()

		if err != nil {
			return err
		}
	}

	return nil
}
