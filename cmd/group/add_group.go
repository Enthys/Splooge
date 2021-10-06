package group

import (
	"errors"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewCreateGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name> [project...]",
		Short: "Create project group",
		Long: `Create new project group.

Project groups are groups which contain projects existing in the configuration.
Groups are primarily used when commands are to be executed on specific projects.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			group := pkg.CreateGroup(config, args[0])

			projectNames := args[1:]
			if len(projectNames) == 0 {
				return config, true, nil
			}

			success := true
			for _, name := range projectNames {
				exists := projectExists(config, name)

				if exists == false {
					emoji.Println(":prohibited: Project ", name, " does not exist in this configuration.")
					success = false
					continue
				}

				if success == true {
					emoji.Println(":check: Project ", name, " has been added to group ", args[0])
					group, _ = group.AddProject(config, name)
				}
			}

			if success == false {
				return nil, false, nil
			}

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}

func projectExists(config *pkg.WildFireConfig, projectName string) bool {
	for _, project := range config.Projects {
		if project.Name == projectName {
			return true
		}
	}

	return false
}
