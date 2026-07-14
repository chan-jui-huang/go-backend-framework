package registrar

import (
	"context"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHttpServerConfig(loader *config.Loader) *http.ServerConfig {
	config := &http.ServerConfig{}
	loader.Unmarshal("httpServer", config)

	return config
}

type NewHttpServerParams struct {
	fx.In

	Config  *http.ServerConfig
	Logger  *zap.Logger
	Engine  *gin.Engine
	Routers []route.Router `group:"routers"`
}

func NewHttpServer(params NewHttpServerParams) *http.Server {
	return http.NewServer(
		*params.Config,
		params.Logger,
		params.Engine,
		params.Routers,
	)
}

func HttpServerOnStart(_ context.Context, server *http.Server, logger *zap.Logger) error {
	logger.Info("app is starting")
	go server.Run()
	logger.Info("app is started")

	return nil
}

func HttpServerOnStop(ctx context.Context, server *http.Server, logger *zap.Logger) error {
	logger.Info("app is terminating")
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("app is terminated")

	return nil
}
