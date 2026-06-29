package route

import (
	admincontroller "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/system"
	usercontroller "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/admin"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/user"
	"go.uber.org/fx"
)

// NewModule provides HTTP routes and their controller domains.
func NewModule() fx.Option {
	return fx.Module(
		"http.route",
		middleware.NewModule(),
		admincontroller.NewModule(),
		system.NewModule(),
		usercontroller.NewModule(),
		fx.Provide(
			admin.NewRouter,
			user.NewRouter,
			fx.Annotate(
				NewApiRouter,
				fx.As(new(Router)),
				fx.ResultTags(`group:"routers"`),
			),
			fx.Annotate(
				NewSwaggerRouter,
				fx.As(new(Router)),
				fx.ResultTags(`group:"routers"`),
			),
		),
	)
}
