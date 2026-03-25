package route

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/admin"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/user"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	router  *gin.RouterGroup
	routers []Router
}

func NewApiRouter(router *gin.Engine) *ApiRouter {
	routers := []Router{
		user.NewRouter(router),
		admin.NewRouter(router),
	}

	return &ApiRouter{
		router:  router.Group(""),
		routers: routers,
	}
}

// @produce json
// @success 200 {string} string "{"message": "pong"}"
// @router /api/ping [get]
func (ar *ApiRouter) AttachRoutes() {
	ar.router.GET("api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	for _, router := range ar.routers {
		router.AttachRoutes()
	}
}
