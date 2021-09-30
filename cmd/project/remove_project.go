package project

import (
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

var removeProjectCmd = &cobra.Command{
	Use:   "remove project-name...",
	Short: "Remove project from configuration and all groups.",
	Run: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool) {
		if len(args) == 0 {
			return config, false
		}

		for _, projectName := range args {
			config.RemoveProject(projectName)

			emoji.Println(":cloud: Removing project: ", projectName)
		}

		return config, true
	}),
}
