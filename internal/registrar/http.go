package registrar

import (
	"context"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewHttpServerConfig(loader *booterconfig.Loader) *http.ServerConfig {
	config := &http.ServerConfig{}
	loader.Unmarshal("httpServer", config)

	return config
}

func NewHttpServer(config *http.ServerConfig) *http.Server {
	return http.NewServer(*config)
}

func HttpServerOnStart(_ context.Context, server *http.Server) error {
	logger := deps.Logger()
	logger.Info("app is starting")
	go server.Run()
	logger.Info("app is started")

	return nil
}

func HttpServerOnStop(ctx context.Context, server *http.Server) error {
	logger := deps.Logger()
	logger.Info("app is terminating")
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("app is terminated")

	return nil
}
