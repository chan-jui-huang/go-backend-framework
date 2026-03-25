package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/route"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	config ServerConfig
}

type ServerConfig struct {
	Address             string
	GracefulShutdownTtl time.Duration
}

func NewEngine() (*gin.Engine, error) {
	engine := gin.New()
	engine.RemoteIPHeaders = []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}
	err := engine.SetTrustedProxies([]string{
		"0.0.0.0/0",
		"::/0",
	})

	return engine, err
}

func NewServer(config ServerConfig) *Server {
	srv := &Server{
		server: &http.Server{
			Addr:              config.Address,
			ReadHeaderTimeout: 30 * time.Minute,
		},
		config: config,
	}

	return srv
}

func (srv *Server) Run() {
	engine, err := NewEngine()
	if err != nil {
		panic(err)
	}
	middleware.AttachGlobalMiddleware(engine)

	routers := []route.Router{
		route.NewApiRouter(engine),
		route.NewSwaggerRouter(engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	srv.server.Handler = engine.Handler()
	logger := deps.Logger()

	if err := srv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}

func (srv *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, srv.config.GracefulShutdownTtl)
	defer cancel()

	logger := deps.Logger()
	logger.Info("server start to shutdown")
	if err := srv.server.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("server end to shutdown")

	return nil
}
