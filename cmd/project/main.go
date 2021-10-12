package project

import (
	"bufio"
	"github.com/spf13/cobra"
	"os"
)

var ProjectCmd = &cobra.Command{
	Use: "project",
}

func init() {
	ProjectCmd.AddCommand(NewAddProjectCmd())
	ProjectCmd.AddCommand(NewRemoveProjectCmd())
	ProjectCmd.AddCommand(NewSetProjectCmd(bufio.NewReader(os.Stdin)))
}
