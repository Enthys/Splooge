package pkg

import (
	"github.com/go-git/go-git/v5"
	"os"
)

type ProjectType string

const (
	ProjectTypeGit       ProjectType = "git"
	ProjectTypeGitLab    ProjectType = "gitlab"
	ProjectTypeBitBucket ProjectType = "bitbucket"
)

func (t *ProjectType) ValidType() bool {
	_, ok := t.GetAvailableTypes()[*t]

	return ok
}

func (t ProjectType) GetAvailableTypes() map[ProjectType]ProjectType {
	return map[ProjectType]ProjectType{
		ProjectTypeGit: ProjectTypeGit,
		ProjectTypeGitLab: ProjectTypeGitLab,
		ProjectTypeBitBucket: ProjectTypeBitBucket,
	}
}

type ProjectPath string

type Project struct {
	Name string
	Type ProjectType
	URL  ProjectPath
	DefaultBranch string

}

func (project *Project) PullProject(path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      string(project.URL),
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}
