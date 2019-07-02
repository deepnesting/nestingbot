package offers

import (
	"testing"
	//	. "github.com/smartystreets/goconvey/convey"
)

func TestOfferValidation(t *testing.T) {
	o := Offer{
		Text: "df",
		Type: SearchNester,
	}
	err := o.Validate()
	if err != nil {
		t.Error(err)
	}
}
