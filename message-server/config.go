package main

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-sync-pr/message-server/kafka"
	"github.com/opensourceways/software-package-sync-pr/syncpr/infrastructure/syncprimpl"
)

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type configuration struct {
	Kafka        kafka.Config      `json:"kafka"         required:"true"`
	SyncPR       syncprimpl.Config `json:"sync_pr"       required:"true"`
	Subscription subscription      `json:"subscription"  required:"true"`
}

func (cfg *configuration) configItems() []interface{} {
	return []interface{}{
		&cfg.Kafka,
		&cfg.SyncPR,
	}
}

func (cfg *configuration) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *configuration) setDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.SetDefault()
		}
	}
}

// subscription
type subscription struct {
	Group     string `json:"group"       required:"true"`
	Topic     string `json:"topic"       required:"true"`
	UserAgent string `json:"user_agent"  required:"true"`
}

func loadConfig(file string) (cfg configuration, err error) {
	if err = utils.LoadFromYaml(file, &cfg); err != nil {
		return
	}

	cfg.setDefault()

	err = cfg.validate()

	return
}
