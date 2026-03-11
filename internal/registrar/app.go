package registrar

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/logger"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/redis"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4"
	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegisterConfigDependenciesParams struct {
	fx.In

	BooterConfig         *booter.Config
	CsrfConfig           *config.CsrfConfig      `optional:"true"`
	RateLimitConfig      *config.RateLimitConfig `optional:"true"`
	AuthenticationConfig *authentication.Config  `optional:"true"`
	DatabaseConfig       *database.Config        `optional:"true"`
	RedisConfig          *redis.Config           `optional:"true"`
	ClickhouseConfig     *clickhouse.Config      `optional:"true"`
	ConsoleLoggerConfig  *logger.Config          `name:"logger.console" optional:"true"`
	FileLoggerConfig     *logger.Config          `name:"logger.file" optional:"true"`
	AccessLoggerConfig   *logger.Config          `name:"logger.access" optional:"true"`
}

func RegisterConfigDependencies(params RegisterConfigDependenciesParams) {
	deps.SetConfig(deps.ConfigState{
		BooterConfig:         params.BooterConfig,
		CsrfConfigValue:      params.CsrfConfig,
		RateLimitConfigValue: params.RateLimitConfig,
		AuthenticationConfig: params.AuthenticationConfig,
		DatabaseConfig:       params.DatabaseConfig,
		RedisConfig:          params.RedisConfig,
		ClickhouseConfig:     params.ClickhouseConfig,
		ConsoleLoggerConfig:  params.ConsoleLoggerConfig,
		FileLoggerConfig:     params.FileLoggerConfig,
		AccessLoggerConfig:   params.AccessLoggerConfig,
	})
}

type RegisterServiceDependenciesParams struct {
	fx.In

	Database            *gorm.DB                      `optional:"true"`
	Redis               *redisClient.Client           `optional:"true"`
	Authenticator       *authentication.Authenticator `optional:"true"`
	CasbinEnforcer      *casbin.SyncedCachedEnforcer  `optional:"true"`
	MapstructureDecoder func(any, any) error          `optional:"true"`
	Clickhouse          driver.Conn                   `optional:"true"`
	FormDecoder         *form.Decoder                 `optional:"true"`
	Modifier            *mold.Transformer             `optional:"true"`
	Logger              *zap.Logger                   `optional:"true"`
	ConsoleLogger       *zap.Logger                   `name:"logger.console" optional:"true"`
	FileLogger          *zap.Logger                   `name:"logger.file" optional:"true"`
	AccessLogger        *zap.Logger                   `name:"logger.access" optional:"true"`
}

func RegisterServiceDependencies(params RegisterServiceDependenciesParams) {
	deps.SetService(deps.ServiceState{
		LoggerValue:         params.Logger,
		ConsoleLogger:       params.ConsoleLogger,
		FileLogger:          params.FileLogger,
		AccessLoggerValue:   params.AccessLogger,
		DatabaseValue:       params.Database,
		RedisValue:          params.Redis,
		AuthenticatorValue:  params.Authenticator,
		CasbinEnforcerValue: params.CasbinEnforcer,
		MapstructureDecoder: params.MapstructureDecoder,
		ClickhouseValue:     params.Clickhouse,
		FormDecoder:         params.FormDecoder,
		Modifier:            params.Modifier,
	})
}
