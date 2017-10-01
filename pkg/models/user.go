package models

import (
	"github.com/jinzhu/gorm"
	"github.com/zhuharev/intarr"
)

//go:generate goqueryset -in user.go

// User is main struct
// easyjson:json
// gen:qs
type User struct {
	gorm.Model
	TelegramID int64

	FirstName string
	LastName  string

	Username string
}

func (u User) btsID() []byte {
	return intarr.Uint64ToBytes(uint64(u.TelegramID))
}

func UserGetOrCreate(telegramID int64) (*User, error) {
	u := new(User)
	u.TelegramID = telegramID
	err := NewUserQuerySet(gormDB).TelegramIDEq(telegramID).One(u)
	if err == gorm.ErrRecordNotFound {
		err = u.Create(gormDB)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return u, err
}

func UserUpdate(user *User) error {
	return user.Update(gormDB)
}
