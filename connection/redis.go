package connection

import (
	"fmt"
	"oms/config"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func RedisConnect(cfg config.Config) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})
}

func GetRedis(cfg config.Config) *redis.Client {
	if redisClient == nil {
		RedisConnect(cfg)
	}
	return redisClient
}
