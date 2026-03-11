package registrar

import (
	"context"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/scheduler"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHttpServerConfig(loader *booterconfig.Loader) *http.ServerConfig {
	config := &http.ServerConfig{}
	loader.Unmarshal("httpServer", config)

	return config
}

func NewHttpServer(config *http.ServerConfig, logger *zap.Logger, lc fx.Lifecycle) *http.Server {
	server := http.NewServer(*config)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("app is starting")
			go server.Run()
			logger.Info("app is started")
			scheduler.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("app is terminating")
			scheduler.Stop()
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			logger.Info("app is terminated")

			return nil
		},
	})

	return server
}
