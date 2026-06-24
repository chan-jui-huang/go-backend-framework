package route

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/user"
	"go.uber.org/fx"
)

// NewModule provides HTTP routes and their controller domains.
func NewModule() fx.Option {
	return fx.Module(
		"http.route",
		admin.NewModule(),
		user.NewModule(),
		fx.Provide(
			NewApiRouter,
			NewSwaggerRouter,
		),
	)
}
