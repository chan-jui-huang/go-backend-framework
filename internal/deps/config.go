package deps

import (
	"sync"

	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/logger"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/redis"
)

type ConfigState struct {
	BooterConfig         *booter.Config
	CsrfConfigValue      any
	RateLimitConfigValue any
	AuthenticationConfig *authentication.Config
	DatabaseConfig       *database.Config
	RedisConfig          *redis.Config
	ClickhouseConfig     *clickhouse.Config
	ConsoleLoggerConfig  *logger.Config
	FileLoggerConfig     *logger.Config
	AccessLoggerConfig   *logger.Config
}

var (
	configMu    sync.RWMutex
	configState ConfigState
)

func SetConfig(next ConfigState) {
	configMu.Lock()
	defer configMu.Unlock()
	configState = next
}

func CurrentConfig() ConfigState {
	configMu.RLock()
	defer configMu.RUnlock()
	return configState
}

func BooterConfig() booter.Config {
	return *CurrentConfig().BooterConfig
}

func CsrfConfig() any {
	return CurrentConfig().CsrfConfigValue
}

func RateLimitConfig() any {
	return CurrentConfig().RateLimitConfigValue
}

func DatabaseConfig() database.Config {
	return *CurrentConfig().DatabaseConfig
}

func ClickhouseConfig() clickhouse.Config {
	return *CurrentConfig().ClickhouseConfig
}
