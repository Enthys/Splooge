package project

import (
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewRemoveProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <project name>...",
		Short: "Remove project from configuration and all groups.",
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			if len(args) == 0 {
				return config, false, nil
			}

			projectService := pkg.NewProjectService(config)
			groupService := pkg.NewGroupService(config)

			for _, projectName := range args {
				projectService.RemoveProject(projectName)

				emoji.Println(":cloud: Removed project: ", projectName)

				for groupName, _ := range config.Groups {
					group := groupService.GetGroup(groupName)
					if groupService.HasProject(group, projectName) == true {
						group = groupService.RemoveProject(group, projectName)
						emoji.Println(emoji.Sprintf(":dash: Removed project '%s' from group '%s'", projectName, groupName))
					}
				}
			}

			return config, true, nil
		}),
	}
}
