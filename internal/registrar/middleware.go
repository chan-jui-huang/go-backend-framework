package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewCsrfConfig(loader *booterconfig.Loader) *middleware.CsrfConfig {
	config := &middleware.CsrfConfig{}
	loader.Unmarshal("middleware.csrf", config)

	return config
}

func NewRateLimitConfig(loader *booterconfig.Loader) *middleware.RateLimitConfig {
	config := &middleware.RateLimitConfig{}
	loader.Unmarshal("middleware.rateLimit", config)

	return config
}
