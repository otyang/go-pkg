package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/leebenson/conform"
)

var _ IValidators = (*GoPlayground)(nil)

type GoPlayground struct {
	vald   *validator.Validate
	uni    *ut.UniversalTranslator
	Trans  ut.Translator
	locale locales.Translator
}

func NewGoPlaygroundValidator() *GoPlayground {
	o := &GoPlayground{}

	o.locale = en.New()
	o.uni = ut.New(o.locale, o.locale)

	t, found := o.uni.GetTranslator("en")
	if !found {
		fmt.Println("validator goplayground translator not found")
	}

	o.Trans = t
	en_translations.RegisterDefaultTranslations(o.vald, o.Trans)

	o.vald.RegisterTagNameFunc(o.goPlayExtractNamesOfFormFromTags) // enables us use the form tags as field names e.g form:”name"

	_ = o.vald.RegisterValidation("alpha_space", o.isAlphaSpace) // lets register custom validation

	return o
}

func (o *GoPlayground) ValidateStruct(vPtr any) error {
	conform.Strings(vPtr) // confirm library
	return o.vald.Struct(vPtr)
}

func (o *GoPlayground) Translator(err error) string {
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			var translatedErrs []string
			for _, e := range errs {
				translatedErrs = append(translatedErrs, e.Translate(o.Trans))
			}
			return strings.Join(translatedErrs, "\n")
		}
		return err.Error()
	}
	return ""
}

// enables us use the form tags as field names e.g form:”name"
func (o *GoPlayground) goPlayExtractNamesOfFormFromTags(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func (o *GoPlayground) isAlphaSpace(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^[a-zA-Z ]+$").MatchString(fl.Field().String())
}
