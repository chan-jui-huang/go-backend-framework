package admin

import (
	httpapi "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/http_api"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.RouterGroup
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{
		router: router.Group(
			"api/admin",
			middleware.Authenticate(),
			middleware.Authorize(),
		),
	}
}

func (r *Router) AttachRoutes() {
	r.AttachHttpApiRoutes()
	r.AttachPermissionRoutes()
	r.AttachUserRoutes()
}

func (r *Router) AttachHttpApiRoutes() {
	httpApiRouter := r.router.Group("http-api")
	httpApiRouter.GET("", httpapi.Search)
}

func (r *Router) AttachPermissionRoutes() {
	permissionRouter := r.router.Group("permission")
	permissionRouter.POST("", permission.Create)
	permissionRouter.GET("", permission.Search)
	permissionRouter.GET(":id", permission.Get)
	permissionRouter.PUT(":id", permission.Update)
	permissionRouter.DELETE("", permission.Delete)
	permissionRouter.POST("reload", permission.Reload)

	roleRouter := r.router.Group("role")
	roleRouter.POST("", permission.CreateRole)
	roleRouter.GET("", permission.SearchRoles)
	roleRouter.PUT(":id", permission.UpdateRole)
	roleRouter.DELETE("", permission.DeleteRoles)
}

func (r *Router) AttachUserRoutes() {
	userRoleRouter := r.router.Group("user-role")
	userRoleRouter.PUT("", user.UpdateUserRole)
}
