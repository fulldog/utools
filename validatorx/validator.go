package validatorx

import (
	"context"
	"errors"
	"github.com/fulldog/utools"
	"github.com/fulldog/utools/timex"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/google/uuid"
	"log"
	"reflect"
	"strings"
)

type DiyValidatorInf interface {
	DiyValidator() error
}
type Validatorx interface {
	Validate(ctx context.Context, data interface{}, lang string) error
}

var validatorxErr = "validatorxErr"

type validatorx struct {
	validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     map[string]ut.Translator
}

var v *validatorx

func GetValidatorx() *validatorx {
	return v
}

func init() {
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

	_ = v.validator.RegisterValidation("phone", validatePhone)
	_ = v.validator.RegisterValidation("uuid", validateUuidx)
	_ = v.validator.RegisterValidation("timex", validateTimex)

	_ = v.validator.RegisterTranslation("timex", zhTrans, func(ut ut.Translator) error {
		return ut.Add("timex", "{0}必须是一个有效的时间格式", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("timex", fe.Field())
		return t
	})

	_ = v.validator.RegisterTranslation("phone", zhTrans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0}必须是一个有效的手机号码", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	err := en_translations.RegisterDefaultTranslations(v.validator, enTrans)
	if err != nil {
		log.Fatalln(err)
	}
	err = zh_translations.RegisterDefaultTranslations(v.validator, zhTrans)
	if err != nil {
		log.Fatalln(err)
	}
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

func (v validatorx) Validate(ctx context.Context, data interface{}, lang string) error {
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

// 手机号码
func validatePhone(fl validator.FieldLevel) bool {
	return utools.VerifyFormat(fl.Field().Interface().(string), "^(13|14|15|16|17|18|19)[0-9]{9}$")
}

// uuid
func validateUuidx(fl validator.FieldLevel) bool {
	if _, err := uuid.Parse(fl.Field().String()); err != nil {
		return false
	}
	return true
}

// time
func validateTimex(fl validator.FieldLevel) bool {
	return timex.Parse(fl.Field().String()).IsZero()
}
