package stage

import (
	"shiphand/app/autobuild"

	"errors"
	"log"
	"os"
)

var JOB_URL string = os.Getenv("JOB_URL")

func getBaseConfig() *StageConfig {
	cfg := &StageConfig{}
	cfg.Image = "debian:stable-slim"

	return cfg
}

func getBaseStage() Stage {
	instance := Stage{}

	instance.Complete = false
	instance.Success = false
	instance.Config = getBaseConfig()

	return instance
}

func NewStage(name string, payload interface{}) (Stage, error) {
	instance := getBaseStage()
	final := make(map[string]interface{})

	for k, v := range payload.(map[interface{}]interface{}) {
		switch k := k.(type) {
		case string:
			switch k {
			case "auto-build":
				final[k] = v.(autobuild.AutoBuildConfig)
			}
		}
	}

	log.Println(final)

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
