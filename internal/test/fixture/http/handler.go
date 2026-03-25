package http

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	pkgHttp "github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/route"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
}

func New(handlerFuncs ...gin.HandlerFunc) *Handler {
	engine, err := pkgHttp.NewEngine()
	if err != nil {
		panic(err)
	}

	handler := &Handler{
		engine: engine,
	}
	handler.attachGlobalMiddleware()
	if len(handlerFuncs) > 0 {
		handler.engine.Use(handlerFuncs...)
	}
	routers := []route.Router{
		route.NewApiRouter(handler.engine),
		route.NewSwaggerRouter(handler.engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	return handler
}

func (handler *Handler) attachGlobalMiddleware() {
	csrfConfig := deps.CsrfConfig()

	handler.engine.Use(
		middleware.VerifyCsrfToken(csrfConfig),
	)
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.engine.ServeHTTP(w, req)
}
