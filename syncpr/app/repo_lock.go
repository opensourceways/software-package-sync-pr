package app

import "sync"

const (
	repoBusy = "busy"
	repoFree = "free"
)

type repoLock struct {
	lock  sync.RWMutex
	repos map[string]string
}

func (rl *repoLock) tryLock(repo string) (locked bool) {
	if !rl.canLock(repo) {
		return
	}

	// try lock
	rl.lock.Lock()

	if !rl.isLocked(repo) {
		rl.repos[repo] = repoBusy
		locked = true
	}

	rl.lock.Unlock()

	return

}

func (rl *repoLock) canLock(repo string) (r bool) {
	rl.lock.RLock()
	r = !rl.isLocked(repo)
	rl.lock.RUnlock()

	return
}

func (rl *repoLock) isLocked(repo string) bool {
	s, ok := rl.repos[repo]
	return ok && rl.isBusy(s)
}

func (rl *repoLock) unlock(repo string) {
	rl.lock.Lock()
	rl.repos[repo] = repoFree
	rl.lock.Unlock()
}

func (rl *repoLock) isBusy(s string) bool {
	return s == repoBusy
}
