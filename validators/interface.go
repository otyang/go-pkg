package validators

import "errors"

type IValidators interface {
	// ValidateStruct used to validate a struct pointer
	ValidateStruct(vPtr any) error
	// Used to transate the errors from VaalidateStruct to error
	Translator(err error) string
}

// avoids repetition
func HelperValidateStructGoKit(payloadPtr any) error {
	o := NewGooKitValidator()
	err := o.ValidateStruct(payloadPtr)
	if err != nil {
		msg := o.Translator(err)
		return errors.New(msg)
	}
	return nil
}

// avoids repetition
func HelperValidateStructGoPlayground(payloadPtr any) error {
	o := NewGoPlaygroundValidator()

	err := o.ValidateStruct(payloadPtr)
	if err != nil {
		msg := o.Translator(err)
		return errors.New(msg)
	}
	return nil
}

// Difference between Go Kit and Go Playground is the way their tags are defined.
// Gookit is easy and more friendly though
//
// Take note conform tag is from the confirm library

// go play ground
// type DogRequest struct {
// 	Name      string `json:"name" validate:"required,min=3,max=12" conform:"trim"`
// 	Age       *int   `json:"age" validate:"required,numeric"`
// 	IsGoodBoy *bool  `json:"isGoodBoy" validate:"required"`
// }

// // go kit
// type ProfileRequest struct {
// 	Name  string `validate:"required|min_len:7" conform:"trim"`
// 	Email string `validate:"email"  message:"email is invalid." label:"User Email"`
// 	Phone string `validate:"required|isE164PhoneNumber"  label:"User Phone"`
// 	Age   int    `validate:"required|int|min:1|max:99"`
// }
