package main

import (
	appregistrar "github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
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
			appregistrar.NewMapstructureDecoder,
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
