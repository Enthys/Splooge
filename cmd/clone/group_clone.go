package clone

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"wildfire/pkg"
	"wildfire/pkg/project_repository"
)

var someProjects bool

type UserInput interface {
	PlainInput(msg string) (string, error)
	PickOne(msg string, options []string) (string, error)
	PickMultiple(msg string, options []string) ([]string, error)
}

type SurveyUserInput struct {}

func (s SurveyUserInput) PlainInput(msg string) (string, error) {
	var response string
	err := survey.AskOne(&survey.Input{Message: msg}, &response)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (s SurveyUserInput) PickOne(msg string, options []string) (string, error) {
	var response string
	err := survey.AskOne(
		&survey.Select{Message: msg, Options: options},
		&response,
	)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (s SurveyUserInput) PickMultiple(msg string, options []string) ([]string, error) {
	var response []string
	err := survey.AskOne(
		&survey.MultiSelect{Message: msg, Options: options},
		&response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type pullGroupExecutor struct {
	projectService pkg.ProjectService
	groupService   pkg.GroupService
	repoService    project_repository.ProjectRepositoryService
	userInput      UserInput
}

func (executor *pullGroupExecutor) Execute(groupName string, path string, partialClone bool) error {
	group := executor.groupService.GetGroup(groupName)
	if group == nil {
		return emoji.Errorf("Group '%s' does not exist in configuration.", groupName)
	}

	if partialClone == true {
		projects, err := executor.pickProjectsFromGroup(group)

		if err != nil {
			return err
		}

		_, _ = emoji.Println(":star: Selected: ", strings.Join(projects, ", "))
		*group = projects
	}

	err := executor.cloneGroupProjects(group, path)
	if err != nil {
		clearErr := executor.clearPath(path)

		if clearErr != nil {
			err = fmt.Errorf("%s\n%s", err.Error(), clearErr.Error())
		}

		return err
	}

	fmt.Println(emoji.Sprintf(":ocean: Projects have been cloned to '%s'", path))

	var repoActionScope string

	for repoActionScope != "Exit" {
		action, err := executor.userInput.PickOne(
			"Do you wish to take any further action?",
			[]string{
				"Run command",
				"Clear clones and Exit",
				"Exit",
			},
		)

		if err == terminal.InterruptErr {
			return nil
		}

		if err != nil {
			return err
		}
		if action == "Run command" {
			err = executor.runCommand(*group, path)

			if err != nil && err != terminal.InterruptErr {
				return err
			}
		}
		if action == "Exit" {
			break
		}
	}

	return nil
}

func (executor *pullGroupExecutor) pickProjectsFromGroup(group *pkg.GroupConfig) ([]string, error) {
	projects, err := executor.userInput.PickMultiple("Select projects:", *group)

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (executor *pullGroupExecutor) cloneGroupProjects(
	group *pkg.GroupConfig,
	pullPath string,
) error {
	errString := ""
	var wg sync.WaitGroup

	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(50))
	cloningBar := executor.createBarForGroup("Cloning repositories:", p, group)
	wg.Add(len(*group))

	for _, projectName := range *group {
		go func(projectName string) {
			project := executor.projectService.GetProject(projectName)
			err := executor.repoService.PullProject(
				filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)),
				project,
			)

			if err != nil {
				errString = fmt.Sprintf(
					"%s\nFailed to clone project '%s'. Error: %s",
					errString,
					projectName,
					err.Error(),
				)
			}

			cloningBar.Increment()
			wg.Done()
		}(projectName)
	}
	p.Wait()

	if errString != "" {
		return errors.New(strings.Trim(errString, "\n"))
	}

	return nil
}

func (executor *pullGroupExecutor) clearPath(path string) error {
	return os.RemoveAll(filepath.FromSlash(path))
}

func (executor *pullGroupExecutor) createBarForGroup(name string, p *mpb.Progress, group *pkg.GroupConfig) *mpb.Bar {
	return p.Add(
		int64(len(*group)),
		mpb.NewBarFiller(mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding(" ").Rbound("]")),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.BarRemoveOnComplete(),
	)
}

func (executor *pullGroupExecutor) runCommand(group pkg.GroupConfig, path string) error {
	scope, err := executor.userInput.PickOne("Select scope:", []string{"All", "Select projects"})
	if err != nil {
		return err
	}

	if scope == "Select projects" {
		selectedProjects, err := executor.userInput.PickMultiple("Select clones:", group)
		if err != nil {
			return err
		}

		emoji.Println("Selected: ", strings.Join(selectedProjects, ", "))
		group = selectedProjects
	}

	actionString, err := executor.userInput.PlainInput("Command:")
	if err != nil {
		return err
	}

	action, actionArgs, err := func(actionString string) (string, []string, error) {
		r := csv.NewReader(strings.NewReader(actionString))
		r.Comma = ' ' // space
		parts, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return "", nil, err
		}

		return parts[0], parts[1:], nil
	}(actionString)


	var wg sync.WaitGroup
	results := sync.Map{}
	var errors []error

	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(50))
	executionProgressBar := executor.createBarForGroup("Running command:", p, &group)
	wg.Add(len(group))

	for _, project := range group {
		fmt.Println(project)
		go func(project string) {
			defer wg.Done()

			command := exec.Command(action, actionArgs...)
			command.Dir = filepath.FromSlash(fmt.Sprintf("%s/%s", path, project))
			var buf bytes.Buffer
			command.Stdout = &buf

			err := command.Run()
			if err != nil {
				errors = append(errors, err)
			}

			results.Store(project, buf.String())
			executionProgressBar.Increment()
		}(project)
	}

	p.Wait()

	results.Range(func(projectName, result interface{}) bool {
		fmt.Println(projectName)
		fmt.Println(result)

		return true
	})

	return nil
}

func NewPullGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group <group name> [path]",
		Short: "Pull group projects from their repositories",
		Long:  `Pull the projects stored in the specified group in the current directory or a specified directory.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("invalid number of arguments provided")
			}

			return nil
		},
		RunE: pkg.ProjectFunc(func(config *pkg.WildFireConfig, cmd *cobra.Command, args []string) (*pkg.WildFireConfig, bool, error) {
			projectService := pkg.NewProjectService(config)
			groupService := pkg.NewGroupService(config)
			projectRepoService := project_repository.NewProjectRepositoryService(
				&projectService,
				&groupService,
				project_repository.NewCloner(nil),
			)

			var input SurveyUserInput
			executor := &pullGroupExecutor{
				projectService: projectService,
				groupService:   groupService,
				repoService:    projectRepoService,
				userInput:      input,
			}

			groupName := args[0]

			var pullPath string

			if len(args) > 1 {
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", args[1], groupName))
			} else {
				currentWD, _ := os.Getwd()
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", currentWD, groupName))
			}

			err := executor.Execute(groupName, pullPath, someProjects)

			return config, false, err
		}),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolVarP(&someProjects, "some", "s", false, "Only clone some projects from group")

	return cmd
}
