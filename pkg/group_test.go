package pkg_test

import (
	"splooge/pkg"
	"testing"
)

//fooProject = pkg.Project{"foo", pkg.ProjectTypeGit, "https://github.com/foo"}
//config = &pkg.SploogeConfig{Projects: map[string]pkg.Project{ }}

func TestProjectGroup_AddProject_should_add_the_name_of_project_to_its_collection(t *testing.T) {
	config := &pkg.SploogeConfig{Projects: map[string]pkg.Project{
		"foo": pkg.Project{"foo", pkg.ProjectTypeGit, "https://github.com/foo"},
	}}
	group := pkg.ProjectGroup{}
	group.AddProject(config, "foo")

	if len(group) == 0 {
		t.Error("Project was not added to ProjectGroup collection")
	}
}
