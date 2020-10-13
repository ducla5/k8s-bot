package k8s


type Job struct {
	Action     string
	Source     string
	ImageName  string
	BranchName string
	EnvName    string
	Pic        string
	Date       int64
	Ops        map[string]string
}

type Jobs []Job
