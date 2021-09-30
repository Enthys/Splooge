package pkg_test

import (
	"wildfire/pkg"
	"testing"
)

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
