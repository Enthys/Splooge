package pkg_test

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"testing"
	"wildfire/pkg"
)

func TestProjectFunc(t *testing.T) {
	t.Run("should execute the provided function", func(t *testing.T) {
		executed := false

		testFunc := pkg.ProjectFunc(
			func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
				executed = true

				return config, false, nil
			},
		)

		testFunc(&cobra.Command{RunE: testFunc}, []string{})

		if executed != true {
			t.Error("Provided function is not executed")
		}
	})

	t.Run("should update the configuration if second argument is true", func(t *testing.T) {
		cfgFile := getConfigFilePath("new.wildfire.yaml")
		_ = deleteConfig(cfgFile)
		_ = setConfig(cfgFile)

		testFunc := pkg.ProjectFunc(
			func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
				return config, true, nil
			},
		)

		testFunc(&cobra.Command{RunE: testFunc}, []string{})

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			t.Error("Configuration was not created when 'true' was returned as second argument")
		}

		err := deleteConfig(cfgFile)
		if err != nil {
			t.Error("Failed to cleanup test environment. Error: ", err)
		}
	})

	t.Run("returns the provided error is something goes wrong", func(t *testing.T) {
		testFunc := pkg.ProjectFunc(
			func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			return config, false, fmt.Errorf("Error Message")
		})

		err := testFunc(&cobra.Command{RunE: testFunc}, []string{})
		if err == nil {
			t.Error("Test function should have returned an error, instead it returned nil")
		}
	})
}
