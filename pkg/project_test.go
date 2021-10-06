package pkg_test

import (
	"fmt"
	"testing"
	"wildfire/pkg"
)

func TestProjectType_ValidType(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		tests := []struct{
			Project pkg.Project
			Expected bool
		}{
			{ pkg.Project{"git", pkg.ProjectTypeGit, "git.com/foo"}, true },
			{ pkg.Project{"gitlab", pkg.ProjectTypeGitLab, "git.com/foo"}, true },
			{ pkg.Project{"bitbucket", pkg.ProjectTypeBitBucket, "git.com/foo"}, true },
			{ pkg.Project{"fake", pkg.ProjectType("fake"), "git.com/foo"}, false },
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
