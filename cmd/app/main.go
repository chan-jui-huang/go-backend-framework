package main

import (
	"time"

	internalhttp "github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	booterConfig := booter.NewConfigWithCommand()

	fx.New(
		fx.StopTimeout(60*time.Second),
		fx.Supply(booterConfig),
		fx.Provide(
			registrar.NewConfigLoader,
			registrar.NewHttpServerConfig,
			fx.Annotate(
				registrar.NewHttpServer,
				fx.OnStart(registrar.HttpServerOnStart),
				fx.OnStop(registrar.HttpServerOnStop),
			),
			registrar.NewCsrfConfig,
			registrar.NewRateLimitConfig,
			registrar.NewAuthenticationConfig,
			registrar.NewAuthenticator,
			registrar.NewDatabaseConfig,
			registrar.NewDatabase,
			registrar.NewRedisConfig,
			registrar.NewRedis,
			registrar.NewClickhouseConfig,
			registrar.NewClickhouse,
			registrar.NewLoggerConfigs,
			registrar.NewLoggers,
			registrar.NewCasbinEnforcer,
			registrar.NewMapstructureDecoder,
			form.NewDecoder,
			modifiers.New,
		),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(registrar.ValidatorOnStart),
			),
			fx.Annotate(
				func() {},
				fx.OnStart(registrar.SchedulerOnStart),
				fx.OnStop(registrar.SchedulerOnStop),
			),
			registrar.RegisterConfigDependencies,
			registrar.RegisterServiceDependencies,
			func(*internalhttp.Server) {},
		),
	).Run()
}
