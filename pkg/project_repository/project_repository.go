package project_repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"wildfire/pkg"
)

type ProjectRepositoryService interface {
	PullProject(path string, project *pkg.ProjectConfig) error
	PullGroup(path string, group *pkg.GroupConfig) error
	PullProjectsFromGroup(path string, group *pkg.GroupConfig, projectName ...string) error
}

type ProjectRepository struct {
	projectService pkg.ProjectService
	groupService   pkg.GroupService
	cloner Cloner
}

func NewProjectRepositoryService(
	projectService *pkg.ProjectService,
	groupService *pkg.GroupService,
	cloner Cloner,
) ProjectRepositoryService {
	return &ProjectRepository{
		projectService: *projectService,
		groupService:   *groupService,
		cloner:         cloner,
	}
}

func (p *ProjectRepository) PullProject(path string, project *pkg.ProjectConfig) error {
	err := p.cloner.CloneProject(path, project)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) PullGroup(path string, group *pkg.GroupConfig) error {
	wg := &sync.WaitGroup{}

	wg.Add(len(*group))

	pullErrorsLock := &sync.Mutex{}
	var pullErrors []error

	addError := func(err error) {
		pullErrorsLock.Lock()
		pullErrors = append(pullErrors, err)
		pullErrorsLock.Unlock()
	}

	for _, projectName := range *group {
		go func(projectName string) {
			defer wg.Done()

			project := p.projectService.GetProject(projectName)
			if project == nil {
				addError(fmt.Errorf("project '%s' does not exist in configuration", projectName))
				return
			}

			err := p.PullProject(path, project)
			if err != nil {
				addError(err)
			}
		}(projectName)
	}

	wg.Wait()

	if len(pullErrors) != 0 {
		errorString := ""
		for _, err := range pullErrors {
			errorString = fmt.Sprintf("%s\n%s", errorString, err.Error())
		}

		return errors.New(strings.Trim(errorString, "\n"))
	}

	return nil
}

func (p *ProjectRepository) PullProjectsFromGroup(path string, group *pkg.GroupConfig, projectNames ...string) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(projectNames))

	pullErrorsLock := &sync.Mutex{}
	var pullErrors []error
	addError := func(err error) {
		pullErrorsLock.Lock()
		pullErrors = append(pullErrors, err)
		pullErrorsLock.Unlock()
	}

	for _, projectName := range projectNames {
		go func(projectName string) {
			defer wg.Done()

			if p.groupService.HasProject(group, projectName) == false {
				addError(fmt.Errorf("group does not contain project '%s'", projectName))
				return
			}

			project := p.projectService.GetProject(projectName)
			if project == nil {
				addError(fmt.Errorf("project '%s' does not exist", projectName))
				return
			}

			err := p.PullProject(path, project)
			if err != nil {
				addError(err)
			}
		}(projectName)
	}

	wg.Wait()

	if len(pullErrors) != 0 {
		errorString := ""
		for _, err := range pullErrors {
			errorString = fmt.Sprintf("%s\n%s", errorString, err.Error())
		}

		return errors.New(strings.Trim(errorString, "\n"))
	}

	return nil
}
