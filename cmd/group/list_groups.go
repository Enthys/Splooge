package group

import (
	"fmt"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewListGroupsCommand() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			groupService := pkg.NewGroupService(config)
			groupNames := groupService.GetGroupNames()

			if len(groupNames) != 0 {
				fmt.Println(fmt.Sprintf("Found %d groups in configuration:", len(groupNames)))
				for _, groupName := range groupNames {
					group := groupService.GetGroup(groupName)
					fmt.Println(fmt.Sprintf("- %s (%d projects)", groupName, len(*group)))
				}
			} else {
				fmt.Println("No groups were found in configuration.")
			}

			return config, false, nil
		}),
	}
}
