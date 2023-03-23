package synclockimpl

import "github.com/opensourceways/software-package-sync-pr/syncpr/domain"

func NewRepoSyncLock() syncLock {
	return syncLock{}
}

type syncLock struct {
}

func (impl syncLock) TryLock(*domain.PullRequestBasic) error {
	return nil
}

func (impl syncLock) Unlock(*domain.PullRequestBasic) error {
	return nil
}
