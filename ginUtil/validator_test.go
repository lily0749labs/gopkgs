package ginUtil

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestRegisterChineseAndTranslate(t *testing.T) {
	validate := validator.New()
	trans, err := RegisterChinese(validate)
	if err != nil {
		t.Fatalf("RegisterChinese() error = %v", err)
	}

	type request struct {
		Name string `validate:"required"`
	}
	err = validate.Struct(request{})
	translated := Translate(err, trans)
	if !strings.Contains(translated, "Name") || !strings.HasSuffix(translated, "</br>") {
		t.Fatalf("Translate() = %q", translated)
	}
}

func TestTranslateNonValidationError(t *testing.T) {
	err := errors.New("plain error")
	if got := Translate(err, nil); got != err.Error() {
		t.Fatalf("Translate() = %q, want %q", got, err.Error())
	}
}

func TestRegisterChineseNilValidator(t *testing.T) {
	_, err := RegisterChinese(nil)
	if !errors.Is(err, ErrNilValidator) {
		t.Fatalf("RegisterChinese() error = %v, want ErrNilValidator", err)
	}
}
