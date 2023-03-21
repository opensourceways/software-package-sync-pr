package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-github/v50/github"

	"github.com/opensourceways/software-package-sync-pr/syncpr/app"
	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
)

const (
	msgHeaderUUID        = "X-GitHub-Delivery"
	msgHeaderUserAgent   = "User-Agent"
	msgHeaderEventType   = "X-GitHub-Event"
	eventTypePullRequest = "pull_request"
)

type msgToHandlePR struct {
	userAgent string
}

func (msg *msgToHandlePR) toCmd(payload []byte, header map[string]string) (
	cmd app.CmdToSyncPR, err error,
) {
	eventType, err := msg.parseRequest(header)
	if err != nil {
		err = fmt.Errorf("invalid msg, err:%s", err.Error())

		return
	}

	if eventType != eventTypePullRequest {
		err = errors.New("not pull_request event")

		return
	}

	e := new(github.PullRequestEvent)
	if err = json.Unmarshal(payload, e); err != nil {
		return
	}

	if e.GetAction() != "opened" {
		err = errors.New("not opened")

		return
	}

	cmd = msg.genCmd(e)

	return
}

func (msg *msgToHandlePR) genCmd(e *github.PullRequestEvent) (cmd app.CmdToSyncPR) {
	repo := e.GetRepo()
	pr := e.GetPullRequest()

	return app.CmdToSyncPR{
		PullRequestBasic: domain.PullRequestBasic{
			Org:  repo.GetOwner().GetLogin(),
			Repo: repo.GetName(),
			Num:  pr.GetNumber(),
		},
		Body:     pr.GetBody(),
		Base:     pr.GetBase().GetRef(),
		RepoLink: repo.GetURL(),
	}
}

func (msg *msgToHandlePR) parseRequest(header map[string]string) (
	eventType string, err error,
) {
	if header == nil {
		err = errors.New("no header")

		return
	}

	if header[msgHeaderUserAgent] != msg.userAgent {
		err = errors.New("unknown " + msgHeaderUserAgent)

		return
	}

	if eventType = header[msgHeaderEventType]; eventType == "" {
		err = errors.New("missing " + msgHeaderEventType)

		return
	}

	if header[msgHeaderUUID] == "" {
		err = errors.New("missing " + msgHeaderUUID)
	}

	return
}
