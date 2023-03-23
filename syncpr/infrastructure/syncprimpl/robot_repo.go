package syncprimpl

import (
	"fmt"

	"github.com/opensourceways/go-gitee/gitee"
	sdk "github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
)

type giteeClient interface {
	GetRepo(org, repo string) (gitee.Project, error)

	ForkRepo(org, repo string) (gitee.Project, error)

	// GetPullRequests
	GetPullRequests(
		org, repo string, opts sdk.ListPullRequestOpt,
	) ([]gitee.PullRequest, error)

	// CreatePullRequest
	CreatePullRequest(
		org, repo, title, body, head, base string, canModify bool,
	) (gitee.PullRequest, error)
}

func newRobotRepo(cfg *Config) robotRepo {
	cli := sdk.NewClient(func() []byte {
		return []byte(cfg.Robot.Credential.Token)
	})

	return robotRepo{
		repoCache: newRepoCache(),
		cli:       cli,
		robot:     cfg.Robot.Credential.UserName,
		gitURL:    cfg.Robot.remoteURL(),
	}
}

type robotRepo struct {
	repoCache *repoCache
	cli       giteeClient
	robot     string
	gitURL    string
}

func (h *robotRepo) remoteURL(repo string) string {
	return h.gitURL + repo
}

func (h *robotRepo) headBranch(localBranch string) string {
	return fmt.Sprintf("%s:%s", h.robot, localBranch)
}

func (h *robotRepo) hasCreatedPR(pr *domain.PullRequest, localBranch string) (has bool, err error) {
	opt := sdk.ListPullRequestOpt{
		State: "open",
		Head:  h.headBranch(localBranch),
		Base:  pr.Base,
	}

	prs, err := h.cli.GetPullRequests(pr.Org, pr.Repo, opt)
	if err != nil {
		return
	}

	switch len(prs) {
	case 0:
		return
	case 1:
		return true, nil
	}

	return false, fmt.Errorf(
		"There are more than one prs in repo(%s/%s) which are open and created by %s",
		pr.Org, pr.Repo, opt.Head,
	)
}

func (h *robotRepo) createPR(pr *domain.PullRequest, localBranch string) error {
	title := ""

	_, err := h.cli.CreatePullRequest(
		pr.Org, pr.Repo, title, pr.Body,
		h.headBranch(localBranch),
		pr.Base, false,
	)

	return err
}

func (h *robotRepo) tryFork(pr *domain.PullRequest) error {
	if h.hasRepo(pr) {
		return nil
	}

	_, err := h.cli.ForkRepo(pr.Org, pr.Repo)
	if err == nil {
		h.repoCache.addRepo(pr.Repo)
	} else {
		if h.hasRepo(pr) {
			err = nil
		}
	}

	return err
}

func (h *robotRepo) hasRepo(pr *domain.PullRequest) bool {
	if h.repoCache.hasRepo(pr.Repo) {
		return true
	}

	if _, err := h.cli.GetRepo(h.robot, pr.Repo); err != nil {
		return false
	}

	h.repoCache.addRepo(pr.Repo)

	return true
}
