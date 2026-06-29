package admin

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/http_api"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router            *gin.RouterGroup
	searchHttpApis    *httpapi.SearchHandler
	createPermission  *permission.CreateHandler
	searchPermissions *permission.SearchHandler
	getPermission     *permission.GetHandler
	updatePermission  *permission.UpdateHandler
	deletePermission  *permission.DeleteHandler
	reloadPermission  *permission.ReloadHandler
	createRole        *permission.CreateRoleHandler
	searchRoles       *permission.SearchRolesHandler
	updateRole        *permission.UpdateRoleHandler
	deleteRoles       *permission.DeleteRolesHandler
	updateUserRole    *user.UpdateUserRoleHandler
}

func NewRouter(
	router *gin.Engine,
	authentication *middleware.AuthenticationMiddleware,
	authorization *middleware.AuthorizationMiddleware,
	searchHttpApis *httpapi.SearchHandler,
	createPermission *permission.CreateHandler,
	searchPermissions *permission.SearchHandler,
	getPermission *permission.GetHandler,
	updatePermission *permission.UpdateHandler,
	deletePermission *permission.DeleteHandler,
	reloadPermission *permission.ReloadHandler,
	createRole *permission.CreateRoleHandler,
	searchRoles *permission.SearchRolesHandler,
	updateRole *permission.UpdateRoleHandler,
	deleteRoles *permission.DeleteRolesHandler,
	updateUserRole *user.UpdateUserRoleHandler,
) *Router {
	return &Router{
		router: router.Group(
			"api/admin",
			authentication.Handle(),
			authorization.Handle(),
		),
		searchHttpApis:    searchHttpApis,
		createPermission:  createPermission,
		searchPermissions: searchPermissions,
		getPermission:     getPermission,
		updatePermission:  updatePermission,
		deletePermission:  deletePermission,
		reloadPermission:  reloadPermission,
		createRole:        createRole,
		searchRoles:       searchRoles,
		updateRole:        updateRole,
		deleteRoles:       deleteRoles,
		updateUserRole:    updateUserRole,
	}
}

func (r *Router) AttachRoutes() {
	r.AttachHttpApiRoutes()
	r.AttachPermissionRoutes()
	r.AttachUserRoutes()
}

func (r *Router) AttachHttpApiRoutes() {
	httpApiRouter := r.router.Group("http-api")
	httpApiRouter.GET("", r.searchHttpApis.Handle)
}

func (r *Router) AttachPermissionRoutes() {
	permissionRouter := r.router.Group("permission")
	permissionRouter.POST("", r.createPermission.Handle)
	permissionRouter.GET("", r.searchPermissions.Handle)
	permissionRouter.GET(":id", r.getPermission.Handle)
	permissionRouter.PUT(":id", r.updatePermission.Handle)
	permissionRouter.DELETE("", r.deletePermission.Handle)
	permissionRouter.POST("reload", r.reloadPermission.Handle)

	roleRouter := r.router.Group("role")
	roleRouter.POST("", r.createRole.Handle)
	roleRouter.GET("", r.searchRoles.Handle)
	roleRouter.PUT(":id", r.updateRole.Handle)
	roleRouter.DELETE("", r.deleteRoles.Handle)
}

func (r *Router) AttachUserRoutes() {
	userRoleRouter := r.router.Group("user-role")
	userRoleRouter.PUT("", r.updateUserRole.Handle)
}
