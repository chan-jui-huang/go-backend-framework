package seeder

import (
	"github.com/casbin/casbin/v3"
	httpapi "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/http_api"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/permission"
	adminuser "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/admin/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/system"
	usercontroller "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	adminroute "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/admin"
	userroute "github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route/user"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	pkgpermission "github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HttpApiSeeder struct {
	database       *gorm.DB
	logger         *zap.Logger
	authenticator  *authentication.Authenticator
	casbinEnforcer *casbin.SyncedCachedEnforcer
	booterConfig   *booter.Config
}

func NewHttpApiSeeder(
	database *gorm.DB,
	logger *zap.Logger,
	authenticator *authentication.Authenticator,
	casbinEnforcer *casbin.SyncedCachedEnforcer,
	booterConfig *booter.Config,
) *HttpApiSeeder {
	return &HttpApiSeeder{
		database:       database,
		logger:         logger,
		authenticator:  authenticator,
		casbinEnforcer: casbinEnforcer,
		booterConfig:   booterConfig,
	}
}

func (s *HttpApiSeeder) Run(tx *gorm.DB) error {
	gin.SetMode(gin.ReleaseMode)
	httpApis, err := pkgpermission.GetHttpApis(tx, "")
	if err != nil {
		return err
	}

	engine := gin.New()
	authentication := middleware.NewAuthenticationMiddleware(s.logger, s.authenticator)
	authorization := middleware.NewAuthorizationMiddleware(s.logger, s.casbinEnforcer)
	adminRouter := adminroute.NewRouter(
		engine, authentication, authorization,
		httpapi.NewSearchHandler(s.database, s.logger),
		permission.NewCreateHandler(s.database, s.casbinEnforcer, s.logger),
		permission.NewSearchHandler(s.database, s.logger),
		permission.NewGetHandler(s.database, s.logger),
		permission.NewUpdateHandler(s.database, s.casbinEnforcer, s.logger),
		permission.NewDeleteHandler(s.database, s.casbinEnforcer, s.logger),
		permission.NewReloadHandler(s.casbinEnforcer, s.logger),
		permission.NewCreateRoleHandler(s.database, s.logger),
		permission.NewSearchRolesHandler(s.database, s.logger),
		permission.NewUpdateRoleHandler(s.database, s.casbinEnforcer, s.logger),
		permission.NewDeleteRolesHandler(s.database, s.casbinEnforcer, s.logger),
		adminuser.NewUpdateUserRoleHandler(s.database, s.casbinEnforcer, s.logger),
	)
	userRouter := userroute.NewRouter(
		engine,
		usercontroller.NewRegisterHandler(s.database, s.authenticator, s.logger),
		usercontroller.NewLoginHandler(s.database, s.authenticator, s.logger),
		usercontroller.NewGetMeHandler(s.database, s.logger),
		usercontroller.NewUpdateHandler(s.database, s.logger),
		usercontroller.NewUpdatePasswordHandler(s.database, s.logger),
		authentication,
	)
	routers := []route.Router{
		route.NewApiRouter(engine, system.NewPingHandler(), adminRouter, userRouter),
		route.NewSwaggerRouter(engine, system.NewSwaggerHandler(s.booterConfig)),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	newHttpApis := []model.HttpApi{}
	doesNotInsert := true
	for _, routeInfo := range engine.Routes() {
		for _, httpApi := range httpApis {
			if routeInfo.Method == httpApi.Method && routeInfo.Path == httpApi.Path {
				doesNotInsert = false
				break
			}
		}
		if doesNotInsert {
			newHttpApis = append(newHttpApis, model.HttpApi{Method: routeInfo.Method, Path: routeInfo.Path})
		}
		doesNotInsert = true
	}

	if len(newHttpApis) == 0 {
		return nil
	}

	if err := pkgpermission.CreateHttpApi(tx, &newHttpApis); err != nil {
		return err
	}

	return nil
}
