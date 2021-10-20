package unit_test

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"
	"wildfire/pkg"
)



func TestGetConfig(t *testing.T) {
	t.Run("should return a wildfire config even if file does not exist", func(t *testing.T) {
		_ = setConfig("missing")
		config := pkg.GetConfig()

		if len(config.Projects) > 0 {
			t.Error("Expected config projects to be empty. Found ", len(config.Projects), " number of projects")
		}

		if len(config.Groups) > 0 {
			t.Error("Expected config groups to be empty. Found ", len(config.Groups), " number of groups")
		}
	})

	t.Run("should return a wildfire config with filled out projects", func(t *testing.T) {
		cfgFile := getConfigFilePath("projects.wildfire.yaml")

		err := setConfig(cfgFile)
		if err != nil {
			t.Error(err)
		}

		config := pkg.GetConfig()

		if len(config.Projects) != 5 {
			t.Error("Invalid number of projects found. Expected 5 received ", len(config.Projects))
		}

		if len(config.Groups) != 0 {
			t.Error("Invalid number of groups found. Expected 0 received", len(config.Groups))
		}
	})

	t.Run("should return an empty wildfire config if selected config is invalid", func(t *testing.T) {
		fmt.Println(viper.ConfigFileUsed())
		_ = setConfig("invalid.wildfire.yaml")
		config := pkg.GetConfig()

		if len(config.Projects) > 0 {
			t.Error("Expected config projects to be empty. Found ", len(config.Projects), " number of projects")
		}

		if len(config.Groups) > 0 {
			t.Error("Expected config groups to be empty. Found ", len(config.Groups), " number of groups")
		}
	})
}

func TestWildFireConfig(t *testing.T) {
	t.Run("SaveConfig", func(t *testing.T) {
		t.Run("should create the config file if it is missing", func(t *testing.T) {
			cfgFile := getConfigFilePath("new.wildfire.yaml")
			_ = deleteConfig(cfgFile)

			if _, err := os.Stat(cfgFile); os.IsExist(err) {
				t.Errorf("Failed to remove old test configuration file '%s'", cfgFile)
			}

			_ = setConfig(cfgFile)
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups:   map[string]*pkg.GroupConfig{},
			}

			config.Groups["foo"] = &pkg.GroupConfig{}
			config.Projects["example"] = &pkg.ProjectConfig{
				Name: "bar",
				Type: pkg.ProjectTypeGit,
				URL:  "git.com/bar/zaz",
			}

			err := config.SaveConfig()

			if err != nil {
				t.Errorf("Failed to save configuration file '%s'. Error: %s", cfgFile, err)
			}

			updatedConfig := pkg.GetConfig()

			if len(updatedConfig.Projects) != 1 {
				t.Errorf(
					"Invalid number of projects found in configuration. Expected 1 received %d", len(updatedConfig.Projects),
				)
			}

			if len(updatedConfig.Groups) != 1 {
				t.Errorf(
					"Invalid number of groups found in configuration. Expected 1 receivec %d", len(updatedConfig.Groups),
				)
			}

			err = deleteConfig(cfgFile)
			if err != nil {
				t.Error("Failed to clear test env")
			}
		})
	})
}
