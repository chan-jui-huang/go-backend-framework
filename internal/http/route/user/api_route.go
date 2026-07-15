package user

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/requestlog"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router          *gin.RouterGroup
	registerHandler *user.RegisterHandler
	loginHandler    *user.LoginHandler
	getMeHandler    *user.GetMeHandler
	updateHandler   *user.UpdateHandler
	updatePassword  *user.UpdatePasswordHandler
	authentication  *middleware.AuthenticationMiddleware
}

func NewRouter(
	router *gin.Engine,
	registerHandler *user.RegisterHandler,
	loginHandler *user.LoginHandler,
	getMeHandler *user.GetMeHandler,
	updateHandler *user.UpdateHandler,
	updatePassword *user.UpdatePasswordHandler,
	authentication *middleware.AuthenticationMiddleware,
) *Router {
	return &Router{
		router:          router.Group("api/user"),
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
		getMeHandler:    getMeHandler,
		updateHandler:   updateHandler,
		updatePassword:  updatePassword,
		authentication:  authentication,
	}
}

func (r *Router) AttachRoutes() {
	r.router.POST("register", requestlog.WithPolicy(requestlog.Policy{
		ErrorLog: []string{"name", "email"},
	}), r.registerHandler.Handle)
	r.router.POST("login", requestlog.WithPolicy(requestlog.Policy{
		ErrorLog: []string{"email"},
	}), r.loginHandler.Handle)
	r.router.GET("me", r.authentication.Handle(), r.getMeHandler.Handle)
	r.router.PUT("", r.authentication.Handle(), r.updateHandler.Handle)
	r.router.PUT("password", requestlog.WithPolicy(requestlog.Policy{}), r.authentication.Handle(), r.updatePassword.Handle)
}
