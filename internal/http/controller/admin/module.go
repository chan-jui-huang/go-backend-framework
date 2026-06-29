package admin

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/http_api"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/user"
	"go.uber.org/fx"
)

// NewModule combines the admin controller domain modules.
func NewModule() fx.Option {
	return fx.Module(
		"http.controller.admin",
		httpapi.NewModule(),
		permission.NewModule(),
		user.NewModule(),
	)
}
