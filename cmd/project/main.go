package project

import (
	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use: "project",
}

func init() {
	ProjectCmd.AddCommand(addProjectCmd)
	ProjectCmd.AddCommand(removeProjectCmd)
	ProjectCmd.AddCommand(setProjectCmd)
}
