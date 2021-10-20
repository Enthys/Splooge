package group

import (
	"errors"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewAddProjectToGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-project <name> <project>...",
		Short: "Add project to group",
		Long: `Add project to group.

Will add projects to group if group exists.
If some project name does not exist in the configuration then no projects will be added. 
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			projectService := pkg.NewProjectService(config)
			groupService := pkg.NewGroupService(config)

			groupName := args[0]

			group := groupService.GetGroup(groupName)
			if group == nil {
				return config, false, emoji.Errorf("Group '%s' does not exist", groupName)
			}

			projectNames := args[1:]

			success := true
			for _, name := range projectNames {
				exists := projectService.HasProject(name)

				if exists == false {
					emoji.Println(":prohibited: ProjectConfig ", name, " does not exist in this configuration.")
					success = false
					continue
				}

				if success == true {
					emoji.Printf(":ocean: ProjectConfig '%s' has been added to group '%s'\n", name, groupName)
					group, _ = groupService.AddProject(group, name)
				}
			}

			if success == false {
				emoji.Println(":error: Group was not updated. Resolve issues and try again.")
				return config, false, nil
			}

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}
