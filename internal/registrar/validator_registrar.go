package registrar

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidatorRegistrar configures gin's validator behavior.
// It maps ValidationErrors field names to json/form tags.
type ValidatorRegistrar struct{}

func (*ValidatorRegistrar) Boot() {}

func (vr *ValidatorRegistrar) Register() {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		if name := vr.pickTagName(field.Tag.Get("json")); name != "" {
			return name
		}
		if name := vr.pickTagName(field.Tag.Get("form")); name != "" {
			return name
		}
		return field.Name
	})
}

func (*ValidatorRegistrar) pickTagName(tagValue string) string {
	if tagValue == "" || tagValue == "-" {
		return ""
	}
	name := strings.Split(tagValue, ",")[0]
	return strings.TrimSpace(name)
}
