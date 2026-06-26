package registrar_test

import (
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	testconfig "github.com/chan-jui-huang/go-backend-framework/v3/internal/test/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegistrarTestSuite struct {
	suite.Suite
	app          *fxtest.App
	dependencies registrarDependencies
}

type registrarDependencies struct {
	BooterConfig     *booter.Config
	CsrfConfig       *config.CsrfConfig
	DatabaseConfig   *database.Config
	ClickhouseConfig *clickhouse.Config
	Database         *gorm.DB
	Redis            *redisClient.Client
	Logger           *zap.Logger
	Authenticator    *authentication.Authenticator
	CasbinEnforcer   *casbin.SyncedCachedEnforcer
	ClickhouseConn   driver.Conn
}

func (suite *RegistrarTestSuite) SetupSuite() {
	files := testconfig.NewFiles("../..")
	testconfig.LoadEnv(files)

	suite.app = fxtest.New(
		suite.T(),
		fx.StopTimeout(60*time.Second),
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booter.NewConfig(files.WorkDir, files.ConfigFile, false)),
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
		),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(registrar.ValidatorOnStart),
			),
		),
		fx.Populate(
			&suite.dependencies.BooterConfig,
			&suite.dependencies.CsrfConfig,
			&suite.dependencies.DatabaseConfig,
			&suite.dependencies.ClickhouseConfig,
			&suite.dependencies.Database,
			&suite.dependencies.Redis,
			&suite.dependencies.Logger,
			&suite.dependencies.Authenticator,
			&suite.dependencies.CasbinEnforcer,
			&suite.dependencies.ClickhouseConn,
		),
	)
	suite.app.RequireStart()
}

func (suite *RegistrarTestSuite) TestDependenciesRegistered() {
	suite.NotNil(suite.dependencies.BooterConfig)
	suite.NotNil(suite.dependencies.CsrfConfig)
	suite.NotNil(suite.dependencies.DatabaseConfig)
	suite.NotNil(suite.dependencies.ClickhouseConfig)
	suite.NotNil(suite.dependencies.Authenticator)
	suite.NotNil(suite.dependencies.Database)
	suite.NotNil(suite.dependencies.Redis)
	suite.NotNil(suite.dependencies.ClickhouseConn)
	suite.NotNil(suite.dependencies.CasbinEnforcer)
	suite.NotNil(suite.dependencies.Logger)
}

func (suite *RegistrarTestSuite) TestValidatorTagNameFunc() {
	type dummy struct {
		Email string `json:"email" binding:"required,email"`
		Page  int    `form:"page" binding:"required,gt=0"`
	}

	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		suite.Fail("validator engine is not *validator.Validate")
		return
	}

	err := engine.Struct(&dummy{})
	suite.Error(err)

	ves, ok := err.(validator.ValidationErrors)
	if !ok {
		suite.Fail("error is not ValidationErrors")
		return
	}

	fields := map[string]string{}
	for _, fe := range ves {
		fields[fe.Field()] = fe.Tag()
	}

	suite.Equal(map[string]string{
		"email": "required",
		"page":  "required",
	}, fields)
}

func (suite *RegistrarTestSuite) TearDownSuite() {
	suite.app.RequireStop()
}

func TestRegistrarTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrarTestSuite))
}
