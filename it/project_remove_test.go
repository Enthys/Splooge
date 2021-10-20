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

		config := pkg.GetConfig()
		if len(config.Projects) != 2 {
			t.Errorf("Failed to remove project from configuration. Expected '%d' projects found '%d'", 2, len(config.Projects))
		}
	})


	t.Run("Should remove the all provided projects from the configuration file", func(t *testing.T) {
		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		_, _ = ps.AddProject("foo", "github.com/example", pkg.ProjectTypeGit)
		_ = config.SaveConfig()

		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{"foo", "bar"})
		_ = cmd.Execute()

		config = pkg.GetConfig()
		if len(config.Projects) != 1 {
			t.Errorf("Failed to remove projects from configuration. Expected '%d' projects found '%d'", 1, len(config.Projects))
		}
	})

	t.Run("Should not update the configuration if no project names have been provided", func(t *testing.T) {
		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)

		_, _ = ps.AddProject("foo", "github.com/example", pkg.ProjectTypeGit)
		_, _ = ps.AddProject("bar", "github.com/example", pkg.ProjectTypeGit)
		_ = config.SaveConfig()

		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{})
		_ = cmd.Execute()

		config = pkg.GetConfig()
		if len(config.Projects) != 3 {
			t.Errorf("Invalid number of projects from configuration. Expected '%d' projects found '%d'", 3, len(config.Projects))
		}
	})

	t.Run("Should remove removed project from groups which contain said project", func(t *testing.T) {
		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		gs := pkg.NewGroupService(config)

		_, _ = ps.AddProject("foo", "github.com/example", pkg.ProjectTypeGitLab)
		_, _ = ps.AddProject("bar", "github.com/example", pkg.ProjectTypeGitLab)
		_, _ = ps.AddProject("zaz", "github.com/example", pkg.ProjectTypeGitLab)
		group, _ := gs.CreateGroup("foo")
		group, _ = gs.AddProject(group, "foo")
		group, _ = gs.AddProject(group, "bar")
		group, _ = gs.AddProject(group, "zaz")
		_ = config.SaveConfig()

		cmd := project.NewRemoveProjectCmd()
		cmd.SetArgs([]string{"foo"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("RemoveProject should have not returned an error. Error: %s", err)
		}

		config = pkg.GetConfig()
		gs = pkg.NewGroupService(config)
		group = gs.GetGroup("foo")
		for _, projectName := range *group {
			if projectName == "foo" {
				t.Error("Should not have found project 'foo' in group")
			}
		}
	})
}
