package clone

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"wildfire/pkg"
	"wildfire/pkg/project_repository"
	"github.com/AlecAivazis/survey/v2"
)

var someProjects bool

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

			groupName := args[0]
			group := groupService.GetGroup(groupName)

			if group == nil {
				return config, false, emoji.Errorf("Group '%s' does not exist in configuration.", groupName)
			}

			var pullPath string

			if len(args) > 1 {
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", args[1], groupName))
			} else {
				currentWD, _ := os.Getwd()
				pullPath = filepath.FromSlash(fmt.Sprintf("%s/%s", currentWD, groupName))
			}

			if someProjects == true {
				projects, err := selectProjectsFromGroup(group)

				if err != nil {
					return nil, false, err
				}

				emoji.Println(":star: Selected: ", strings.Join(projects, ", "))
				*group = projects
			}

			err := cloneProjects(projectService, projectRepoService, pullPath, group)
			fmt.Println(err)
			if err != nil {
				err = os.RemoveAll(filepath.FromSlash(pullPath))

				return config, false, err
			}

			fmt.Println(emoji.Sprintf(":ocean: Projects have been cloned to '%s'", pullPath))

			for {
				var action string
				err = survey.AskOne(&survey.Select{
					Message: "Do you wish to take any further action?",
					Options: []string{
						"Run action on all clones",
						"Run action on specific clones",
						"Clear clones and Exit",
						"Exit",
					},
				}, &action)

				if err != nil {
					return config, false, err
				}

				if action == "Exit" {
					break
				}
			}

			return config, false, nil
		}),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolVarP(&someProjects, "some", "s", false, "Only clone some projects from group")

	return cmd
}

func selectProjectsFromGroup(group *pkg.GroupConfig) ([]string, error) {
	var projects []string
	err := survey.AskOne(&survey.MultiSelect{
		Message: "Select projects:",
		Options: *group,
	}, &projects)

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func cloneProjects(
	projectService pkg.ProjectService,
	projectRepoService project_repository.ProjectRepositoryService,
	pullPath string,
	group *pkg.GroupConfig,
) error {
	errString := ""
	var wg sync.WaitGroup

	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithWidth(50),
	)
	cloningBar := createCloneProjectsBar(p, group)
	wg.Add(len(*group))

	for _, projectName := range *group {
		go func(projectName string) {
			project := projectService.GetProject(projectName)
			err := projectRepoService.PullProject(
				filepath.FromSlash(fmt.Sprintf("%s/%s", pullPath, projectName)),
				project,
			)

			if err != nil {
				errString = fmt.Sprintf("%s\nFailed to clone project '%s'. Error: %s", errString, projectName, err.Error())
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

func createCloneProjectsBar(p *mpb.Progress, group *pkg.GroupConfig) *mpb.Bar {
	barName := "Cloning repositories:"
	return p.Add(
		int64(len(*group)),
		mpb.NewBarFiller(mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding(" ").Rbound("]")),
		mpb.PrependDecorators(
			decor.Name(barName, decor.WC{W: len(barName) + 1, C: decor.DidentRight}),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.BarRemoveOnComplete(),
	)
}
