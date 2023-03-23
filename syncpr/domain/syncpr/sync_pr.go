package syncpr

import "github.com/opensourceways/software-package-sync-pr/syncpr/domain"

type SyncPR interface {
	Sync(*domain.PullRequest) error
}
