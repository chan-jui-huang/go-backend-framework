package requestlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

const maxBodySize = 1 << 20

var (
	errBodyTooLarge          = errors.New("request body exceeds logging size limit")
	errInvalidJSON           = errors.New("request body is not a valid JSON object")
	errInvalidAllowlistPath  = errors.New("request body allowlist contains an invalid path")
	errForbiddenAllowlistKey = errors.New("request body allowlist contains a forbidden field")
	errSensitiveBody         = errors.New("request body contains a forbidden field")
)

var forbiddenFields = map[string]struct{}{
	"password":        {},
	"currentpassword": {},
	"confirmpassword": {},
	"token":           {},
	"accesstoken":     {},
	"refreshtoken":    {},
	"authorization":   {},
	"secret":          {},
	"clientsecret":    {},
	"clienttoken":     {},
	"apikey":          {},
	"privatekey":      {},
	"credential":      {},
	"credentials":     {},
}

var forbiddenFieldSuffixes = []string{
	"password",
	"secret",
	"token",
	"authorization",
	"apikey",
	"privatekey",
	"passphrase",
	"credential",
	"credentials",
}

type fieldTree map[string]fieldTree

// Filter returns a compact JSON body containing only fields allowed for the destination.
func Filter(c *gin.Context, destination Destination) ([]byte, error) {
	fields, hasPolicy := fieldsFor(c, destination)
	if hasPolicy && len(fields) == 0 {
		return nil, nil
	}

	source, err := decodeBody(c)
	if err != nil || source == nil {
		return nil, err
	}

	if !hasPolicy {
		if containsForbiddenField(source) {
			return nil, errSensitiveBody
		}

		return marshalBody(source)
	}

	tree, err := compileFieldTree(fields)
	if err != nil {
		return nil, err
	}

	filtered := selectObject(source, tree)
	if len(filtered) == 0 {
		return nil, nil
	}
	if containsForbiddenField(filtered) {
		return nil, errForbiddenAllowlistKey
	}

	return marshalBody(filtered)
}

func decodeBody(c *gin.Context) (map[string]any, error) {
	bodyValue, ok := c.Get(gin.BodyBytesKey)
	if !ok {
		return nil, nil
	}

	body, ok := bodyValue.([]byte)
	if !ok || len(body) == 0 {
		return nil, nil
	}
	if len(body) > maxBodySize {
		return nil, errBodyTooLarge
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()

	var source map[string]any
	if err := decoder.Decode(&source); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidJSON, err)
	}
	if source == nil {
		return nil, errInvalidJSON
	}

	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, errInvalidJSON
	}

	return source, nil
}

func marshalBody(body map[string]any) ([]byte, error) {
	result, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body for logging: %w", err)
	}

	return result, nil
}

func compileFieldTree(fields []string) (fieldTree, error) {
	tree := fieldTree{}
	for _, field := range fields {
		parts := strings.Split(field, ".")
		current := tree
		for _, part := range parts {
			if part == "" {
				return nil, errInvalidAllowlistPath
			}
			if isForbiddenField(part) {
				return nil, fmt.Errorf("%w: %s", errForbiddenAllowlistKey, field)
			}

			next, ok := current[part]
			if !ok {
				next = fieldTree{}
				current[part] = next
			}
			current = next
		}
	}

	return tree, nil
}

func selectObject(source map[string]any, tree fieldTree) map[string]any {
	result := make(map[string]any, len(tree))
	for field, children := range tree {
		value, ok := source[field]
		if !ok {
			continue
		}

		if len(children) == 0 {
			result[field] = value
			continue
		}

		filtered, ok := selectValue(value, children)
		if ok {
			result[field] = filtered
		}
	}

	return result
}

func selectValue(value any, tree fieldTree) (any, bool) {
	switch typedValue := value.(type) {
	case map[string]any:
		filtered := selectObject(typedValue, tree)
		return filtered, len(filtered) > 0
	case []any:
		filtered := make([]any, 0, len(typedValue))
		for _, item := range typedValue {
			selected, ok := selectValue(item, tree)
			if ok {
				filtered = append(filtered, selected)
			}
		}
		return filtered, len(filtered) > 0
	default:
		return nil, false
	}
}

func containsForbiddenField(value any) bool {
	switch typedValue := value.(type) {
	case map[string]any:
		for field, nestedValue := range typedValue {
			if isForbiddenField(field) || containsForbiddenField(nestedValue) {
				return true
			}
		}
	case []any:
		for _, nestedValue := range typedValue {
			if containsForbiddenField(nestedValue) {
				return true
			}
		}
	}

	return false
}

func isForbiddenField(field string) bool {
	normalized := strings.Map(func(r rune) rune {
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, field)
	if _, ok := forbiddenFields[normalized]; ok {
		return true
	}
	for _, suffix := range forbiddenFieldSuffixes {
		if strings.HasSuffix(normalized, suffix) {
			return true
		}
	}

	return false
}
