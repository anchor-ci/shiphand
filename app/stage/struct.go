package stage

import (
  "shiphand/app/autobuild"
)

type Stage struct {
	Name         string
	Complete     bool
	Success      bool
    Config       *StageConfig
}

type StageConfig struct {
	Clone     bool
	AutoBuild autobuild.AutoBuildConfig
	Script    []string
	Image     string
}

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
