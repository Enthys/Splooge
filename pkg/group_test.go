package pkg_test

import (
	"fmt"
	"reflect"
	"testing"
	"wildfire/pkg"
)

//fooProject = pkg.Project{"foo", pkg.ProjectTypeGit, "https://github.com/foo"}
//config = &pkg.WildFireConfig{Projects: map[string]pkg.Project{ }}

func TestProjectGroup_AddProject_should_add_the_name_of_project_to_its_collection(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{
		"foo": pkg.Project{"foo", pkg.ProjectTypeGit, "https://github.com/foo"},
	}}
	group := pkg.ProjectGroup{}
	group, err := group.AddProject(config, "foo")

	if err != nil {
		t.Error("Failed to add project to group")
	}

	if len(group) == 0 {
		t.Error("Project was not added to ProjectGroup collection")
	}
}

func TestProjectGroup_AddProject_should_return_an_error_if_project_does_not_exist_in_provided_configuration(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
	group := &pkg.ProjectGroup{}
	_, err := group.AddProject(config, "foo")

	if err == nil {
		t.Error("An error should have been returned.")
	}
}

func TestProjectGroup_AddProject_should_not_add_projects_which_are_already_in_the_group(t *testing.T) {
	config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{"foo": pkg.Project{"", "", ""}}}
	group := pkg.ProjectGroup{"foo", "bar"}
	group, _ = group.AddProject(config, "foo")

	if len(group) > 2 {
		t.Error("Project should not have been added to group twice.")
	}
}

func TestProjectGroup_RemoveProject_should_remove_the_project_from_the_group(t *testing.T) {
	group := &pkg.ProjectGroup{"foo", "bar", "zaz"}
	result := group.RemoveProject("bar")

	if len(result) != 2 {
		t.Error(fmt.Sprintf("Group has invalid number of projects. Expected %d received %d", len(result), 3))
	}

	expectedGroup := pkg.ProjectGroup{"foo", "zaz"}
	if !reflect.DeepEqual(result, expectedGroup) {
		t.Error(fmt.Sprintf(
			"Expected group does not match result. Expected %+v received %+v",
			group,
			expectedGroup,
		))
	}
}

func TestProjectGroup_HasProject_should_return_true_if_project_exists_in_group(t *testing.T) {
	group := &pkg.ProjectGroup{"foo", "bar"}
	projectName := "foo"

	if group.HasProject(&projectName) == false {
		t.Error(
			fmt.Sprintf("HasProject should have returned true. Expected '%t' Received '%t'", true, false),
		)
	}
}

func TestProjectGroup_HasProject_should_return_false_if_project_does_not_exist_in_group(t *testing.T) {
	group := &pkg.ProjectGroup{"foo", "bar"}
	projectName := "zaz"

	if group.HasProject(&projectName) == true {
		t.Error(
			fmt.Sprintf("HasProject should have returned false. Expected '%t' Received '%t'", false, true),
		)
	}
}
