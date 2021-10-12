package unit_test

import (
	"fmt"
	"reflect"
	"testing"
	"wildfire/pkg"
)

func TestProjectGroup(t *testing.T) {
	t.Run("CreateGroup", func(t *testing.T) {
		t.Run("should add the group to the provided configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]pkg.Project{},
				Groups:   map[string]pkg.ProjectGroup{},
			}

			_ = pkg.CreateGroup(config, "foo")

			_, ok := config.Groups["foo"]

			if !ok {
				t.Error("Failed to find group in configuration")
			}
		})

		t.Run("should return a empty project group", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]pkg.Project{},
				Groups:   map[string]pkg.ProjectGroup{},
			}

			group := pkg.CreateGroup(config, "foo")

			if len(*group) != 0 {
				t.Errorf("Group should have been returned empty. Group has length of %d", len(*group))
			}
		})
	})

	t.Run("RemoveGroup", func(t *testing.T) {
		t.Run("should remove the group from the provided configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{Groups: map[string]pkg.ProjectGroup{
				"foo": {"bar", "zaz"},
			}}

			pkg.RemoveGroup(config, "foo")

			if len(config.Groups) != 0 {
				t.Error("Group was not removed from configuration")
			}
		})
	})

	t.Run("AddProject", func(t *testing.T) {
		t.Run("should add the name of project to its collection", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{
				"foo": {"foo", pkg.ProjectTypeGit, "https://github.com/foo"},
			}}
			group := &pkg.ProjectGroup{}
			group, err := group.AddProject(config, "foo")

			if err != nil {
				t.Error("Failed to add project to group")
			}

			if len(*group) == 0 {
				t.Error("Project was not added to ProjectGroup collection")
			}
		})

		t.Run("should_return_an_error_if_project_does_not_exist_in_provided_configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{}}
			group := &pkg.ProjectGroup{}
			_, err := group.AddProject(config, "foo")

			if err == nil {
				t.Error("An error should have been returned.")
			}
		})

		t.Run("should not add projects which are already in the group", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]pkg.Project{"foo": {"", "", ""}}}
			group := &pkg.ProjectGroup{"foo", "bar"}
			group, _ = group.AddProject(config, "foo")

			if len(*group) > 2 {
				t.Error("Project should not have been added to group twice.")
			}
		})
	})

	t.Run("RemoveProject", func(t *testing.T) {
		t.Run("should remove the project from the group", func(t *testing.T) {
			group := &pkg.ProjectGroup{"foo", "bar", "zaz"}
			result := group.RemoveProject("bar")

			if len(result) != 2 {
				t.Error(fmt.Sprintf("Group has invalid number of projects. Expected %d received %d", len(result), 3))
			}

			expectedGroup := pkg.ProjectGroup{"foo", "zaz"}
			if !reflect.DeepEqual(result, expectedGroup) {
				t.Errorf(
					"Expected group does not match result. Expected %+v received %+v",
					expectedGroup,
					group,
				)
			}
		})

		t.Run("should return the same project group if project does not exist", func(t *testing.T) {
			group := &pkg.ProjectGroup{"foo", "bar"}
			result := group.RemoveProject("zaz")

			if !reflect.DeepEqual(group, &result) {
				t.Errorf(
					"Expected group does not match result. Expected %+v received %+v",
					group,
					&result,
				)
			}
		})
	})

	t.Run("HasProject", func(t *testing.T) {
		t.Run("should return true if project exists in group", func(t *testing.T) {
			group := &pkg.ProjectGroup{"foo", "bar"}
			projectName := "foo"

			if group.HasProject(&projectName) == false {
				t.Error(
					fmt.Sprintf("HasProject should have returned true. Expected '%t' Received '%t'", true, false),
				)
			}
		})

		t.Run("should return false if project does not exist in group", func(t *testing.T) {
			group := &pkg.ProjectGroup{"foo", "bar"}
			projectName := "zaz"

			if group.HasProject(&projectName) == true {
				t.Error(
					fmt.Sprintf("HasProject should have returned false. Expected '%t' Received '%t'", false, true),
				)
			}
		})
	})
}
