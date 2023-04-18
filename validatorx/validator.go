package validatorx

import (
	"context"
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"strings"
	"sync"
)

type DiyValidatorInf interface {
	DiyValidator() error
}

var validatorxErr = "validatorxErr"

type validatorx struct {
	validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     map[string]ut.Translator
}

var v *validatorx
var once sync.Once

func NewValidator() *validatorx {
	once.Do(func() {
		v = &validatorx{}
		enx := en.New()
		zh := zh_Hans.New()
		v.Uni = ut.New(zh, enx, zh)
		v.validator = validator.New()
		enTrans, _ := v.Uni.GetTranslator("en")
		zhTrans, _ := v.Uni.GetTranslator("zh")
		v.Trans = make(map[string]ut.Translator)
		v.Trans["en"] = enTrans
		v.Trans["zh"] = zhTrans

		v.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})
		err := en_translations.RegisterDefaultTranslations(v.validator, enTrans)
		if err != nil {
			log.Fatalln(validatorxErr, err)
		}
		err = zh_translations.RegisterDefaultTranslations(v.validator, zhTrans)
		if err != nil {
			log.Fatalln(validatorxErr, err)
		}
	})
	return v
}

func (v *validatorx) validate(ctx context.Context, data interface{}, lang string) error {
	err := v.validator.StructCtx(ctx, data)
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		trans, ok := v.Trans[lang]
		if !ok {
			trans = v.Trans["zh"]
		}
		transData := errs.Translate(trans)
		s := strings.Builder{}
		for _, v := range transData {
			s.WriteString(v)
			s.WriteString("\n")
		}
		return errors.New(s.String())
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return errors.New(invalid.Error())
	}

	return nil
}

var infRtf = reflect.TypeOf((*DiyValidatorInf)(nil)).Elem()

func (v *validatorx) Validate(ctx context.Context, data interface{}, lang string) error {
	//先调用接口验证参数
	rvf := reflect.Indirect(reflect.ValueOf(data))
	rtf := rvf.Type()
	if rtf.Implements(infRtf) {
		//获取反射方法的 err
		method, _ := rtf.MethodByName("DiyValidator")
		res := method.Func.Call([]reflect.Value{rvf})
		if !res[0].IsNil() {
			return res[0].Interface().(error)
		}
	}
	return v.validate(ctx, data, lang)
}
