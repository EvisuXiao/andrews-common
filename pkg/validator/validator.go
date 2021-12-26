package validator

import (
	"errors"
	"reflect"
	"time"

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
	var err error
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		return Translate(validate.Struct(v))
	case reflect.Array, reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if rv.Index(i).CanInterface() {
				if err = Check(rv.Index(i).Interface()); utils.HasErr(err) {
					return err
				}
			}
		}
		return nil
	case reflect.Map:
		for _, mrv := range rv.MapKeys() {
			if rv.MapIndex(mrv).CanInterface() {
				if err = Check(rv.MapIndex(mrv).Interface()); utils.HasErr(err) {
					return err
				}
			}
		}
		return nil
	default:
		if rv.Type() == reflect.TypeOf(time.Time{}) {
			return Translate(validate.Struct(v))
		}
		return errors.New("not struct type")
	}
}

func Translate(err error) error {
	if !utils.HasErr(err) {
		return nil
	}
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, vErr := range validationErr {
			return errors.New(vErr.Translate(GetTranslator()))
		}
	}
	return err
}
