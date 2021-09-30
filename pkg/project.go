package pkg

type ProjectType string

const (
	ProjectTypeGit       ProjectType = "git"
	ProjectTypeGitLab    ProjectType = "gitlab"
	ProjectTypeBitBucket ProjectType = "bitbucket"
)

type ProjectPath string

type Project struct {
	Name string
	Type ProjectType
	URL  ProjectPath
}
