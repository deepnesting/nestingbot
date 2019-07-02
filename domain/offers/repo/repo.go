package repo

import (
	"time"

	"github.com/asdine/storm"
	"github.com/deepnesting/nestingbot/domain/offers"
	"github.com/deepnesting/nestingbot/domain/user"
	"github.com/zhuharev/errors"
)

type repo struct {
	db       *storm.DB
	userRepo user.Repo
}

// New returns repository
func New(db *storm.DB, userRepo user.Repo) offers.Repository {
	return &repo{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *repo) Create(userID int, typ, text string, images []string, contacts string) (id int, err error) {
	user, err := r.userRepo.Get(userID)
	if err != nil {
		return 0, err
	}
	offer := offers.Offer{
		UserID:    userID,
		Text:      text,
		Images:    images,
		Contacts:  contacts,
		City:      user.City,
		CreatedAt: time.Now(),
		Type:      typ,
	}
	err = r.db.Save(&offer)
	return offer.ID, err
}

func (r *repo) Update(o *offers.Offer) error {
	return r.db.Update(o)
}

func (r *repo) GetByID(id int) (*offers.Offer, error) {
	o := offers.Offer{}
	err := r.db.One("ID", id, &o)
	return &o, err
}

func (r *repo) GetByUserID(userID int) (offrs []offers.Offer, err error) {
	err = r.db.Find("UserID", userID, &offrs, storm.Reverse(), storm.Limit(3))
	return
}

func (r *repo) ToggleTag(offerID int, tag string) error {
	o := offers.Offer{}
	err := r.db.One("ID", offerID, &o)
	if err != nil {
		return TagNotFound.New("get tag for offer from db").
			Int("offer_id", offerID).
			String("err", err.Error())
	}

	var found bool
	for _, t := range o.Tags {
		if t == tag {
			found = true
		}
	}

	if !found {
		o.Tags = append(o.Tags, tag)
		err = r.db.Update(&o)
		if err != nil {
			return TagUpdate.New("update offer").
				Int("offer_id", o.ID)
		}
		return nil
	}

	var newTags []string
	for _, t := range o.Tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	o.Tags = newTags
	return r.db.Update(&o)
}

// ErrorType copy type
type ErrorType = errors.ErrorType

const (
	// TagNotFound not found
	TagNotFound ErrorType = ErrorType(iota)
	// TagUpdate update tag in db
	TagUpdate
)
