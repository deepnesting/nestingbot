package user

import (
	"github.com/asdine/storm"
)

type User struct {
	TID      int `storm:"id"`
	City     int // 1 - spb 2 - msk
	Username string
}

type Repo interface {
	Get(id int) (*User, error)
	Create(id int, username string, city int) error
	Update(*User) error
}

type repo struct {
	db *storm.DB
}

func NewRepository(db *storm.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(id int) (*User, error) {
	u := new(User)
	err := r.db.One("TID", id, u)
	return u, err
}

func (r *repo) Create(id int, username string, city int) error {
	u := User{
		TID:      id,
		City:     city,
		Username: username,
	}
	return r.db.Save(&u)
}

func (r *repo) Update(u *User) error {
	return r.db.Update(u)
}
