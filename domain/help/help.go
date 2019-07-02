package help

import (
	"time"

	"github.com/asdine/storm"
)

type Help struct {
	ID         int `storm:"id,increment"`
	UserID     int
	Text       string
	IsIncoming bool
	CreatedAt  time.Time
}

type Repo interface {
	Create(userID int, text string, isIncoming bool) error
	List(userID int) ([]Help, error)
}

type repo struct {
	db *storm.DB
}

func NewRepo(db *storm.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(userID int, text string, isIncoming bool) error {
	h := Help{
		UserID:     userID,
		Text:       text,
		IsIncoming: isIncoming,
		CreatedAt:  time.Now(),
	}
	return r.db.Save(&h)
}

func (r *repo) List(userID int) (messages []Help, err error) {
	err = r.db.Find("UserID", userID, &messages, storm.Limit(10), storm.Reverse())
	return
}
