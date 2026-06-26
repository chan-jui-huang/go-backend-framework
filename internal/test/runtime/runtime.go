package runtime

import (
	"testing"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	testconfig "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/config"
	dbfixture "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/fixture/db"
	domainfixture "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/fixture/domain"
	httpfixture "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/fixture/http"
	scenariofixture "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/fixture/scenario"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
)

type Runtime struct {
	app         *fxtest.App
	options     RuntimeOptions
	HTTP        *httpfixture.Handler
	Rdbms       *dbfixture.RdbmsMigration
	Clickhouse  *dbfixture.ClickhouseMigration
	Users       *domainfixture.UserFixture
	Permissions *domainfixture.PermissionFixture
	UserAPI     *scenariofixture.UserAPI
	AdminAPI    *scenariofixture.AdminAPI
}

type RuntimeOptions struct {
	UseRdbms      bool
	UseClickhouse bool
}

func NewRuntime(tb testing.TB, options RuntimeOptions) *Runtime {
	tb.Helper()

	files := testconfig.NewFiles("../../..")
	testconfig.LoadEnv(files)

	booterConfig := booter.NewConfig(files.WorkDir, files.ConfigFile, false)
	rt := &Runtime{
		options: options,
	}

	app := fxtest.New(
		tb,
		fx.StopTimeout(60*time.Second),
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booterConfig),
		fx.Provide(
			registrar.NewConfigLoader,
			registrar.NewHttpServerConfig,
			registrar.NewCsrfConfig,
			registrar.NewRateLimitConfig,
			registrar.NewAuthenticationConfig,
			registrar.NewAuthenticator,
			registrar.NewDatabaseConfig,
			registrar.NewDatabase,
			registrar.NewRedisConfig,
			registrar.NewRedis,
			registrar.NewClickhouseConfig,
			registrar.NewClickhouse,
			registrar.NewLoggerConfigs,
			registrar.NewLoggers,
			registrar.NewCasbinEnforcer,
			form.NewDecoder,
			modifiers.New,
			httpfixture.NewEngine,
			dbfixture.NewRdbmsMigration,
			dbfixture.NewClickhouseMigration,
			domainfixture.NewUserFixture,
			domainfixture.NewPermissionFixture,
			httpfixture.New,
			scenariofixture.NewUserAPI,
			scenariofixture.NewAdminAPI,
		),
		route.NewModule(),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(registrar.ValidatorOnStart),
			),
		),
		fx.Populate(
			&rt.HTTP,
			&rt.Rdbms,
			&rt.Clickhouse,
			&rt.Users,
			&rt.Permissions,
			&rt.UserAPI,
			&rt.AdminAPI,
		),
	)
	app.RequireStart()

	rt.app = app

	if rt.options.UseRdbms {
		rt.Rdbms.Run()
	}

	if rt.options.UseClickhouse {
		rt.Clickhouse.Run()
	}

	return rt
}

func NewBaseRuntime(tb testing.TB) *Runtime {
	tb.Helper()

	return NewRuntime(tb, RuntimeOptions{})
}

func NewRdbmsRuntime(tb testing.TB) *Runtime {
	tb.Helper()

	return NewRuntime(tb, RuntimeOptions{UseRdbms: true})
}

func NewClickhouseRuntime(tb testing.TB) *Runtime {
	tb.Helper()

	return NewRuntime(tb, RuntimeOptions{UseClickhouse: true})
}

func NewFullRuntime(tb testing.TB) *Runtime {
	tb.Helper()

	return NewRuntime(tb, RuntimeOptions{
		UseRdbms:      true,
		UseClickhouse: true,
	})
}

func (rt *Runtime) Close() {
	if rt == nil {
		return
	}

	if rt.options.UseClickhouse && rt.Clickhouse != nil {
		rt.Clickhouse.Reset()
	}

	if rt.options.UseRdbms && rt.Rdbms != nil {
		rt.Rdbms.Reset()
	}

	if rt.app != nil {
		rt.app.RequireStop()
		rt.app = nil
	}
}
