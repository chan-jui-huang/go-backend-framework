package registrar

import (
	"reflect"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/mitchellh/mapstructure"
)

type MapstructureDecoderRegistrar struct {
	decodeFunc func(input any, output any) error
}

func (mdr *MapstructureDecoderRegistrar) Boot() {
	timeHookFunc := func() mapstructure.DecodeHookFuncType {
		return func(from reflect.Type, to reflect.Type, data any) (any, error) {
			if to != reflect.TypeOf(time.Time{}) {
				return data, nil
			}

			switch from.Kind() {
			case reflect.String:
				return time.Parse(time.RFC3339, data.(string))
			default:
				return data, nil
			}
		}
	}

	mdr.decodeFunc = func(input any, output any) error {
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

func (mdr *MapstructureDecoderRegistrar) Register() {
	current := deps.CurrentService()
	current.MapstructureDecoder = mdr.decodeFunc
	deps.SetService(current)
}
