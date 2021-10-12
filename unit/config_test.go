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
				Projects: map[string]pkg.Project{},
				Groups:   map[string]pkg.ProjectGroup{},
			}

			config.Groups["foo"] = pkg.ProjectGroup{}
			config.Projects["example"] = pkg.Project{
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

	t.Run("AddProject", func(t *testing.T) {
		t.Run("should return nil if project was added", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			err := config.AddProject(&pkg.Project{
				Name: "foo",
				Type: pkg.ProjectTypeGit,
				URL:  "github.com/foo/bar",
			})

			if err != nil {
				t.Error("AddProject should have returned nil when adding a project which does not exist. Error", err)
			}
		})

		t.Run("should return an error if project already exists in config", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			_ = config.AddProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar"})
			err := config.AddProject(&pkg.Project{Name: "foo", Type: pkg.ProjectTypeGit, URL: "github.com/foo/bar"})

			if err == nil {
				t.Error("AddProject should have returned an error when adding a project which already exists in config.")
			}
		})
	})

	t.Run("GetProject", func(t *testing.T) {
		t.Run("should return nil if it does not find the requested project", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			retrievedProject := config.GetProject("foo")

			if retrievedProject != nil {
				t.Error("Retrieved project should have been nil")
			}
		})
	})

	t.Run("RemoveProject", func(t *testing.T) {
		t.Run("should remove the project from the configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{
				"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
				"bar": {"bar", pkg.ProjectTypeGit, "github.com/bar/bar"},
			}}

			config.RemoveProject("foo")
			fooProject := config.GetProject("foo")
			if fooProject != nil {
				t.Error("Project was not removed from configuration")
			}

			barProject := config.GetProject("bar")
			if barProject == nil {
				t.Error("Project which should not have been removed has been removed")
			}
		})
	})

	t.Run("SetProject", func(t *testing.T) {
		t.Run("should add project if it does not exist", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar"})
			fooProject := config.GetProject("foo")

			if fooProject == nil {
				t.Error("Project was not set in configuration")
			}
		})
		t.Run("should replace a project if it already exists", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar"})
			config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGitLab, "github.com/foo/bar/zaz"})
			fooProject := config.GetProject("foo")

			if fooProject.Type != pkg.ProjectTypeGitLab ||
				fooProject.URL != "github.com/foo/bar/zaz" {
				t.Error("Project has not overwritten old project")
			}
		})
	})

	t.Run("HasProject", func(t *testing.T) {
		config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{
			"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
			"bar": {"bar", pkg.ProjectTypeGit, "github.com/bar/bar"},
		}}

		t.Run("Should return true if the project exists in the configuration", func(t *testing.T) {
			result := config.HasProject("foo")

			if result != true {
				t.Errorf("HasProject should have returned true. Expected '%t' received '%t'", true, result)
			}
		})

		t.Run("Should return false if the project does not exist in the configuration", func(t *testing.T) {
			result := config.HasProject("zaz")

			if result != false {
				t.Errorf("HasProject should have returned false. Expected '%t' received '%t'", false, result)
			}
		})
	})
}
