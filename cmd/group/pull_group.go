package group

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
	"wildfire/pkg"
)

func NewPullGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <group name> [path]",
		Short: "Pull group projects from their repositories",
		Long:  `Pull the projects stored in the specified group in the current directory or a specified directory.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			group, ok := config.Groups[args[0]]
			if ok == false {
				return config, false, fmt.Errorf("group '%s' does not exist in configuration", args[0])
			}

			var pullPath string

			if len(args) > 1 {
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", args[1], args[0]))
			} else {
				currentWD, _ := os.Getwd()
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", currentWD, args[0]))
			}

			var err error
			var wg sync.WaitGroup
			wg.Add(len(group))

			for _, projectName := range group {
				go func(projectName string) {
					err = config.GetProject(projectName).PullProject(
						filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)),
					)

					if err != nil {
						err = emoji.Errorf("Failed to clone project '%s'. Error: ", err)
					}
					wg.Done()
				}(projectName)
			}

			wg.Wait()

			if err != nil {
				for _, projectName := range group {
					err = os.RemoveAll(filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)))
					if os.IsNotExist(err) {
						break;
					}
				}
			}

			return config, false, nil
		}),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}
