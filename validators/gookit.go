package validators

import (
	"regexp"

	"github.com/gookit/validate"
	"github.com/leebenson/conform"
)

var _ IValidators = (*GooKit)(nil)

type GooKit struct{}

func NewGooKitValidator() *GooKit {
	return &GooKit{}
}

func (o *GooKit) ValidateStruct(vPtr any) error {
	conform.Strings(vPtr) // confirm

	v := validate.New(vPtr)
	v.SkipOnEmpty = false // should be false
	v.StopOnError = true

	o.initGlobalValidator()
	o.initGlobalMsg(v)

	v.Validate()

	if v.Errors.Empty() {
		return nil
	}
	return v.Errors
}

func (o *GooKit) Translator(err error) string {
	if err != nil {
		val, ok := err.(validate.Errors)
		if ok {
			return string(val.String())
		}
		return err.Error()
	}
	return ""
}

func (o *GooKit) initGlobalMsg(v *validate.Validation) {
	v.AddMessages(map[string]string{
		"minLength":         "{field} min length is %d",
		"isE164PhoneNumber": "Invalid Phone Number Format.",
	})
}

func (o *GooKit) initGlobalValidator() {
	validate.AddValidator(
		"isE164PhoneNumber",
		func(val string) bool {
			return regexp.MustCompile(
				`^\+?[1-9]\d{1,14}$`).Match([]byte(val))
		},
	)
}
