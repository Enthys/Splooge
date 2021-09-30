package pkg

import (
	"errors"
	"fmt"
)

type ProjectGroup []string

func (group *ProjectGroup) AddProject(config *WildFireConfig, projectName string) (ProjectGroup, error) {
	if _, ok := config.Projects[projectName]; ok == false {
		return nil, errors.New(fmt.Sprintf("Project with name `%s` does not exist", projectName))
	}

	if group.HasProject(&projectName) {
		return *group, nil
	}

	return append(*group, projectName), nil
}

func (group *ProjectGroup) HasProject(projectName *string) bool {
	for _, project := range *group {
		if project == *projectName {
			return true
		}
	}

	return false
}

func (group *ProjectGroup) RemoveProject(projectName string) ProjectGroup {
	g := *group
	for index, groupProjectName := range g {
		if groupProjectName == projectName {
			return append(g[:index], g[index+1:]...)
		}
	}

	return *group
}
