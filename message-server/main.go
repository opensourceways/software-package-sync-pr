package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-sync-pr/message-server/kafka"
	"github.com/opensourceways/software-package-sync-pr/syncpr/app"
	"github.com/opensourceways/software-package-sync-pr/syncpr/infrastructure/syncprimpl"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	_ = fs.Parse(args)

	return o
}

const component = "software-pkg-sync-pr"

func main() {
	logrusutil.ComponentInit(component)
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		log.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// load config
	cfg, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		log.Errorf("Error loading config, err:%v", err)

		return
	}

	// init kafka
	if err := kafka.Init(&cfg.Kafka, log); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())

		return
	}

	defer kafka.Exit()

	// sync service
	service := app.NewSyncService(
		nil,
		syncprimpl.NewSyncPR(&cfg.SyncPR),
	)

	server := server{
		service:   service,
		userAgent: cfg.Subscription.UserAgent,
	}

	// run
	run(&cfg, &server)
}

func run(cfg *configuration, s *server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")
			return

		case <-sig:
			logrus.Info("receive exit signal")
			called = true
			done()
			return
		}
	}(ctx)

	if err := s.run(&cfg.Subscription, ctx); err != nil {
		logrus.Errorf("server exited, err:%s", err.Error())
	}
}
