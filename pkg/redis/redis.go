package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	MinIdleConns int    `yaml:"minIdleConns"`
	PoolSize     int    `yaml:"poolSize"`
	PoolTimeout  int    `yaml:"poolTimeout"`
}

func NewRedisConn(ctx *context.Context, cfg *Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})

	if err := redisClient.Ping(*ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
