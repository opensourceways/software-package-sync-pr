package syncprimpl

import (
	"sync"
)

func newRepoCache() *repoCache {
	return &repoCache{
		repos: make(map[string]struct{}),
	}
}

type repoCache struct {
	lock  sync.RWMutex
	repos map[string]struct{}
}

func (c *repoCache) hasRepo(v string) (b bool) {
	c.lock.RLock()
	_, b = c.repos[v]
	c.lock.RUnlock()

	return
}

func (c *repoCache) addRepo(v string) {
	c.lock.Lock()
	c.repos[v] = struct{}{}
	c.lock.Unlock()

	return
}
