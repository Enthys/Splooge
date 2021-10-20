package clone

import "github.com/spf13/cobra"

var CloneCmd = &cobra.Command{
	Use: "clone",
}

func init() {
	CloneCmd.AddCommand(NewPullProjectCmd())
	CloneCmd.AddCommand(NewPullGroupCmd())
}
