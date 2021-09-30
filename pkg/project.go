package pkg

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
}
