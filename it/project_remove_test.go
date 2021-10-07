//+build integration

package it_test

import (
	"os"
	"testing"
	"wildfire/cmd/project"
	"wildfire/pkg"
)

func TestRemoveProject(t *testing.T) {
	cfgFile := getConfigFilePath("project_remove.wildfire.yaml")
	err := initiateConfiguration(cfgFile)
	if err != nil {
		t.Errorf("Failed to initiate configuration. Error: %s", err)
	}

	config := pkg.GetConfig()

	defer func() {
		err = os.Remove(cfgFile)
		if err != nil {
			t.Errorf("Failed to tear down test environment. Error: %s", err)
		}
	}()

	t.Run("should remove the project from the configuration file", func(t *testing.T) {
		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{"foo"})
		_ = cmd.Execute()

		config = pkg.GetConfig()
		if len(config.Projects) != 2 {
			t.Errorf("Failed to remove project from configuration. Expected '%d' projects found '%d'", 2, len(config.Projects))
		}
	})

	_ = config.AddProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/example"})
	_ = config.SaveConfig()

	t.Run("Should remove the all provided projects from the configuration file", func(t *testing.T) {
		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{"foo", "bar"})
		_ = cmd.Execute()

		config = pkg.GetConfig()
		if len(config.Projects) != 1 {
			t.Errorf("Failed to remove projects from configuration. Expected '%d' projects found '%d'", 1, len(config.Projects))
		}
	})

	_ = config.AddProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/url"})
	_ = config.AddProject(&pkg.Project{"bar", pkg.ProjectTypeGit, "github.com/url"})
	_ = config.SaveConfig()

	t.Run("Should not update the configuration if no project names have been provided", func(t *testing.T) {
		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{})
		_ = cmd.Execute()

		config = pkg.GetConfig()
		if len(config.Projects) != 3 {
			t.Errorf("Invalid number of projects from configuration. Expected '%d' projects found '%d'", 3, len(config.Projects))
		}
	})
}
