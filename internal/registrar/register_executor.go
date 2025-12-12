package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var RegisterExecutor = registerExecutor{
	registrarCenter: booter.NewRegistrarCenter([]booter.Registrar{
		&ValidatorRegistrar{},
		&LoggerRegistrar{},
		&DatabaseRegistrar{},
		&RedisRegistrar{},
		&AuthenticationRegistrar{},
		&CasbinRegistrar{},
		&MapstructureDecoderRegistrar{},
		&ClickhouseRegistrar{},
	}),
}

type registerExecutor struct {
	registrarCenter *booter.RegistrarCenter
}

func (*registerExecutor) BeforeExecute() {
	config.Registry.RegisterMany(map[string]any{
		"httpServer":           &http.ServerConfig{},
		"middleware.csrf":      &middleware.CsrfConfig{},
		"middleware.rateLimit": &middleware.RateLimitConfig{},
	})
}

func (re *registerExecutor) Execute() {
	re.registrarCenter.Execute()
}

func (*registerExecutor) AfterExecute() {
	service.Registry.SetMany(map[string]any{
		"formDecoder": form.NewDecoder(),
		"modifier":    modifiers.New(),
	})
}
