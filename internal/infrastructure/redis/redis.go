package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfing struct {
	Host string
	Port int
}

func Init(cfg RedisConfing) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return rc
}
