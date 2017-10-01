package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zhuharev/boltutils"
)

var (
	db *boltutils.DB

	gormDB *gorm.DB

	boltPrefix  = "nst_"
	usersBucket = []byte(boltPrefix + "users")
)

func NewContext() (err error) {
	db, err = boltutils.Open("data/db.bolt", 0777, nil)
	if err != nil {
		return
	}
	db.EnableGzip = true
	err = db.CreateBucket(usersBucket)

	// sql

	gormDB, err = gorm.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		return err
	}
	gormDB.LogMode(false)

	gormDB.AutoMigrate(&User{}, &Subscription{}, &UserVote{})

	return
}
