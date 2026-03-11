package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/redis"
)

type RedisRegistrar struct {
	config redis.Config
}

func (rr *RedisRegistrar) Boot() {
	config.Registry.Register("redis", &rr.config)
}

func (rr *RedisRegistrar) Register() {
	current := deps.CurrentConfig()
	current.RedisConfig = &rr.config
	deps.SetConfig(current)

	serviceState := deps.CurrentService()
	serviceState.RedisValue = redis.New(rr.config)
	deps.SetService(serviceState)
}
