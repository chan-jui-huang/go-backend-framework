package main

import (
	"fmt"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/route"
	appregistrar "github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	routers := []route.Router{
		route.NewApiRouter(engine),
		route.NewSwaggerRouter(engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	for _, routeInfo := range engine.Routes() {
		fmt.Printf("method: [%s], path: [%s], handler: [%s]\n", routeInfo.Method, routeInfo.Path, routeInfo.Handler)
	}
}
