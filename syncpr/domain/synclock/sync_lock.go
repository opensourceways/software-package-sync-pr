package synclock

import "github.com/opensourceways/software-package-sync-pr/syncpr/domain"

type RepoSyncLock interface {
	TryLock(*domain.PullRequestBasic) error
	Unlock(*domain.PullRequestBasic) error
}
