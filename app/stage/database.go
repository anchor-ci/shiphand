package stage

import (
  "log"
)

func GetStageDatabase() map[string]func(int) {
  log.Println("Database being grabbed.")

  db := make(map[string]func(int), 1)
  return db
}
