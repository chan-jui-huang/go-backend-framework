package registrar

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func NewMapstructureDecoder() func(any, any) error {
	timeHookFunc := func() mapstructure.DecodeHookFuncType {
		return func(from reflect.Type, to reflect.Type, data any) (any, error) {
			if to != reflect.TypeOf(time.Time{}) {
				return data, nil
			}
			if from.Kind() != reflect.String {
				return data, nil
			}

			return time.Parse(time.RFC3339, data.(string))
		}
	}

	return func(input any, output any) error {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata:   nil,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(timeHookFunc()),
			Result:     output,
		})
		if err != nil {
			return err
		}

		return decoder.Decode(input)
	}
}
