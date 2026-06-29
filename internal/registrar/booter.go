package registrar

import (
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewConfigLoader(booterConfig *booter.Config) *config.Loader {
	return booter.BootConfigLoader(booterConfig)
}
