package pkg

import (
	"fmt"
)

type GroupService interface {
	GetGroup(name string) *GroupConfig
	CreateGroup(name string) (*GroupConfig, error)
	DeleteGroup(name string)
	HasProject(group *GroupConfig, projectName string) bool
	AddProject(group *GroupConfig, projectName string) (*GroupConfig, error)
	RemoveProject(group *GroupConfig, projectName string) *GroupConfig
}

type Group struct {
	Config *WildFireConfig
}

func NewGroupService(config *WildFireConfig) GroupService {
	return &Group{config}
}

func (g *Group) GetGroup(name string) *GroupConfig {
	return g.Config.Groups[name]
}

func (g *Group) CreateGroup(name string) (*GroupConfig, error) {
	if _, ok := g.Config.Groups[name]; ok == true {
		return nil, fmt.Errorf("group with name '%s' already exists", name)
	}

	g.Config.Groups[name] = &GroupConfig{}

	return g.Config.Groups[name], nil
}

func (g *Group) DeleteGroup(name string) {
	delete(g.Config.Groups, name)
}

func (g *Group) HasProject(group *GroupConfig, projectName string) bool {
	for _, project := range *group {
		if project == projectName {
			return true
		}
	}

	return false
}

func (g *Group) AddProject(group *GroupConfig, projectName string) (*GroupConfig, error) {
	_, ok := g.Config.Projects[projectName]
	if ok == false {
		return nil, fmt.Errorf("project with name '%s' does not exist", projectName)
	}

	grp := *group
	grp = append(grp, projectName)
	*group = grp

	return group, nil
}

func (g *Group) RemoveProject(group *GroupConfig, projectName string) *GroupConfig {

	for index, groupProjectName := range *group {
		if groupProjectName == projectName {
			grp := *group
			grp = append(grp[:index], grp[index+1:]...)
			*group = grp

			return group
		}
	}

	return group
}
