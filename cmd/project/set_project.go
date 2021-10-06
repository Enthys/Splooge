package project

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"wildfire/pkg"
)

func NewSetProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set name type url",
		Short: "Update or create a Project in the configuration.",
		Long: fmt.Sprintf(
			`Update or create a Project in the configuration.
Warning! If a project already exists with the provided name it will be overwritten.

name - The name of the project. Will be used to store in the configuration groups.
type - The project location type.
    Available options are: %s
url - The location through which to retrieve clone the project

		`,
			func() string {
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
			project := config.GetProject(args[0])
			if project != nil &&
				!requestUserApproval(emoji.Sprintf(
					"Project '%s' already exists. Do you wish to overwrite the project configuration?",
					args[0],
				)) {
				return config, false, nil
			}

			config.SetProject(&pkg.Project{
				Name: args[0],
				Type: pkg.ProjectType(args[1]),
				URL:  pkg.ProjectPath(args[2]),
			})

			emoji.Println(":fire: Setting project!")
			fmt.Println("    -> Name: ", args[0])
			fmt.Println("    -> Type: ", args[1])
			fmt.Println("    -> URL: ", args[2])

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}

func requestUserApproval(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message, "[y/N]: ")

	text, _, _ := reader.ReadRune()

	return strings.Compare(strings.ToLower(string(text)), "y") == 0
}
