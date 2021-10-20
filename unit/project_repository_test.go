package unit_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"wildfire/pkg"
	"wildfire/pkg/project_repository"
)

type ClonerMock struct {
	StubCloneProject func(path string, project *pkg.ProjectConfig) error
}

func (c *ClonerMock) CloneProject(path string, project *pkg.ProjectConfig) error {
	return c.StubCloneProject(path, project)
}

func TestProjectRepository(t *testing.T) {
	t.Run("PullProject", func(t *testing.T) {
		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		gs := pkg.NewGroupService(config)
		project := &pkg.ProjectConfig{
			Name: "foo",
			Type: pkg.ProjectTypeGit,
			URL:  "github.com/foo/bar",
		}
		ps.UpdateOrCreate(project)

		t.Run("Should use the provided cloner object when cloning a project", func(t *testing.T) {
			used := false
			cloner := &ClonerMock{
				StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
					used = true
					return nil
				},
			}

			pr := project_repository.NewProjectRepositoryService(&ps, &gs, cloner)

			err := pr.PullProject("./testdata", project)
			if err != nil {
				t.Errorf("Test should not have returned an error. Error: %s", err)
			}

			if used != true {
				t.Errorf("Variable 'used' should have been set to true. Expected '%t' received '%t'", true, used)
			}
		})

		t.Run("should return an error if the Cloner returns an error", func(t *testing.T) {
			expectedErr := errors.New("test error")
			cloner := &ClonerMock{
				StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
					return expectedErr
				},
			}

			pr := project_repository.NewProjectRepositoryService(&ps, &gs, cloner)

			err := pr.PullProject("./testdata", project)
			if err == nil {
				t.Errorf("Expected error to be returned, instead method resolved.")
			}

			if err != expectedErr {
				t.Errorf(
					"Different error than expected has been returned. Expected '%s' received '%s'",
					expectedErr,
					err,
				)
			}
		})
	})

	t.Run("PullGroup", func(t *testing.T) {
		t.Run("should return nil if no issues occurred which cloning group", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "bar", "zaz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"bar": {"bar", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"zaz": {"zaz", pkg.ProjectTypeGit, "github.com/foo/bar"},
				},
				Groups: map[string]*pkg.GroupConfig{
					"foo": group,
				},
			}
			ps := pkg.NewProjectService(config)
			gs := pkg.NewGroupService(config)
			pr := project_repository.NewProjectRepositoryService(
				&ps,
				&gs,
				&ClonerMock{
					StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
						return nil
					},
				},
			)
			err := pr.PullGroup("./testdata", group)
			if err != nil {
				t.Errorf("PullGroup should have returned 'nil' instead it returned an error. Error: '%s'", err)
			}
		})

		t.Run("should return an error if group contains project which does not exist in configuration", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "bar", "zaz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"bar": {"bar", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"zaz": {"zaz", pkg.ProjectTypeGit, "github.com/foo/bar"},
				},
				Groups: map[string]*pkg.GroupConfig{
					"foo": group,
				},
			}
			ps := pkg.NewProjectService(config)
			gs := pkg.NewGroupService(config)
			pr := project_repository.NewProjectRepositoryService(
				&ps,
				&gs,
				&ClonerMock{
					StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
						return nil
					},
				},
			)

			tests := []struct {
				Group            *pkg.GroupConfig
				ExpectedInString []string
			}{
				{
					&pkg.GroupConfig{"foo", "bar", "missing_project"},
					[]string{"project 'missing_project' does not exist in configuration"},
				},
				{
					&pkg.GroupConfig{"foo", "missing_project", "missing_project_1"},
					[]string{
						"project 'missing_project' does not exist in configuration",
						"project 'missing_project_1' does not exist in configuration",
					},
				},
				{
					&pkg.GroupConfig{"foo", "bar", "missing_project", "taz", "paz"},
					[]string{
						"project 'missing_project' does not exist in configuration",
						"project 'taz' does not exist in configuration",
						"project 'paz' does not exist in configuration",
					},
				},
			}

			for _, testCase := range tests {
				err := pr.PullGroup("./testdata", testCase.Group)
				if err == nil {
					t.Errorf("PullGroup should have returned an error instead it resolved.")
				}

				for _, expected := range testCase.ExpectedInString {
					if strings.Contains(err.Error(), expected) == false {
						t.Errorf("Expected error to contain project name '%s'. Error: %s", expected, err)
					}
				}
			}
		})

		t.Run("should return error if the cloner returns an error while cloning the group project", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "bar", "zaz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"bar": {"bar", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"zaz": {"zaz", pkg.ProjectTypeGit, "github.com/foo/bar"},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			ps := pkg.NewProjectService(config)
			gs := pkg.NewGroupService(config)
			pr := project_repository.NewProjectRepositoryService(
				&ps,
				&gs,
				&ClonerMock{
					StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
						return fmt.Errorf("failed to clone project '%s'", project.Name)
					},
				},
			)
			tests := []struct {
				Group            *pkg.GroupConfig
				ExpectedInString []string
			}{
				{
					&pkg.GroupConfig{"foo", "bar"},
					[]string{
						"failed to clone project 'foo'",
						"failed to clone project 'bar'",
					},
				},
				{
					&pkg.GroupConfig{"foo"},
					[]string{
						"failed to clone project 'foo'",
					},
				},
				{
					&pkg.GroupConfig{"foo", "bar", "zaz"},
					[]string{
						"failed to clone project 'foo'",
						"failed to clone project 'bar'",
						"failed to clone project 'zaz'",
					},
				},
			}

			for _, testCase := range tests {
				err := pr.PullGroup("./testdata", group)
				if err == nil {
					t.Errorf("PullGroup should have returned an error instead it resolved.")
				}

				for _, expected := range testCase.ExpectedInString {
					if strings.Contains(err.Error(), expected) == false {
						t.Errorf("Expected error to contain project name '%s'. Error: %s", expected, err)
					}
				}
			}
		})
	})

	t.Run("PullProjectsFromGroup", func(t *testing.T) {
		t.Run("Should return nil if to issues have occurred", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "bar", "zaz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"bar": {"bar", pkg.ProjectTypeGit, "github.com/foo/bar"},
					"zaz": {"zaz", pkg.ProjectTypeGit, "github.com/foo/bar"},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			ps := pkg.NewProjectService(config)
			gs := pkg.NewGroupService(config)
			pr := project_repository.NewProjectRepositoryService(
				&ps,
				&gs,
				&ClonerMock{
					StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
						return nil
					},
				},
			)

			err := pr.PullProjectsFromGroup("./testdata", group, "foo", "bar", "zaz")
			if err != nil {
				t.Errorf("PullProjectsFromGroup should have resolved, instead it returned an error. Error '%s'", err)
			}
		})

		t.Run("should return an error if we've passed a project which does not exist in group or config", func(t *testing.T) {
			group := &pkg.GroupConfig{"foo", "taz"}
			config := &pkg.WildFireConfig{
				Projects: map[string]*pkg.ProjectConfig{
					"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
				},
				Groups: map[string]*pkg.GroupConfig{},
			}
			ps := pkg.NewProjectService(config)
			gs := pkg.NewGroupService(config)
			pr := project_repository.NewProjectRepositoryService(
				&ps,
				&gs,
				&ClonerMock{
					StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
						return nil
					},
				},
			)

			tests := []struct {
				Projects []string
				ExpectedStrings []string
			}{
				{
					Projects: []string{"foo", "bar", "taz"},
					ExpectedStrings: []string{
						"group does not contain project 'bar'",
						"project 'taz' does not exist",
					},
				},
				{
					Projects: []string{"zaz"},
					ExpectedStrings: []string{
						"group does not contain project 'zaz'",
					},
				},
				{
					Projects: []string{"foo", "bar", "zaz", "taz"},
					ExpectedStrings: []string{
						"group does not contain project 'bar'",
						"group does not contain project 'zaz'",
						"project 'taz' does not exist",
					},
				},
			}

			for _, testCase := range tests {
				err := pr.PullProjectsFromGroup("./testdata", group, testCase.Projects...)
				for _, expectedString := range testCase.ExpectedStrings {
					if strings.Contains(err.Error(), expectedString) == false {
						t.Errorf("Expected error from PullProjectsFromGroup to contain '%s'. Error: '%s'", expectedString, err)
					}
				}
			}
		})
	})

	t.Run("should return an error if the cloner fails to clone the projects", func(t *testing.T) {
		group := &pkg.GroupConfig{"foo", "bar", "zaz"}
		config := &pkg.WildFireConfig{
			Projects: map[string]*pkg.ProjectConfig{
				"foo": {"foo", pkg.ProjectTypeGit, "github.com/foo/bar"},
				"bar": {"bar", pkg.ProjectTypeGit, "github.com/foo/bar"},
				"zaz": {"zaz", pkg.ProjectTypeGit, "github.com/foo/bar"},
			},
			Groups: map[string]*pkg.GroupConfig{},
		}
		ps := pkg.NewProjectService(config)
		gs := pkg.NewGroupService(config)
		pr := project_repository.NewProjectRepositoryService(
			&ps,
			&gs,
			&ClonerMock{
				StubCloneProject: func(path string, project *pkg.ProjectConfig) error {
					return fmt.Errorf("failed to clone project '%s'", project.Name)
				},
			},
		)

		tests := []struct {
			Projects []string
			ExpectedStrings []string
		}{
			{
				Projects: []string{"foo"},
				ExpectedStrings: []string{
					"failed to clone project 'foo'",
				},
			},
			{
				Projects: []string{"foo", "bar"},
				ExpectedStrings: []string{
					"failed to clone project 'foo'",
					"failed to clone project 'bar'",
				},
			},
			{
				Projects: []string{"foo", "bar", "zaz"},
				ExpectedStrings: []string{
					"failed to clone project 'foo'",
					"failed to clone project 'bar'",
					"failed to clone project 'zaz'",
				},
			},
		}

		for _, testCase := range tests {
			err := pr.PullProjectsFromGroup("./testdata", group, testCase.Projects...)
			for _, expectedErrString := range testCase.ExpectedStrings {
				if strings.Contains(err.Error(), expectedErrString) == false {
					t.Errorf("Expected error to contain '%s'. Error '%s'", expectedErrString, err)
				}
			}
		}
	})

}
