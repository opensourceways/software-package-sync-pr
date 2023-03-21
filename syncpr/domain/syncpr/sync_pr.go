package syncpr

type PullRequest struct {
	Org      string
	Repo     string
	Num      int
	Body     string
	RepoLink string
	// TargetBranch is the branch to be merged to
	TargetBranch string
}

type SyncPR interface {
	Sync(*PullRequest) error
}
