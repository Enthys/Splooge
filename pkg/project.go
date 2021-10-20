package pkg

import (
	"errors"
	"fmt"
)

type ProjectService interface {
	AddProject(name string, url ProjectPath, projectType ProjectType) (*ProjectConfig, error)
	RemoveProject(name string)
	HasProject(name string) bool
	GetProject(name string) *ProjectConfig
	UpdateOrCreate(project *ProjectConfig)
}

type Project struct {
	Config *WildFireConfig
}

func NewProjectService(config *WildFireConfig) ProjectService {
	return &Project{config}
}

func (p *Project) AddProject(name string, url ProjectPath, projectType ProjectType) (*ProjectConfig, error) {
	if _, ok := p.Config.Projects[name]; ok != false {
		return nil, errors.New(fmt.Sprintf("ProjectConfig with name `%s` already exists", name))
	}

	p.Config.Projects[name] = &ProjectConfig{
		Name: name,
		Type: projectType,
		URL:  url,
	}

	return p.Config.Projects[name], nil
}

func (p *Project) RemoveProject(name string) {
	delete(p.Config.Projects, name)
}

func (p *Project) HasProject(name string) bool {
	_, ok := p.Config.Projects[name]

	return ok
}

func (p *Project) GetProject(name string) *ProjectConfig {
	if _, ok := p.Config.Projects[name]; ok == false {
		return nil
	}

	return p.Config.Projects[name]
}

func (p *Project) UpdateOrCreate(project *ProjectConfig) {
	p.Config.Projects[project.Name] = project
}
