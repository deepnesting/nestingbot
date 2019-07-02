package offers

import (
	"time"
)

const (
	SearchNest      = "searchnest"
	SearchNester    = "searchnester"
	SearchCompanion = "searchcompanion"
	SearchHookUp    = "searchhookup"
)

type Offer struct {
	ID        int    `storm:"id,increment"`
	Text      string `validate:"required"`
	Images    []string
	Contacts  string `validate:"required"`
	City      int
	Type      string //searchnest searchnesting
	CreatedAt time.Time
	Paid      bool
	Published bool
	UserID    int
	Tags      []string
}
