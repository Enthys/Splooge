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
	cfgFile := getConfigFilePath("set_project.wildfire.yaml")
	_ = os.Remove(cfgFile)
	_ = setConfig(cfgFile)
	config := pkg.GetConfig()
	_ = config.AddProject(&pkg.Project{"foo", pkg.ProjectTypeGit, "github.com/url"})
	_ = config.AddProject(&pkg.Project{"bar", pkg.ProjectTypeGit, "github.com/url"})
	_ = config.AddProject(&pkg.Project{"zaz", pkg.ProjectTypeGit, "github.com/url"})
	err := config.SaveConfig()
	if err != nil {
		t.Error("Failed to initialize configuration")
	}
	config = pkg.GetConfig()
	if len(config.Projects) != 3 {
		t.Errorf(
			"Failed to initialize configuration projects. Expected to have 3 projects found '%d'", len(config.Projects),
		)
	}

	defer func() {
		err = os.Remove(cfgFile)
		if err != nil {
			t.Errorf("Failed to tear down test environment. Error: %s", err)
		}
	}()

	t.Run("Should replace the project configuration", func(t *testing.T) {
		cmd := project.NewSetProjectCmd(NewMockInputReader('y'))
		cmd.SetArgs([]string{"foo", string(pkg.ProjectTypeGitLab), "github.com/example/new"})
		err = cmd.Execute()
		if err != nil {
			t.Errorf("Failed to Set Project. Command should have resolved. Error: %s", err)
		}

		config := pkg.GetConfig()
		updatedPackage := config.GetProject("foo")
		if updatedPackage.Type != pkg.ProjectTypeGitLab {
			t.Errorf(
				"Project type was not updated. Expected '%s' received '%s'",
				pkg.ProjectTypeGitLab,
				updatedPackage.Type,
			)
		}
		if updatedPackage.URL != "github.com/example/new" {
			t.Errorf(
				"Project URL was not updated. Expected '%s' received '%s'",
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
		newProject := config.GetProject("baz")
		if newProject.Type != pkg.ProjectTypeBitBucket {
			t.Errorf(
				"New Project has invalid type. Expected '%s' received '%s'",
				pkg.ProjectTypeBitBucket,
				newProject.Type,
			)
		}
		if newProject.URL != "bitbucket.com/example/bar" {
			t.Errorf(
				"New Project has invalid type. Expected '%s' received '%s'",
				pkg.ProjectTypeBitBucket,
				newProject.URL,
			)
		}
	})
}
