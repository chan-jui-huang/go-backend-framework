package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var SimpleRegisterExecutor = simpleRegisterExecutor{
	registrarCenter: booter.NewRegistrarCenter([]booter.Registrar{
		&LoggerRegistrar{},
		&AuthenticationRegistrar{},
		&MapstructureDecoderRegistrar{},
	}),
}

type simpleRegisterExecutor struct {
	registrarCenter *booter.RegistrarCenter
}

func (*simpleRegisterExecutor) BeforeExecute() {
	config.Registry.RegisterMany(map[string]any{
		"httpServer":           &http.ServerConfig{},
		"middleware.csrf":      &middleware.CsrfConfig{},
		"middleware.rateLimit": &middleware.RateLimitConfig{},
	})
	currentConfig := deps.CurrentConfig()
	currentConfig.CsrfConfigValue = config.Registry.Get("middleware.csrf").(middleware.CsrfConfig)
	currentConfig.RateLimitConfigValue = config.Registry.Get("middleware.rateLimit").(middleware.RateLimitConfig)
	deps.SetConfig(currentConfig)
}

func (sre *simpleRegisterExecutor) Execute() {
	sre.registrarCenter.Execute()
}

func (*simpleRegisterExecutor) AfterExecute() {
	current := deps.CurrentService()
	current.FormDecoder = form.NewDecoder()
	current.Modifier = modifiers.New()
	deps.SetService(current)
}
