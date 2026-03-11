package registrar

import (
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewConfigLoader(booterConfig *booter.Config) *booterconfig.Loader {
	return booter.BootConfigLoader(booterConfig)
}
