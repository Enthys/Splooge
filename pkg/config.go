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
		return &WildFireConfig{Projects: map[string]Project{}}
	}

	return &config
}

func (config *WildFireConfig) SaveConfig() error {
	viper.Set("projects", config.Projects)
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
