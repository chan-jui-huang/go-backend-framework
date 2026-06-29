package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	server            *http.Server
	config            ServerConfig
	logger            *zap.Logger
	engine            *gin.Engine
	globalMiddlewares *middleware.GlobalMiddlewares
	routers           []route.Router
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

func NewServer(
	config ServerConfig,
	logger *zap.Logger,
	engine *gin.Engine,
	globalMiddlewares *middleware.GlobalMiddlewares,
	routers []route.Router,
) *Server {
	srv := &Server{
		server: &http.Server{
			Addr:              config.Address,
			ReadHeaderTimeout: 30 * time.Minute,
		},
		config:            config,
		logger:            logger,
		engine:            engine,
		globalMiddlewares: globalMiddlewares,
		routers:           routers,
	}

	return srv
}

func (srv *Server) Run() {
	srv.globalMiddlewares.Attach(srv.engine)
	for _, router := range srv.routers {
		router.AttachRoutes()
	}

	srv.server.Handler = srv.engine.Handler()
	if err := srv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		srv.logger.Error(err.Error())
	}
}

func (srv *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, srv.config.GracefulShutdownTtl)
	defer cancel()

	srv.logger.Info("server start to shutdown")
	if err := srv.server.Shutdown(ctx); err != nil {
		return err
	}
	srv.logger.Info("server end to shutdown")

	return nil
}
