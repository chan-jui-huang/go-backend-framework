package registrar

import (
	internalhttp "github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"go.uber.org/fx"
)

// NewModule provides the application's infrastructure dependencies.
func NewModule() fx.Option {
	return fx.Module(
		"registrar",
		fx.Provide(
			NewConfigLoader,
			internalhttp.NewEngine,
			NewHttpServerConfig,
			fx.Annotate(
				NewHttpServer,
				fx.OnStart(HttpServerOnStart),
				fx.OnStop(HttpServerOnStop),
			),
			NewCsrfConfig,
			NewRateLimitConfig,
			NewAuthenticationConfig,
			NewAuthenticator,
			NewDatabaseConfig,
			NewDatabase,
			NewRedisConfig,
			NewRedis,
			NewClickhouseConfig,
			NewClickhouse,
			NewLoggerConfigs,
			NewLoggers,
			NewCasbinEnforcer,
			form.NewDecoder,
			modifiers.New,
		),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(ValidatorOnStart),
			),
		),
	)
}
