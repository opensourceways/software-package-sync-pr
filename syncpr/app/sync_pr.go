package app

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
	"github.com/opensourceways/software-package-sync-pr/syncpr/domain/synclock"
	"github.com/opensourceways/software-package-sync-pr/syncpr/domain/syncpr"
	"github.com/opensourceways/software-package-sync-pr/utils"
)

type CmdToSyncPR = domain.PullRequest

type SyncService interface {
	SyncRepo(*CmdToSyncPR) error
}

func NewSyncService(
	l synclock.RepoSyncLock,
	s syncpr.SyncPR,
) *syncService {
	return &syncService{
		lock:   l,
		syncpr: s,
	}
}

type syncService struct {
	rl     repoLock
	lock   synclock.RepoSyncLock
	syncpr syncpr.SyncPR
}

func (s *syncService) SyncRepo(cmd *CmdToSyncPR) (err error) {
	if k := cmd.Org + cmd.Repo; s.rl.tryLock(k) {
		err = s.doSync(cmd)
		s.rl.unlock(k)
	}

	return
}

func (s *syncService) doSync(cmd *CmdToSyncPR) error {
	if err := s.lock.TryLock(&cmd.PullRequestBasic); err != nil {
		return err
	}

	err := s.syncpr.Sync(cmd)

	s.unlock(&cmd.PullRequestBasic)

	return err
}

func (s *syncService) unlock(pr *domain.PullRequestBasic) {
	err := utils.Retry(func() error {
		return s.lock.Unlock(pr)
	})

	if err == nil {
		return
	}

	logrus.Errorf(
		"unlock repo(%s) failed, dead lock happened",
		pr.String(),
	)
}
