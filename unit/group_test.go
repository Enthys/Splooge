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
				Projects: map[string]*pkg.ProjectConfig{},
				Groups:   map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)

			_, _ = groupService.CreateGroup("foo")

			_, ok := config.Groups["foo"]

			if ok == false {
				t.Error("Failed to find group in configuration")
			}
		})

		t.Run("should return a empty project group", func(t *testing.T) {
			groupService := pkg.NewGroupService(&pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups:   map[string]*pkg.GroupConfig{},
			})

			group, _ := groupService.CreateGroup("foo")

			if len(*group) != 0 {
				t.Errorf("Group should have been returned empty. Group has length of %d", len(*group))
			}
		})

		t.Run("should return an error if the group already exists in the configuration", func(t *testing.T) {
			groupService := pkg.NewGroupService(&pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups: map[string]*pkg.GroupConfig{
					"foo": {},
				},
			})

			_, err := groupService.CreateGroup("foo")
			if err == nil {
				t.Error("CreateGroup should have returned an error instead of resolving")
			}

			expectedErrorMessage := "group with name 'foo' already exists"
			if err.Error() != expectedErrorMessage {
				t.Errorf("Unexpected error has been returned. Expected '%s' received '%s'", expectedErrorMessage, err)
			}
		})
	})

	t.Run("GetGroup", func(t *testing.T) {
		t.Run("should return the group if it exists", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "bar", "zaz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups: map[string]*pkg.GroupConfig{"foo": group},
			}
			groupService := pkg.NewGroupService(config)

			g := groupService.GetGroup("foo")
			if g != group {
				t.Errorf("Failed to retrieve group from configuration. Expected '%v' received '%v'", group, g)
			}
		})

		t.Run("should return nil if the group doesnot exist", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)

			g := groupService.GetGroup("foo")
			if g != nil {
				t.Errorf("Failed to retrieve group from configuration. Expected '%v' received '%v'", nil, g)
			}
		})
	})

	t.Run("DeleteGroup", func(t *testing.T) {
		t.Run("should remove the group from the provided configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{Groups: map[string]*pkg.GroupConfig{
				"foo": {"bar", "zaz"},
			}}
			groupService := pkg.NewGroupService(config)

			groupService.DeleteGroup("foo")

			if len(config.Groups) != 0 {
				t.Error("Group was not removed from configuration")
			}
		})
	})

	t.Run("AddProject", func(t *testing.T) {
		t.Run("should add the name of project to its collection", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "https://github.com/foo"},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			group, err := groupService.AddProject(group, "foo")

			if err != nil {
				t.Error("Failed to add project to group")
			}

			if len(*group) == 0 {
				t.Error("ProjectConfig was not added to GroupConfig collection")
			}
		})

		t.Run("should return an error if project does not exist in provided configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			_, err := groupService.AddProject(group, "foo")

			if err == nil {
				t.Error("An error should have been returned.")
			}
		})

		t.Run("should not add projects which are already in the group", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{"foo": {"", "", ""}},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			updatedGroup, _ := groupService.AddProject(group, "foo")

			if len(*updatedGroup) > 2 {
				t.Error("ProjectConfig should not have been added to group twice.")
			}
		})
	})

	t.Run("RemoveProject", func(t *testing.T) {
		t.Run("should remove the project from the group", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"", "", ""},
					"bar": {"", "", ""},
					"zaz": {"", "", ""},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			group, _ = groupService.AddProject(group, "foo")
			group, _ = groupService.AddProject(group, "bar")
			group, _ = groupService.AddProject(group, "zaz")
			result := groupService.RemoveProject(group, "bar")


			if len(*result) != 2 {
				t.Error(fmt.Sprintf("Group has invalid number of projects. Expected %d received %d", 2, len(*result)))
			}

			expectedGroup := pkg.GroupConfig{"foo", "zaz"}
			if !reflect.DeepEqual(*result, expectedGroup) {
				t.Errorf(
					"Expected group does not match result. Expected %+v received %+v",
					expectedGroup,
					*result,
				)
			}
		})

		t.Run("should return the same project group if project does not exist", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"", "", ""},
					"bar": {"", "", ""},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			group, _ = groupService.AddProject(group, "foo")
			group, _ = groupService.AddProject(group, "bar")
			result := groupService.RemoveProject(group, "zaz")

			expected := pkg.GroupConfig{"foo", "bar"}
			if !reflect.DeepEqual(expected, *result) {
				t.Errorf(
					"Expected group does not match result. Expected %+v received %+v",
					expected,
					*result,
				)
			}
		})
	})

	t.Run("HasProject", func(t *testing.T) {
		t.Run("should return true if project exists in group", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"", "", ""},
					"bar": {"", "", ""},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			group, _ = groupService.AddProject(group, "foo")

			if groupService.HasProject(group, "foo") == false {
				t.Error(
					fmt.Sprintf("HasProject should have returned true. Expected '%t' Received '%t'", true, false),
				)
			}
		})

		t.Run("should return false if project does not exist in group", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"", "", ""},
					"bar": {"", "", ""},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			groupService := pkg.NewGroupService(config)
			group, _ := groupService.CreateGroup("foo")
			group, _ = groupService.AddProject(group, "foo")

			if groupService.HasProject(group, "zaz") == true {
				t.Error(
					fmt.Sprintf("HasProject should have returned false. Expected '%t' Received '%t'", false, true),
				)
			}
		})
	})
}
