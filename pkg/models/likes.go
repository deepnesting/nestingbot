package models

import "github.com/jinzhu/gorm"

//go:generate goqueryset -in likes.go

// User is main struct
// easyjson:json
// gen:qs
type UserVote struct {
	UserID   uint
	ObjectID int64
	Vote     int
}

func VotesVote(userID uint, objectID int64, vote int) (err error) {
	var uv UserVote
	qs := NewUserVoteQuerySet(gormDB).UserIDEq(userID).ObjectIDEq(objectID)
	err = qs.One(&uv)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uv.ObjectID = objectID
			uv.UserID = userID
			uv.Vote = vote
			addr := &uv
			return addr.Create(gormDB)
		}
	}
	return qs.GetUpdater().SetVote(vote).Update()
}

func VotesCount(objectID int64, vote int) (int, error) {
	return NewUserVoteQuerySet(gormDB).ObjectIDEq(objectID).VoteEq(vote).Count()
}
