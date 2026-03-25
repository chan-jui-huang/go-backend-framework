package user

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.RouterGroup
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{
		router: router.Group("api/user"),
	}
}

func (r *Router) AttachRoutes() {
	r.router.POST("register", user.Register)
	r.router.POST("login", user.Login)
	r.router.GET("me", middleware.Authenticate(), user.Me)
	r.router.PUT("", middleware.Authenticate(), user.Update)
	r.router.PUT("password", middleware.Authenticate(), user.UpdatePassword)
}
