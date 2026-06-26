package http

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	pkgHttp "github.com/chan-jui-huang/go-backend-framework/v3/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Dependencies struct {
	fx.In

	Engine     *gin.Engine
	CsrfConfig *config.CsrfConfig
}

type RouteParams struct {
	fx.In

	Routers []route.Router `group:"routers"`
}

type Handler struct {
	engine     *gin.Engine
	csrfConfig *config.CsrfConfig
}

func NewEngine(globalMiddlewares *middleware.GlobalMiddlewares) *gin.Engine {
	engine, err := pkgHttp.NewEngine()
	if err != nil {
		panic(err)
	}
	globalMiddlewares.Attach(engine)

	return engine
}

func New(dependencies Dependencies, routeParams RouteParams) *Handler {
	handler := &Handler{
		engine:     dependencies.Engine,
		csrfConfig: dependencies.CsrfConfig,
	}

	for _, router := range routeParams.Routers {
		router.AttachRoutes()
	}

	return handler
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.engine.ServeHTTP(w, req)
}
