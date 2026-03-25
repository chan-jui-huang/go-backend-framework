package registrar

import (
	"context"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidatorOnStart(context.Context) error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		for _, tagName := range []string{"json", "form"} {
			tagValue := field.Tag.Get(tagName)
			if tagValue == "" || tagValue == "-" {
				continue
			}

			name := strings.TrimSpace(strings.Split(tagValue, ",")[0])
			if name != "" {
				return name
			}
		}

		return field.Name
	})

	return nil
}
