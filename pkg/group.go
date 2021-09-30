package pkg

import (
	"errors"
	"fmt"
)

type ProjectGroup []string

func (group *ProjectGroup) AddProject(config *SploogeConfig, projectName string) error {
	if _, ok := config.Projects[projectName]; ok == false {
		return errors.New(fmt.Sprintf("Project with name `%s` does not exist", projectName))
	}

	group = ProjectGroup(append([]string(*group), projectName))
}
