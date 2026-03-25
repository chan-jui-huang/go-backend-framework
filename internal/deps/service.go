package deps

import (
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4"
	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceState struct {
	LoggerValue         *zap.Logger
	ConsoleLogger       *zap.Logger
	FileLogger          *zap.Logger
	AccessLoggerValue   *zap.Logger
	DatabaseValue       *gorm.DB
	RedisValue          *redisClient.Client
	AuthenticatorValue  *authentication.Authenticator
	CasbinEnforcerValue *casbin.SyncedCachedEnforcer
	MapstructureDecoder func(any, any) error
	ClickhouseValue     driver.Conn
	FormDecoder         *form.Decoder
	Modifier            *mold.Transformer
}

var (
	serviceMu    sync.RWMutex
	serviceState ServiceState
)

func SetService(next ServiceState) {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	serviceState = next
}

func CurrentService() ServiceState {
	serviceMu.RLock()
	defer serviceMu.RUnlock()
	return serviceState
}

func Logger() *zap.Logger {
	return CurrentService().LoggerValue
}

func AccessLogger() *zap.Logger {
	return CurrentService().AccessLoggerValue
}

func Database() *gorm.DB {
	return CurrentService().DatabaseValue
}

func Redis() *redisClient.Client {
	return CurrentService().RedisValue
}

func Authenticator() *authentication.Authenticator {
	return CurrentService().AuthenticatorValue
}

func CasbinEnforcer() *casbin.SyncedCachedEnforcer {
	return CurrentService().CasbinEnforcerValue
}

func MapstructureDecoder() func(any, any) error {
	return CurrentService().MapstructureDecoder
}

func Clickhouse() driver.Conn {
	return CurrentService().ClickhouseValue
}
