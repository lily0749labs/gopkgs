// Package ginUtil 提供 Gin 框架相关的通用工具。
package ginUtil

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslation "github.com/go-playground/validator/v10/translations/zh"
)

var (
	ErrNilValidator       = errors.New("validator must not be nil")
	ErrTranslatorNotFound = errors.New("Chinese translator not found")
	translator            ut.Translator
)

// RegisterChinese 为指定 validator 注册中文默认翻译。
func RegisterChinese(validate *validator.Validate) (ut.Translator, error) {
	if validate == nil {
		return nil, ErrNilValidator
	}

	locale := zh.New()
	trans, found := ut.New(locale, locale).GetTranslator("zh")
	if !found {
		return nil, ErrTranslatorNotFound
	}
	if err := zhtranslation.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}
	return trans, nil
}

// InitValidatorTranslation 为 Gin 的全局 validator 注册中文翻译。
// 为兼容原有服务，注册错误会被忽略。
func InitValidatorTranslation() {
	translator, _ = RegisterChinese(binding.Validator.Engine().(*validator.Validate))
}

// Translate 使用指定翻译器转换 validator.ValidationErrors。
// 非验证错误保持原错误文本；多个字段继续使用旧版 </br> 分隔格式。
func Translate(err error, trans ut.Translator) string {
	if err == nil {
		return ""
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var result strings.Builder
	for _, fieldError := range validationErrors {
		result.WriteString(fieldError.Translate(trans))
		result.WriteString("</br>")
	}
	return result.String()
}

// VaTrans 使用 Gin 全局 validator 对验证错误进行中文翻译。
func VaTrans(err error) string {
	return Translate(err, translator)
}
