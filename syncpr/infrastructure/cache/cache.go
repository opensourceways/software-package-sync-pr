package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func (s *service) Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := s.cli.SetNX(ctx, key, value, expiration).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	return result, nil
}

func (s *service) UnLock(ctx context.Context, keys []string, args ...interface{}) (interface{}, error) {
	cl := redis.NewScript(
		`
if 
redis.call("GET", KEYS[1]) == ARGV[1] 
then 
return redis.call("DEL", KEYS[1]) 
else 
return 0 
end
	`,
	)
	result, err := cl.Eval(ctx, s.cli, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
