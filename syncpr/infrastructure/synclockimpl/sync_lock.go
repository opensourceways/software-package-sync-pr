package synclockimpl

import (
	"errors"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
)

func NewRepoSyncLock(cli dbClient) syncLock {
	return syncLock{cli: cli}
}

type syncLock struct {
	cli dbClient
}

func (impl syncLock) TryLock(i *domain.PullRequestBasic) (err error) {
	var do PRInfoDO
	toPRInfoDO(i, &do)

	filter := PRInfoDO{
		Org:  i.Org,
		Repo: i.Repo,
		Num:  i.Num,
	}

	err = impl.cli.FirstOrCreate(&filter, &do)
	if err == nil || !impl.cli.IsRowExists(err) {
		return
	}

	if do.isBusy() {
		err = errors.New("record busy")

		return
	}

	filter.Status = free
	err = impl.cli.UpdateRecord(&filter, &PRInfoDO{Status: busy})

	return
}

func (impl syncLock) Unlock(i *domain.PullRequestBasic) error {
	return impl.cli.DeleteRecord(
		&PRInfoDO{Org: i.Org, Repo: i.Repo, Num: i.Num, Status: busy},
	)
}
