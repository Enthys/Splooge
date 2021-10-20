//+build integration

package it_test

import (
	"os"
	"testing"
	"wildfire/cmd/project"
	"wildfire/pkg"
)

type MockInputReader struct {
	Char rune
	project.CharacterInputReader
}

func NewMockInputReader(char rune) MockInputReader {
	return MockInputReader{
		Char: char,
	}
}

func (r MockInputReader) ReadRune() (rune, int, error) {
	return r.Char, 0, nil
}

func TestSetProject(t *testing.T) {
	cfgFile := getConfigFilePath("project_set.wildfire.yaml")
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

	t.Run("Should not replace the project configuration if we don't have user approval", func(t *testing.T) {
		cmd := project.NewSetProjectCmd(NewMockInputReader('N'))
		cmd.SetArgs([]string{"foo", string(pkg.ProjectTypeGitLab), "github.com/example/new"})
		err = cmd.Execute()
		if err != nil {
			t.Errorf("Failed to Set ProjectConfig. Command should have resolved. Error: %s", err)
		}

		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		updatedPackage := ps.GetProject("foo")
		if updatedPackage.Type != pkg.ProjectTypeGit {
			t.Errorf(
				"ProjectConfig type was not updated. Expected '%s' received '%s'",
				pkg.ProjectTypeGit,
				updatedPackage.Type,
			)
		}
		if updatedPackage.URL != "github.com/url" {
			t.Errorf(
				"ProjectConfig URL was not updated. Expected '%s' received '%s'",
				"github.com/url",
				updatedPackage.URL,
			)
		}
	})

	t.Run("Should replace the project configuration if it already exists and we have user approval", func(t *testing.T) {
		cmd := project.NewSetProjectCmd(NewMockInputReader('y'))
		cmd.SetArgs([]string{"foo", string(pkg.ProjectTypeGitLab), "github.com/example/new"})
		err = cmd.Execute()
		if err != nil {
			t.Errorf("Failed to Set ProjectConfig. Command should have resolved. Error: %s", err)
		}

		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		updatedPackage := ps.GetProject("foo")
		if updatedPackage.Type != pkg.ProjectTypeGitLab {
			t.Errorf(
				"ProjectConfig type was not updated. Expected '%s' received '%s'",
				pkg.ProjectTypeGitLab,
				updatedPackage.Type,
			)
		}
		if updatedPackage.URL != "github.com/example/new" {
			t.Errorf(
				"ProjectConfig URL was not updated. Expected '%s' received '%s'",
				"github.com/example/new",
				updatedPackage.URL,
			)
		}
	})

	t.Run("Should not request input if the project does not exist", func(t *testing.T) {
		cmd := project.NewSetProjectCmd(NewMockInputReader('y'))
		cmd.SetArgs([]string{"baz", string(pkg.ProjectTypeBitBucket), "bitbucket.com/example/bar"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Command should not have returned an error. Error %s", err)
		}

		config := pkg.GetConfig()
		ps := pkg.NewProjectService(config)
		newProject := ps.GetProject("baz")
		if newProject.Type != pkg.ProjectTypeBitBucket {
			t.Errorf(
				"New ProjectConfig has invalid type. Expected '%s' received '%s'",
				pkg.ProjectTypeBitBucket,
				newProject.Type,
			)
		}
		if newProject.URL != "bitbucket.com/example/bar" {
			t.Errorf(
				"New ProjectConfig has invalid type. Expected '%s' received '%s'",
				pkg.ProjectTypeBitBucket,
				newProject.URL,
			)
		}
	})

	t.Run("should return an error if the provided project type is invalid", func(t *testing.T) {
		cmd := project.NewSetProjectCmd(NewMockInputReader('y'))
		cmd.SetArgs([]string{"baz", "invalid_type", "bitbucket.com/example/bar"})
		err := cmd.Execute()
		if err == nil {
			t.Error("Command should have returned an error, instead it resolved")
		}
	})

	t.Run("should return an error if the provided arguments are less than required", func(t *testing.T) {
		testCall := func(args []string) {
			cmd := project.NewSetProjectCmd(NewMockInputReader('y'))
			cmd.SetArgs([]string{"foo"})
			err := cmd.Execute()
			if err == nil {
				t.Error("Command should have returned an error, instead it resolved")
			}
		}
		testArgs := [][]string{
			{"foo"},
			{"foo", "bar"},
			{"foo", "bar", "zaz", "baz"},
			{"foo", "bar", "zaz", "baz", "far"},
		}

		for _, args := range testArgs {
			testCall(args)
		}
	})
}
