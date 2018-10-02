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
	db, err = boltutils.New(boltutils.OpenPath("data/db.bolt"), boltutils.Compression(boltutils.GzipCompressor))
	if err != nil {
		return
	}
	err = db.CreateBucket(usersBucket)
	if err != nil {
		return err
	}
	// sql

	gormDB, err = gorm.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		return err
	}
	gormDB.LogMode(false)

	gormDB.AutoMigrate(&User{}, &Subscription{}, &UserVote{})

	return
}
