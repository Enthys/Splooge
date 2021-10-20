package clone

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"wildfire/pkg"
	"wildfire/pkg/project_repository"
)

func NewPullProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "project <project name> [path]",
		Short: "Pull group projects from their repositories",
		Long:  `Pull the projects stored in the specified group in the current directory or a specified directory.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			projectService := pkg.NewProjectService(config)
			groupService := pkg.NewGroupService(config)
			projectRepoService := project_repository.NewProjectRepositoryService(
				&projectService,
				&groupService,
				project_repository.NewCloner(os.Stdout),
			)
			projectName := args[0]

			project := projectService.GetProject(projectName)

			if project == nil {
				return config, false, emoji.Errorf("Project '%s' does not exist in configuration.", projectName)
			}

			var pullPath string

			if len(args) > 1 {
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", args[1], args[0]))
			} else {
				currentWD, _ := os.Getwd()
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", currentWD, args[0]))
			}

			err := projectRepoService.PullProject(
				filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)),
				project,
			)

			if err != nil {
				err = emoji.Errorf("Failed to clone project '%s'. Error: ", err)
				err = os.RemoveAll(filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)))
			}

			return config, false, nil
		}),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}
