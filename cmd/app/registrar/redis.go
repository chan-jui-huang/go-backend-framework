package registrar

import (
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/redis"
	redisClient "github.com/redis/go-redis/v9"
)

func NewRedisConfig(loader *booterconfig.Loader) *redis.Config {
	config := &redis.Config{}
	loader.Unmarshal("redis", config)

	return config
}

func NewRedis(config *redis.Config) *redisClient.Client {
	return redis.New(*config)
}
