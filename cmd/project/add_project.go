package project

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"strings"
	"wildfire/pkg"
)

var addProjectCmd = &cobra.Command{
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
	Run: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool) {
		err := config.AddProject(&pkg.Project{
			Name: args[0],
			Type: pkg.ProjectType(args[1]),
			URL:  pkg.ProjectPath(args[2]),
		})
		cobra.CheckErr(err)

		emoji.Println(":fire: Adding new project!")
		emoji.Println("    -> Name: ", args[0])
		emoji.Println("    -> Type: ", args[1])
		emoji.Println("    -> URL: ", args[2])

		return config, true
	}),
}
