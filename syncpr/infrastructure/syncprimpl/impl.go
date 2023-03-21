package syncprimpl

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain/syncpr"
	"github.com/opensourceways/software-package-sync-pr/utils"
)

func NewSyncPR(cfg *Config) *syncPR {
	return &syncPR{
		shell:     cfg.SyncRepoShell,
		workDir:   cfg.WorkDir,
		robotRepo: cfg.RobotRepo.remoteURL(),
	}
}

type syncPR struct {
	shell     string
	workDir   string
	robotRepo string
	giteePR   giteePR
}

func (impl *syncPR) Sync(pr *syncpr.PullRequest) error {
	if err := impl.syncPRBranch(pr); err != nil {
		return err
	}

	branch := impl.localBranch(pr)

	if b, err := impl.giteePR.hasPR(pr, branch); err != nil || b {
		return err
	}

	if err := impl.giteePR.tryFork(pr); err != nil {
		return err
	}

	return impl.giteePR.createPR(pr, branch)
}

func (impl *syncPR) syncPRBranch(pr *syncpr.PullRequest) error {
	params := []string{
		impl.shell,
		impl.workDir,
		strconv.Itoa(pr.Num), pr.RepoLink,
		impl.robotRepo + pr.Repo,
	}

	_, err, _ := utils.RunCmd(params...)
	if err != nil {
		logrus.Errorf(
			"run sync shell, err=%s, params=%v",
			err.Error(), params[:len(params)-1],
		)
	}

	return err
}

func (impl *syncPR) localBranch(pr *syncpr.PullRequest) string {
	return fmt.Sprintf("pull%d", pr.Num)
}
