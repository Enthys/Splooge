package group

import (
	"errors"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewRemoveProjectFromGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-project <name> <project>...",
		Short: "Remove project to group",
		Long: `Remove project to group.

Will remove projects from group if group exists.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			groupService := pkg.NewGroupService(config)

			groupName := args[0]

			group := groupService.GetGroup(groupName)
			if group == nil {
				return config, false, emoji.Errorf("Group '%s' does not exist", groupName)
			}

			projectNames := args[1:]

			for _, name := range projectNames {
				group = groupService.RemoveProject(group, name)
			}

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}
