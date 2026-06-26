package registrar

import (
	"context"

	internalhttp "github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHttpServerConfig(loader *booterconfig.Loader) *internalhttp.ServerConfig {
	config := &internalhttp.ServerConfig{}
	loader.Unmarshal("httpServer", config)

	return config
}

type NewHttpServerParams struct {
	fx.In

	Config            *internalhttp.ServerConfig
	Logger            *zap.Logger
	Engine            *gin.Engine
	GlobalMiddlewares *middleware.GlobalMiddlewares
	Routers           []route.Router `group:"routers"`
}

func NewHttpServer(params NewHttpServerParams) *internalhttp.Server {
	return internalhttp.NewServer(
		*params.Config,
		params.Logger,
		params.Engine,
		params.GlobalMiddlewares,
		params.Routers,
	)
}

func HttpServerOnStart(_ context.Context, server *internalhttp.Server, logger *zap.Logger) error {
	logger.Info("app is starting")
	go server.Run()
	logger.Info("app is started")

	return nil
}

func HttpServerOnStop(ctx context.Context, server *internalhttp.Server, logger *zap.Logger) error {
	logger.Info("app is terminating")
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("app is terminated")

	return nil
}
