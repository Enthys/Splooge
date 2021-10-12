package pkg

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type WildFireConfig struct {
	Projects map[string]Project `yaml:"projects"`
	Groups map[string]ProjectGroup `yaml:"groups"`
}

func GetConfig() *WildFireConfig {
	var config WildFireConfig
	err := viper.Unmarshal(&config)

	if err != nil {
		return &WildFireConfig{
			Projects: make(map[string]Project),
			Groups: make(map[string]ProjectGroup),
		}
	}

	if len(config.Projects) == 0 {
		config.Projects = make(map[string]Project)
	}

	if len(config.Groups) == 0 {
		config.Groups = make(map[string]ProjectGroup)
	}

	return &config
}

func (config *WildFireConfig) SaveConfig() error {
	viper.Set("projects", config.Projects)
	viper.Set("groups", config.Groups)

	return viper.WriteConfig()
}

func (config *WildFireConfig) AddProject(project *Project) error {
	if _, ok := config.Projects[project.Name]; ok != false {
		return errors.New(fmt.Sprintf("Project with name `%s` already exists", project.Name))
	}

	config.Projects[project.Name] = *project

	return nil
}

func (config *WildFireConfig) GetProject(projectName string) *Project {
	if _, ok := config.Projects[projectName]; ok == false {
		return nil
	}

	project := config.Projects[projectName]
	return &project
}

func (config *WildFireConfig) RemoveProject(projectName string) {
	delete(config.Projects, projectName)
}

func (config *WildFireConfig) SetProject(project *Project) {
	config.Projects[project.Name] = *project
}

func (config *WildFireConfig) HasProject(projectName string) bool {
	_, ok := config.Projects[projectName];

	return ok
}
