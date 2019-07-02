package offers 

import (
	"gopkg.in/go-playground/validator.v9"
)


func (o Offer) Validate() error {
	validate := validator.New()
	validate.RegisterStructValidation(offerStructLevelValidation, Offer{})
	return validate.Struct(o)
}

func offerStructLevelValidation(sl validator.StructLevel) {

	offer := sl.Current().Interface().(Offer)

	// если сообщение о сдаче, то фото обязательно
	if offer.Type == SearchNester {
		if len(offer.Images) == 0 {
			sl.ReportError(offer.Images, "Images", "images", "images", "")
		}
	}

	// if len(user.FirstName) == 0 && len(user.LastName) == 0 {
	// 	sl.ReportError(user.FirstName, "FirstName", "fname", "fnameorlname", "")
	// 	sl.ReportError(user.LastName, "LastName", "lname", "fnameorlname", "")
	// }

	// plus can do more, even with different tag than "fnameorlname"
}
