package main

import (
	"time"

	internalhttp "github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	booterConfig := booter.NewConfigWithCommand()

	fx.New(
		fx.StopTimeout(60*time.Second),
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booterConfig),
		registrar.NewModule(),
		route.NewModule(),
		fx.Invoke(
			func(*internalhttp.Server) {},
		),
	).Run()
}
