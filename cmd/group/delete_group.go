package group

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"wildfire/pkg"
)

func NewDeleteGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <group name>",
		Short: "Delete project group",
		Long: `Delete a project group from provided configuration`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			_, ok := config.Groups[args[0]]
			if ok == false {
				return config, false, fmt.Errorf("group '%s' does not exist in configuration", args[0])
			}

			delete(config.Groups, args[0])

			return config, true, nil
		}),
		SilenceUsage: true,
		SilenceErrors: true,
	}
}
