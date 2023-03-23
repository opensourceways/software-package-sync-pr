package cache

import (
	"errors"
	"regexp"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9]\d*$`)

type Config struct {
	Address string `json:"address"  required:"true"`
}

func (cfg *Config) Validate() error {
	if reIpPort.MatchString(cfg.Address) {
		return errors.New("invalid mq address")
	}

	return nil
}
