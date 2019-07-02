package offers

//"github.com/asdine/storm"
//	"github.com/deepnesting/nestingbot/domain/offers/repo"

type Repository interface {
	Create(userID int, typ string, text string, images []string, contacts string) (id int, err error)
	GetByID(id int) (*Offer, error)
	GetByUserID(userID int) ([]Offer, error)
	ToggleTag(offerID int, tag string) error
	Update(o *Offer) error

	List() ([]Offer, error)
}

// func NewRepository(db *storm.DB) Repository {
// 	return repo.New(db)
// }
