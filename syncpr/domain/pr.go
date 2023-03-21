package domain

import "fmt"

type PullRequestBasic struct {
	Org  string
	Repo string
	Num  int
}

func (pr *PullRequestBasic) String() string {
	return fmt.Sprintf("%s/%s/%d", pr.Org, pr.Repo, pr.Num)
}

type PullRequest struct {
	PullRequestBasic

	Body     string
	Base     string // Base is the branch to be merged to
	RepoLink string
}
