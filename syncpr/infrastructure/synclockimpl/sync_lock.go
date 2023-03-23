package synclockimpl

import (
	"context"
	"errors"
	"time"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
	"github.com/opensourceways/software-package-sync-pr/syncpr/infrastructure/cache"
)

func NewRepoSyncLock() syncLock {
	return syncLock{
		l: cache.Instance(),
	}
}

type syncLock struct {
	l lock
}

func (impl syncLock) TryLock(p *domain.PullRequestBasic) error {
	success, err := impl.l.Lock(context.Background(), p.String(), p.String(), time.Minute*30)
	if err != nil {
		return err
	}

	if !success {
		return errors.New("busy")
	}

	return nil
}

func (impl syncLock) Unlock(p *domain.PullRequestBasic) error {
	flag, err := impl.l.UnLock(context.Background(), []string{p.String()}, p.String())
	if err != nil {
		return err
	}

	if flag.(int64) == 0 {
		return errors.New("unlock failed")
	}

	return nil
}
