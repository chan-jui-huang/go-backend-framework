package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type HttpRouteRunnerParams struct {
	fx.In

	Engine  *gin.Engine
	Routers []route.Router `group:"routers"`
}

type HttpRouteRunner struct {
	engine  *gin.Engine
	routers []route.Router
}

func NewHttpRouteRunner(params HttpRouteRunnerParams) *HttpRouteRunner {
	return &HttpRouteRunner{
		engine:  params.Engine,
		routers: params.Routers,
	}
}

func (r *HttpRouteRunner) Run() error {
	for _, router := range r.routers {
		router.AttachRoutes()
	}

	for _, routeInfo := range r.engine.Routes() {
		fmt.Printf("method: [%s], path: [%s], handler: [%s]\n", routeInfo.Method, routeInfo.Path, routeInfo.Handler)
	}

	return nil
}

func main() {
	var runner *HttpRouteRunner

	gin.SetMode(gin.ReleaseMode)
	fxApp := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booter.NewDefaultConfig()),
		route.NewModule(),
		fx.Provide(
			http.NewEngine,
			registrar.NewConfigLoader,
			registrar.NewAuthenticationConfig,
			registrar.NewAuthenticator,
			registrar.NewDatabaseConfig,
			registrar.NewDatabase,
			registrar.NewLoggerConfigs,
			registrar.NewLoggers,
			registrar.NewCasbinEnforcer,
			NewHttpRouteRunner,
		),
		fx.Populate(&runner),
	)
	if err := fxApp.Err(); err != nil {
		log.Fatal(err)
	}

	startCtx, cancelStart := context.WithTimeout(context.Background(), 15*time.Second)
	startErr := fxApp.Start(startCtx)
	cancelStart()
	if startErr != nil {
		log.Fatal(startErr)
	}

	runErr := runner.Run()

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 15*time.Second)
	stopErr := fxApp.Stop(stopCtx)
	cancelStop()

	if runErr != nil {
		log.Fatal(runErr)
	}
	if stopErr != nil {
		log.Fatal(stopErr)
	}
}
