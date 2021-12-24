package validator

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"

	"github.com/EvisuXiao/andrews-common/utils"
)

type translator struct {
	localeTranslator    locales.Translator
	universalTranslator ut.Translator
	register            func(*validator.Validate, ut.Translator) error
}

var (
	validate    *validator.Validate
	locale      = "en"
	translators = map[string]*translator{
		"en": {localeTranslator: en.New(), register: enTranslation.RegisterDefaultTranslations},
		"zh": {localeTranslator: zh.New(), register: zhTranslation.RegisterDefaultTranslations},
	}
)

func Init() {
	validate = binding.Validator.Engine().(*validator.Validate)
	for l, t := range translators {
		t.universalTranslator, _ = ut.New(t.localeTranslator).GetTranslator(l)
		_ = t.register(validate, t.universalTranslator)
	}
}

func GetValidator() *validator.Validate {
	return validate
}

func GetTranslator() ut.Translator {
	return translators[locale].universalTranslator
}

func SwitchLocale(l string) error {
	if _, ok := translators[l]; ok {
		locale = l
		return nil
	}
	return errors.New("invalid locale")
}

func Check(v interface{}) error {
	return Translate(GetValidator().Struct(v))
}

func Translate(err error) error {
	if !utils.HasErr(err) {
		return nil
	}
	validationErr := err.(validator.ValidationErrors)
	for _, vErr := range validationErr {
		return errors.New(vErr.Translate(GetTranslator()))
	}
	return err
}
