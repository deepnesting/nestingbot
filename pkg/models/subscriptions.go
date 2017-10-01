package models

func Subscribe(userID int64, channelID string) error {
	sbs := Subscription{
		UserID: userID, ChannelID: channelID,
	}
	err := sbs.Create(gormDB)
	if err != nil {
		return err
	}
	return nil
}

func Unsibscribe(userID int64, channelID string) error {
	return NewSubscriptionQuerySet(gormDB).ChannelIDEq(channelID).UserIDEq(userID).Delete()
}

func GetSubscriptions(userID int64) (res []string, err error) {
	var subs []Subscription
	err = NewSubscriptionQuerySet(gormDB).UserIDEq(userID).All(&subs)
	if err != nil {
		return
	}
	for _, v := range subs {
		res = append(res, v.ChannelID)
	}
	return
}

func GetSubscribers(channelID string) (res []int64, err error) {
	var subs []Subscription
	err = NewSubscriptionQuerySet(gormDB).ChannelIDEq(channelID).All(&subs)
	if err != nil {
		return
	}
	for _, v := range subs {
		res = append(res, v.UserID)
	}
	return
}

//go:generate goqueryset -in subscriptions.go

// Subscription is relation struct
// gen:qs
type Subscription struct {
	// telegram id
	UserID    int64  `gorm:"unique_index:idx_subscriptions"`
	ChannelID string `gorm:"unique_index:idx_subscriptions"`
}

func GetSubscribersCount(channelID string) (int, error) {
	return NewSubscriptionQuerySet(gormDB).ChannelIDEq(channelID).Count()
}
