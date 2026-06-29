package main

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fxApp := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booter.NewDefaultConfig()),
		fx.Provide(
			registrar.NewConfigLoader,
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
		),
	)
	if err := fxApp.Err(); err != nil {
		panic(err)
	}
}
