package stage

import (
  "log"
)

func (s *Stage) DebugRun(name string) error {
  log.Println("Debug runnin!")
  return nil
}
