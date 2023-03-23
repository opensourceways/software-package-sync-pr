package synclockimpl

import (
	"context"
	"time"
)

type lock interface {
	Lock(context.Context, string, interface{}, time.Duration) (bool, error)
	UnLock(context.Context, []string, ...interface{}) (interface{}, error)
}
