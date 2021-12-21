package translation

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"

	"github.com/EvisuXiao/andrews-common/utils"
)

var (
	validate     *validator.Validate
	enTranslator ut.Translator
	zhTranslator ut.Translator
)

func Init() {
	enTranslator, _ = ut.New(en.New()).GetTranslator("en")
	zhTranslator, _ = ut.New(zh.New()).GetTranslator("zh")
	validate := binding.Validator.Engine().(*validator.Validate)
	_ = enTranslation.RegisterDefaultTranslations(validate, enTranslator)
	_ = zhTranslation.RegisterDefaultTranslations(validate, zhTranslator)
}

func GetValidator() *validator.Validate {
	return validate
}

func TranslateByEn(err error) error {
	if !utils.HasErr(err) {
		return nil
	}
	validationErr := err.(validator.ValidationErrors)
	for _, vErr := range validationErr {
		return errors.New(vErr.Translate(enTranslator))
	}
	return err
}

func TranslateByZh(err error) error {
	if !utils.HasErr(err) {
		return nil
	}
	validationErr := err.(validator.ValidationErrors)
	for _, vErr := range validationErr {
		return errors.New(vErr.Translate(zhTranslator))
	}
	return err
}
