package main

import (
	"context"

	"github.com/opensourceways/software-package-sync-pr/message-server/kafka"
	"github.com/opensourceways/software-package-sync-pr/syncpr/app"
)

type server struct {
	service   app.SyncService
	userAgent string
}

func (s *server) run(cfg *subscription, ctx context.Context) error {
	err := kafka.Subscriber().Subscribe(
		cfg.Group,
		map[string]kafka.Handler{
			cfg.Topic: s.handlePR,
		},
	)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *server) handlePR(data []byte, header map[string]string) error {
	msg := msgToHandlePR{s.userAgent}

	cmd, err := msg.toCmd(data, header)
	if err != nil {
		// it is almost the invalid event. ignore it.
		return nil
	}

	return s.service.SyncPR(&cmd)
}
