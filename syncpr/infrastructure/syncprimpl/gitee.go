package syncprimpl

import (
	"fmt"

	"github.com/opensourceways/go-gitee/gitee"
	sdk "github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain/syncpr"
)

type giteeClient interface {
	// GetPullRequests
	GetPullRequests(
		org, repo string, opts sdk.ListPullRequestOpt,
	) ([]gitee.PullRequest, error)

	// CreatePullRequest
	CreatePullRequest(
		org, repo, title, body, head, base string, canModify bool,
	) (gitee.PullRequest, error)
}

type giteePR struct {
	cli   giteeClient
	robot string
}

func (h *giteePR) headBranch(localBranch string) string {
	return fmt.Sprintf("%s:%s", h.robot, localBranch)
}

func (h *giteePR) hasPR(pr *syncpr.PullRequest, localBranch string) (has bool, err error) {
	opt := sdk.ListPullRequestOpt{
		State: "open",
		Head:  h.headBranch(localBranch),
		Base:  pr.TargetBranch,
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

func (h *giteePR) createPR(pr *syncpr.PullRequest, localBranch string) error {
	title := ""

	_, err := h.cli.CreatePullRequest(
		pr.Org, pr.Repo, title, pr.Body,
		h.headBranch(localBranch),
		pr.TargetBranch, false,
	)

	return err
}

func (h *giteePR) tryFork(pr *syncpr.PullRequest) error {
	return nil
}
