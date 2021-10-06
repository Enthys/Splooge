package project

import (
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewRemoveProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove project-name...",
		Short: "Remove project from configuration and all groups.",
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			if len(args) == 0 {
				return config, false, nil
			}

			for _, projectName := range args {
				config.RemoveProject(projectName)

				emoji.Println(":cloud: Removing project: ", projectName)
			}

			return config, true, nil
		}),
	}
}
