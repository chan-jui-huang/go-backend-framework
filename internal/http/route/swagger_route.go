package route

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/system"
	"github.com/gin-gonic/gin"
)

type SwaggerRouter struct {
	router         *gin.RouterGroup
	swaggerHandler *system.SwaggerHandler
}

func NewSwaggerRouter(router *gin.Engine, swaggerHandler *system.SwaggerHandler) *SwaggerRouter {
	return &SwaggerRouter{
		router:         router.Group(""),
		swaggerHandler: swaggerHandler,
	}
}

// type [http://localhost:8080/swagger/index.html] in browser to watch the swagger api doc
func (sr *SwaggerRouter) AttachRoutes() {
	if !sr.swaggerHandler.Enabled() {
		return
	}
	sr.router.GET("swagger/*any", sr.swaggerHandler.Handle)
}
