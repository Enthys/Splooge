package project

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"strings"
	"wildfire/pkg"
)

func NewAddProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add name type url",
		Short: "Add Project to the loaded configuration",
		Long: fmt.Sprintf(`Add a Project to the configuration.

name - The name of the project. Will be used to store in the configuration groups.
type - The project location type.
    Available options are: %s
url - The location through which to retrieve clone the project

`, func() string {
			var t pkg.ProjectType
			var availableTypes []string

			for _, projectType := range t.GetAvailableTypes() {
				availableTypes = append(availableTypes, string(projectType))
			}

			return strings.Join(availableTypes, ", ")
		}()),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("invalid number of arguments provided")
			}

			projectType := pkg.ProjectType(args[1])
			if projectType.ValidType() == false {
				return errors.New("invalid project type has been provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			err := config.AddProject(&pkg.Project{
				Name: args[0],
				Type: pkg.ProjectType(args[1]),
				URL:  pkg.ProjectPath(args[2]),
			})

			if err != nil {
				return nil, false, err
			}

			emoji.Println(":fire: Adding new project!")
			fmt.Println("    -> Name: ", args[0])
			fmt.Println("    -> Type: ", args[1])
			fmt.Println("    -> URL: ", args[2])

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}
