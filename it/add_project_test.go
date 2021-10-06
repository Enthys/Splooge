//+build integration

package it_test

import (
	"os"
	"testing"
	"wildfire/cmd/project"
	"wildfire/pkg"
)

func TestAddProject(t *testing.T) {
	t.Run("Should add project to configuration file", func(t *testing.T) {
		cfgName := getConfigFilePath("add_project_test.wildfire.yaml")
		_ = setConfig(cfgName)

		addProjectCmd := project.NewAddProjectCmd()
		addProjectCmd.SetArgs([]string{"foo", "git", "github.com/example"})
		err := addProjectCmd.Execute()

		if err != nil {
			t.Errorf("Command has failed to execute. Error: %s", err)
		}

		config := pkg.GetConfig()

		if len(config.Projects) != 1 {
			t.Errorf("Invalid number of projects found. Expected '%d' found '%d'", 1, len(config.Projects))
		}

		err = os.Remove(cfgName)

		if err != nil {
			t.Errorf("Failed to clean up test environment. Error: %s", err)
		}
	})

	t.Run("Should throw an error if the provided project type is invalid", func(t *testing.T) {
		cfgName := getConfigFilePath("add_project_test.wildfire.yaml")
		_ = setConfig(cfgName)

		addProjectCmd := project.NewAddProjectCmd()
		addProjectCmd.SetArgs([]string{"foo", "invalid", "github.com/example"})
		err := addProjectCmd.Execute()

		if err == nil {
			t.Errorf("Should have thrown error for invalid project type")
		}
	})

	t.Run("Should throw an error if the provided project already exists", func(t *testing.T) {
		cfgName := getConfigFilePath("add_project_test.wildfire.yaml")
		_ = os.Remove(cfgName)
		_ = setConfig(cfgName)

		createProject := func() error {
			addProjectCmd := project.NewAddProjectCmd()
			addProjectCmd.SetArgs([]string{"foo", "git", "github.com/example"})
			return addProjectCmd.Execute()
		}

		_ = createProject()
		err := createProject()
		if err == nil {
			t.Errorf("Should have thrown error for project name conflict")
		}

		err = os.Remove(cfgName)
		if err != nil {
			t.Errorf("Failed to clean up test environment. Error: %s", err)
		}
	})
}
