package models

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr string
	Pass string
	Db   uint
}

func (cfg *RedisConfig) Default() RedisConfig {
	return RedisConfig{
		Addr: fmt.Sprintf("%s:%s", Env(REDIS_HOST), Env(REDIS_PORT)),
		Pass: Env(REDIS_PASSWORD),
		Db:   0,
	}
}

func (cfg *RedisConfig) Start() *redis.Client {
	cfg.Default()
	redis_client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pass,
		DB:       int(cfg.Db),
	})
	status := redis_client.Ping(context.TODO())
	fmt.Println(status)
	return redis_client
}
