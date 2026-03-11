package test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/logger"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type registryConfigParams struct {
	fx.In

	BooterConfig         *booter.Config
	HttpServerConfig     *http.ServerConfig          `optional:"true"`
	CsrfConfig           *middleware.CsrfConfig      `optional:"true"`
	RateLimitConfig      *middleware.RateLimitConfig `optional:"true"`
	AuthenticationConfig *authentication.Config      `optional:"true"`
	DatabaseConfig       *database.Config            `optional:"true"`
	RedisConfig          *redis.Config               `optional:"true"`
	ClickhouseConfig     *clickhouse.Config          `optional:"true"`
	ConsoleLoggerConfig  *logger.Config              `name:"logger.console" optional:"true"`
	FileLoggerConfig     *logger.Config              `name:"logger.file" optional:"true"`
	AccessLoggerConfig   *logger.Config              `name:"logger.access" optional:"true"`
}

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
			registerTestConfigDependencies,
			registrar.RegisterServiceDependencies,
			registrar.RegisterValidator,
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

func registerTestConfigDependencies(params registryConfigParams) {
	current := deps.CurrentConfig()
	current.BooterConfig = params.BooterConfig
	if params.CsrfConfig != nil {
		current.CsrfConfigValue = *params.CsrfConfig
	}
	if params.RateLimitConfig != nil {
		current.RateLimitConfigValue = *params.RateLimitConfig
	}
	current.AuthenticationConfig = params.AuthenticationConfig
	current.DatabaseConfig = params.DatabaseConfig
	current.RedisConfig = params.RedisConfig
	current.ClickhouseConfig = params.ClickhouseConfig
	current.ConsoleLoggerConfig = params.ConsoleLoggerConfig
	current.FileLoggerConfig = params.FileLoggerConfig
	current.AccessLoggerConfig = params.AccessLoggerConfig
	deps.SetConfig(current)
}

func emptyMockedServices() {
	// If you register a new mock dependency, initialize its empty test value here.
}
