package unit_test

import (
	"fmt"
	"testing"
	"wildfire/pkg"
)

func TestProjectType_ValidType(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		tests := []struct{
			Project pkg.ProjectConfig
			Expected bool
		}{
			{ pkg.ProjectConfig{"git", pkg.ProjectTypeGit, "git.com/foo"}, true },
			{ pkg.ProjectConfig{"gitlab", pkg.ProjectTypeGitLab, "git.com/foo"}, true },
			{ pkg.ProjectConfig{"bitbucket", pkg.ProjectTypeBitBucket, "git.com/foo"}, true },
			{ pkg.ProjectConfig{"fake", pkg.ProjectType("fake"), "git.com/foo"}, false },
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Should return %t if type is %s", test.Expected, test.Project.Type), func(t *testing.T) {
				if test.Project.Type.ValidType() != test.Expected {
					t.Errorf("Type validation returned an incorrect value. Expected for project '%s' to receive '%t' received '%t'", test.Project.Name, test.Expected, test.Project.Type.ValidType())
				}
			})
		}
	})

	t.Run("GetAvailableTypes", func(t *testing.T) {
		var tt pkg.ProjectType
		types := tt.GetAvailableTypes()

		if len(types) != 3 {
			t.Errorf(
				"Expected there to be %d number of types received but received %d(%v)",
				3,
				len(types),
				types,
			)
		}
	})
}

func TestProject(t *testing.T) {
	t.Run("AddProject", func(t *testing.T) {
		t.Run("should return nil if project was added", func(t *testing.T) {
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{},
			}
			projectService := pkg.NewProjectService(config)

			_, err := projectService.AddProject("foo", "github.com/foo/bar", pkg.ProjectTypeGit)

			if err != nil {
				t.Error("AddProject should have returned nil when adding a project which does not exist. Error", err)
			}
		})

		t.Run("should return an error if project already exists in config", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{}}
			projectService := pkg.NewProjectService(config)

			_, _ = projectService.AddProject("foo", "github.com/foo/bar", pkg.ProjectTypeGit)
			_, err := projectService.AddProject("foo", "github.com/foo/bar", pkg.ProjectTypeGit)

			if err == nil {
				t.Error("AddProject should have returned an error when adding a project which already exists in config.")
			}
		})
	})

	t.Run("GetProject", func(t *testing.T) {
		t.Run("should return nil if it does not find the requested project", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{}}
			projectService := pkg.NewProjectService(config)
			retrievedProject := projectService.GetProject("foo")

			if retrievedProject != nil {
				t.Error("Retrieved project should have been nil")
			}
		})
	})

	t.Run("RemoveProject", func(t *testing.T) {
		t.Run("should remove the project from the configuration", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{
				"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
				"bar": {"bar", pkg.ProjectTypeGit, "github.com/bar/bar"},
			}}

			projectService := pkg.NewProjectService(config)

			projectService.RemoveProject("foo")
			fooProject := projectService.GetProject("foo")
			if fooProject != nil {
				t.Error("Project was not removed from configuration")
			}

			barProject := projectService.GetProject("bar")
			if barProject == nil {
				t.Error("Project which should not have been removed has been removed")
			}
		})
	})

	t.Run("SetProject", func(t *testing.T) {
		t.Run("should add project if it does not exist", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{}}
			projectService := pkg.NewProjectService(config)

			projectService.UpdateOrCreate(&pkg.ProjectConfig{"foo", pkg.ProjectTypeGit, "github.com/foo/bar"})
			fooProject := projectService.GetProject("foo")

			if fooProject == nil {
				t.Error("Project was not set in configuration")
			}
		})

		t.Run("should replace a project if it already exists", func(t *testing.T) {
			config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{}}
			projectService := pkg.NewProjectService(config)

			projectService.UpdateOrCreate(&pkg.ProjectConfig{"foo", pkg.ProjectTypeGit, "github.com/foo/bar"})
			projectService.UpdateOrCreate(&pkg.ProjectConfig{"foo", pkg.ProjectTypeGitLab, "github.com/foo/bar/zaz"})
			fooProject := projectService.GetProject("foo")

			if fooProject.Type != pkg.ProjectTypeGitLab ||
				fooProject.URL != "github.com/foo/bar/zaz" {
				t.Error("Project has not overwritten old project")
			}
		})
	})

	t.Run("HasProject", func(t *testing.T) {
		config := &pkg.WildFireConfig{Projects: map[string]*pkg.ProjectConfig{
			"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
			"bar": {"bar", pkg.ProjectTypeGit, "github.com/bar/bar"},
		}}

		projectService := pkg.NewProjectService(config)

		t.Run("Should return true if the project exists in the configuration", func(t *testing.T) {
			result := projectService.HasProject("foo")

			if result != true {
				t.Errorf("HasProject should have returned true. Expected '%t' received '%t'", true, result)
			}
		})

		t.Run("Should return false if the project does not exist in the configuration", func(t *testing.T) {
			result := projectService.HasProject("zaz")

			if result != false {
				t.Errorf("HasProject should have returned false. Expected '%t' received '%t'", false, result)
			}
		})
	})
}
