package test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type Runtime struct {
	app         *fxtest.App
	options     RuntimeOptions
	HTTP        *httpHandler
	Rdbms       *rdbmsMigration
	Clickhouse  *clickhouseMigration
	Users       *UserOperator
	Permissions *PermissionOperator
}

type RuntimeOptions struct {
	UseRdbms      bool
	UseClickhouse bool
}

func NewRuntime(tb testing.TB, options RuntimeOptions) *Runtime {
	tb.Helper()

	wd, envFile, configFile := testConfigFiles()
	loadEnv(wd, envFile)

	booterConfig := booter.NewConfig(wd, configFile, false)

	app := fxtest.New(
		tb,
		fx.StopTimeout(60*time.Second),
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
			registrar.NewMapstructureDecoder,
			form.NewDecoder,
			modifiers.New,
		),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(registrar.ValidatorOnStart),
			),
			registrar.RegisterConfigDependencies,
			registrar.RegisterServiceDependencies,
		),
	)
	app.RequireStart()

	emptyMockedServices()
	httpHandler := NewHttpHandler()

	rt := &Runtime{
		app:         app,
		options:     options,
		HTTP:        httpHandler,
		Rdbms:       NewRdbmsMigration(),
		Clickhouse:  NewClickhouseMigration(),
		Users:       NewUserOperator(httpHandler),
		Permissions: NewPermissionOperator(),
	}

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

func testConfigFiles() (string, string, string) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}

	wd := path.Join(path.Dir(file), "../..")
	env := "test"
	if e := os.Getenv("ENV"); e != "" {
		switch e {
		case "test":
			env = e
		default:
			panic(fmt.Sprintf("unsupported ENV for test setup: %s", e))
		}
	}

	return wd, fmt.Sprintf(".env.%s", env), fmt.Sprintf("config.%s.yml", env)
}

func loadEnv(wd string, envFile string) {
	if err := godotenv.Load(path.Join(wd, envFile)); err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
}

func emptyMockedServices() {
	// If you register a new mock dependency, initialize its empty test value here.
}
