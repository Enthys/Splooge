package it_test

import (
	"os"
	"testing"
	"wildfire/cmd/group"
	"wildfire/pkg"
)

func TestGroupDelete(t *testing.T) {
	cfgFile := getConfigFilePath("group_delete.wildfire.yaml")
	err := initiateConfiguration(cfgFile)
	if err != nil {
		t.Errorf("Failed to initiate configuration. Error: %s", err)
	}

	config := pkg.GetConfig()
	groupService := pkg.NewGroupService(config)
	_, _ = groupService.CreateGroup("foo")
	err = config.SaveConfig()
	if err != nil {
		t.Errorf("Failed to initialize test group. Error: %s", err)
	}

	defer func() {
		err = os.Remove(cfgFile)
		if err != nil {
			t.Errorf("Failed to tear down test environment. Error: %s", err)
		}
	}()

	t.Run("Should throw an error if not enough arguments are passed", func(t *testing.T) {
		cmd := group.NewDeleteGroupCmd()
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected DeleteGroupCMD to return an error instead of resolving")
		}

		expectedErrString := "invalid number of arguments provided"
		if err.Error() != expectedErrString {
			t.Errorf("Invalid error message returned. Expected: %s\nReceived: %s", expectedErrString, err)
		}
	})

	t.Run("Should remove the group from the configuration if the group exists", func(t *testing.T) {
		cmd := group.NewDeleteGroupCmd()
		cmd.SetArgs([]string{"foo"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("DeleteGroupCmd should not have returned an error. Error: %s", err)
		}

		config := pkg.GetConfig()
		if len(config.Groups) != 0 {
			t.Errorf("Invalid number of groups found in configuration. Expected '%d' received '%d'", 0, len(config.Groups))
		}
	})

	t.Run("Should return an error if the group does not exist in configuration", func(t *testing.T) {
		cmd := group.NewDeleteGroupCmd()
		cmd.SetArgs([]string{"foo"})
		err := cmd.Execute()
		if err == nil {
			t.Errorf("DeleteGroupCmd should have returned an error instead of resolving")
		}

		expectedErrString := "Group 'foo' does not exist in configuration."
		if err.Error() != expectedErrString {
			t.Errorf("Invalid error returned. \nExpected Error message '%s'\nReceived '%s'", expectedErrString, err)
		}
	})
}
