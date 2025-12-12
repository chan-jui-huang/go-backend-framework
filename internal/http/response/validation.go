package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// MakeValidationErrorContext extracts validation errors into a map for response context
func MakeValidationErrorContext(err error) map[string]any {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]any, len(ve))
		for _, fe := range ve {
			out[fe.Field()] = fe.Tag()
		}
		return out
	}
	return nil
}
