package pkg

import (
	"errors"
	"fmt"
)

type ProjectGroup []string

func CreateGroup(config *WildFireConfig, name string) *ProjectGroup {
	group := ProjectGroup{}

	if len(config.Groups) == 0 {
		config.Groups = map[string]ProjectGroup{}
	}

	config.Groups[name] = group

	return &group
}

func RemoveGroup(config *WildFireConfig, groupName string) {
	delete(config.Groups, groupName)
}

func (group *ProjectGroup) AddProject(config *WildFireConfig, projectName string) (*ProjectGroup, error) {
	if _, ok := config.Projects[projectName]; ok == false {
		return nil, errors.New(fmt.Sprintf("Project with name `%s` does not exist", projectName))
	}

	if group.HasProject(&projectName) {
		return group, nil
	}

	newGroup := append(*group, projectName)
	return &newGroup, nil
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
