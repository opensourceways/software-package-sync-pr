package syncprimpl

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
	"github.com/opensourceways/software-package-sync-pr/utils"
)

func NewSyncPR(cfg *Config) *syncPR {
	return &syncPR{
		shell:     cfg.SyncRepoShell,
		workDir:   cfg.WorkDir,
		robotRepo: newRobotRepo(cfg),
	}
}

type syncPR struct {
	shell     string
	workDir   string
	robotRepo robotRepo
}

func (impl *syncPR) Sync(pr *domain.PullRequest) error {
	if err := impl.syncPRBranch(pr); err != nil {
		return err
	}

	branch := impl.localBranch(pr)

	if b, err := impl.robotRepo.hasCreatedPR(pr, branch); err != nil || b {
		return err
	}

	if err := impl.robotRepo.tryFork(pr); err != nil {
		return err
	}

	return impl.robotRepo.createPR(pr, branch)
}

func (impl *syncPR) syncPRBranch(pr *domain.PullRequest) error {
	params := []string{
		impl.shell,
		impl.workDir,
		strconv.Itoa(pr.Num), pr.CloneURL,
		impl.robotRepo.remoteURL(pr.Repo),
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

func (impl *syncPR) localBranch(pr *domain.PullRequest) string {
	return fmt.Sprintf("pull%d", pr.Num)
}
