package route

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/system"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/admin"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/user"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	router      *gin.RouterGroup
	pingHandler *system.PingHandler
	routers     []Router
}

func NewApiRouter(
	router *gin.Engine,
	pingHandler *system.PingHandler,
	adminRouter *admin.Router,
	userRouter *user.Router,
) *ApiRouter {
	return &ApiRouter{
		router:      router.Group(""),
		pingHandler: pingHandler,
		routers:     []Router{userRouter, adminRouter},
	}
}

// @produce json
// @success 200 {string} string "{"message": "pong"}"
// @router /api/ping [get]
func (ar *ApiRouter) AttachRoutes() {
	ar.router.GET("api/ping", ar.pingHandler.Handle)

	for _, router := range ar.routers {
		router.AttachRoutes()
	}
}
