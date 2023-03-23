package synclockimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/software-package-sync-pr/syncpr/domain"
)

const (
	busy = "busy"
	free = "free"
)

type PRInfoDO struct {
	Id     uuid.UUID `gorm:"column:uuid;type:uuid"`
	Org    string    `gorm:"column:org"`
	Repo   string    `gorm:"column:repo"`
	Num    int       `gorm:"column:num"`
	Status string    `gorm:"column:status"`
}

func toPRInfoDO(i *domain.PullRequestBasic, do *PRInfoDO) {
	*do = PRInfoDO{
		Id:     uuid.New(),
		Org:    i.Org,
		Repo:   i.Repo,
		Num:    i.Num,
		Status: busy,
	}
}

func (r *PRInfoDO) isBusy() bool {
	return r.Status == busy
}
