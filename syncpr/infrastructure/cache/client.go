package cache

import "github.com/go-redis/redis/v8"

var instance *service

type service struct {
	cli *redis.Client
}

func Instance() *service {
	return instance
}

func Init(cfg *Config) error {
	cli := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
	})

	if err := cli.Ping(cli.Context()).Err(); err != nil {
		return err
	}

	instance = &service{cli: cli}

	return nil
}
