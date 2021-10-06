package pkg_test

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
	"wildfire/pkg"
)

func setConfig(file string) error {
	dir, _ := os.Getwd()
	cfgFile := filepath.FromSlash(fmt.Sprintf("%s/pkg/testdata/%s", dir, file))
	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig()

	return fmt.Errorf("failed to load '%s' configuration file. Error %s", cfgFile, err)
}

func TestGetProject_should_return_a_wildfire_config_even_if_file_does_not_exist(t *testing.T) {
	_ = setConfig("missing")
	config := pkg.GetConfig()

	if len(config.Projects) > 0 {
		t.Error("Expected config projects to be empty. Found ", len(config.Projects), " number of projects")
	}

	if len(config.Groups) > 0 {
		t.Error("Expected config groups to be empty. Found ", len(config.Groups), " number of groups")
	}
}

func TestGetProject_should_return_an_empty_wildfire_config_if_selected_config_is_invalid(t *testing.T) {
	_ = setConfig("invalid.wildfire.yaml")
	config := pkg.GetConfig()

	if len(config.Projects) > 0 {
		t.Error("Expected config projects to be empty. Found ", len(config.Projects), " number of projects")
	}

	if len(config.Groups) > 0 {
		t.Error("Expected config groups to be empty. Found ", len(config.Groups), " number of groups")
	}
}

func TestGetConfig_should_return_a_wildfire_config_with_filled_out_projects(t *testing.T) {
	_ = setConfig("projects.wildfire.yaml")
	config := pkg.GetConfig()

	if len(config.Projects) != 5 {
		t.Error("Invalid number of projects found. Expected 5 received ", len(config.Projects))
	}

	if len(config.Groups) != 0 {
		t.Error("Invalid number of groups found. Expected 0 received", len(config.Groups))
	}
}

func TestWildFireConfig_SaveConfig_should_create_the_config_file_if_it_is_missing(t *testing.T) {
	//Removing the configuration if it exists
	file := "new.wildfire.yaml"
	dir, _ := os.Getwd()
	cfgFile := filepath.FromSlash(fmt.Sprintf("%s/pkg/testdata/%s", dir, file))
	_ = os.Remove(cfgFile)

	if _, err := os.Stat(cfgFile); os.IsExist(err) {
		t.Errorf("Failed to remove old test configuration file '%s'", cfgFile)
	}

	_ = setConfig(file)
	config := &pkg.WildFireConfig{
		Projects: map[string]pkg.Project{},
		Groups:   map[string]pkg.ProjectGroup{},
	}

	config.Groups["foo"] = pkg.ProjectGroup{}
	config.Projects["example"] = pkg.Project{
		Name: "bar",
		Type: pkg.ProjectTypeGit,
		URL:  pkg.ProjectPath("git.com/bar/zaz"),
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

	err = os.Remove(cfgFile)
	if err != nil {
		t.Error("Failed to clear test env")
	}
}

func TestWildFireConfig_AddProject_returns_nil_if_project_was_added(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	err := config.AddProject(&pkg.Project{
		Name: "foo",
		Type: pkg.ProjectTypeGit,
		URL:  "github.com/foo/bar",
	})

	if err != nil {
		t.Error("AddProject should have returned nil when adding a project which does not exist. Error", err)
	}
}

func TestWildFireConfig_AddProject_returns_an_error_if_project_already_exists_in_config(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	_ = config.AddProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar" })
	err := config.AddProject(&pkg.Project{ Name: "foo", Type: pkg.ProjectTypeGit, URL:  "github.com/foo/bar" })

	if err == nil {
		t.Error("AddProject should have returned an error when adding a project which already exists in config.")
	}
}

func TestWildFireConfig_ReplaceConfig_should_replace_the_project_in_the_configuration(t *testing.T) {
	project := &pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar" }
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{"foo": *project}}
	retrievedProject := config.GetProject("foo")

	if retrievedProject == nil {
		t.Error("Failed to retrieve project from configuration")
	}

	if retrievedProject.Name != project.Name ||
		retrievedProject.Type != project.Type ||
		retrievedProject.URL != project.URL {
		t.Error("Retrieved project does not match the requested object")
	}
}

func TestWildFireConfig_GetProject_should_return_nil_if_it_does_not_find_the_requested_project(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	retrievedProject := config.GetProject("foo")

	if retrievedProject != nil {
		t.Error("Retrieved project should have been nil")
	}
}

func TestWildFireConfig_RemoveProject_should_remove_the_project_from_the_configuration(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{
		"foo": pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar" },
		"bar": pkg.Project{"bar", pkg.ProjectTypeGit, "github.com/bar/bar" },
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
}

func TestWildFireConfig_SetProject_should_add_project_if_it_does_not_exist(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar" })
	fooProject := config.GetProject("foo")

	if fooProject == nil {
		t.Error("Project was not set in configuration")
	}
}

func TestWildFireConfig_SetProject_should_replace_a_project_if_it_already_exists(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/foo/bar" })
	config.SetProject(&pkg.Project{"foo", pkg.ProjectTypeGitLab, "github.com/foo/bar/zaz" })
	fooProject := config.GetProject("foo")

	if fooProject.Type != pkg.ProjectTypeGitLab ||
		fooProject.URL != "github.com/foo/bar/zaz" {
		t.Error("Project has not overwritten old project")
	}
}
