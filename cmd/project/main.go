package project

import (
	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use: "project",
}

func init() {
	ProjectCmd.AddCommand(NewAddProjectCmd())
	ProjectCmd.AddCommand(NewRemoveProjectCmd())
	ProjectCmd.AddCommand(NewSetProjectCmd())
}
