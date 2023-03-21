package syncpr

type PullRequest struct {
	Org      string
	Repo     string
	Num      int
	Body     string
	Base     string // Base is the branch to be merged to
	RepoLink string
}

type SyncPR interface {
	Sync(*PullRequest) error
}
