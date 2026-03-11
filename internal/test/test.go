package test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
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

var (
	testApp *fxtest.App
	setupMu sync.Mutex
)

func Setup(tb testing.TB) {
	tb.Helper()

	setupMu.Lock()
	defer setupMu.Unlock()

	if testApp != nil {
		return
	}

	wd, envFile, configFile := testConfigFiles()
	loadEnv(wd, envFile)

	booterConfig := booter.NewConfig(wd, configFile, false)

	testApp = fxtest.New(
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
	testApp.RequireStart()

	emptyMockedServices()

	HttpHandler = NewHttpHandler()
	RdbmsMigration = NewRdbmsMigration()
	ClickhouseMigration = NewClickhouseMigration()
	PermissionService = NewPermissionService()
	UserService = NewUserService()
	AdminService = NewAdminService()
}

func Shutdown() {
	setupMu.Lock()
	defer setupMu.Unlock()

	if testApp != nil {
		testApp.RequireStop()
		testApp = nil
	}

	HttpHandler = nil
	RdbmsMigration = nil
	ClickhouseMigration = nil
	PermissionService = nil
	UserService = nil
	AdminService = nil
}

func testConfigFiles() (string, string, string) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}

	wd := path.Join(path.Dir(file), "../..")
	env := "dev"
	if e := os.Getenv("ENV"); e != "" {
		env = e
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
