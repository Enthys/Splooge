package it_test

import (
	"os"
	"testing"
	"wildfire/cmd/group"
	"wildfire/pkg"
)

func TestGroupCreate(t *testing.T) {
	cfgFile := getConfigFilePath("group_create.wildfire.yaml")
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

	t.Run("should create empty group if no projects args were provided", func(t *testing.T) {
		cmd := group.NewCreateGroupCmd()
		cmd.SetArgs([]string{"foo"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Create group command should not have returned an error. Error: %s", err)
		}

		config := pkg.GetConfig()
		if len(config.Groups) != 1 {
			t.Errorf(
				"Invalid number of groups found in configuration. Expected '%d' found '%d'",
				1,
				len(config.Groups),
			)
		}
		if _, ok := config.Groups["foo"]; ok == false {
			t.Errorf("Could not find group in configuration by map key")
		}
	})

	t.Run("should add provided by name projects to group if projects exist", func(t *testing.T) {
		cmd := group.NewCreateGroupCmd()
		cmd.SetArgs([]string{"foo", "foo", "bar", "zaz"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Command should not have returned an error. Error: %s", err)
		}

		config := pkg.GetConfig()
		if len(config.Groups) != 1 {
			t.Errorf(
				"Invalid number of groups found in configuration. Expected '%d' found '%d'",
				1,
				len(config.Groups),
			)
		}

		newGroup := config.Groups["foo"]
		if len(newGroup) != 3 {
			t.Errorf("Invalid number of projects found in group. Expected '%d' found '%d'", 3, len(newGroup))
		}
	})
}
