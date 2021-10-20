package project_repository

import (
	"github.com/go-git/go-git/v5"
	"io"
	"wildfire/pkg"
)

type Cloner interface {
	CloneProject(path string, project *pkg.ProjectConfig) error
}

type GitCloner struct {
	Output io.Writer
}

func NewCloner(output io.Writer) Cloner {
	return &GitCloner{
		Output: output,
	}
}

func (g *GitCloner) CloneProject(path string, project *pkg.ProjectConfig) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      string(project.URL),
		Progress: g.Output,
	})

	return err
}
