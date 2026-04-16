package main

import (
	appregistrar "github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
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
			appregistrar.NewConfigLoader,
			appregistrar.NewAuthenticationConfig,
			appregistrar.NewAuthenticator,
			appregistrar.NewDatabaseConfig,
			appregistrar.NewDatabase,
			appregistrar.NewRedisConfig,
			appregistrar.NewRedis,
			appregistrar.NewClickhouseConfig,
			appregistrar.NewClickhouse,
			appregistrar.NewLoggerConfigs,
			appregistrar.NewLoggers,
			appregistrar.NewCasbinEnforcer,
		),
		fx.Invoke(
			appregistrar.RegisterConfigDependencies,
			appregistrar.RegisterServiceDependencies,
		),
	)
	if err := fxApp.Err(); err != nil {
		panic(err)
	}
}
