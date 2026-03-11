package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/config"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewCsrfConfig(loader *booterconfig.Loader) *config.CsrfConfig {
	cfg := &config.CsrfConfig{}
	loader.Unmarshal("middleware.csrf", cfg)

	return cfg
}

func NewRateLimitConfig(loader *booterconfig.Loader) *config.RateLimitConfig {
	cfg := &config.RateLimitConfig{}
	loader.Unmarshal("middleware.rateLimit", cfg)

	return cfg
}
