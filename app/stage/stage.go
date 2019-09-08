package stage

import (
	"shiphand/app/autobuild"
)

type Stage struct {
	AutoBuild autobuild.AutoBuildConfig `json:"auto-build"`
	Script    []string                  `json:"script"`
	Image     string                    `json:"image"`
	Name      string                    `json:"name"`
	Complete  bool
	Success   bool
}
