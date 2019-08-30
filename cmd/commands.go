package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func GetCommands() []*cobra.Command {
  return []*cobra.Command{
    &cobra.Command{
      Use: "test",
      Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello :)")
      },
    },
  }
}
