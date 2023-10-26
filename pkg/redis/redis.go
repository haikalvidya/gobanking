package redis

import (
	"context"
	"gobanking/pkg/serializer"
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

func GetDataFromRedis[T any](ctx context.Context, redisClient *redis.Client, key string, result *T) (*T, error) {
	redisResp := redisClient.Get(ctx, key)
	if redisResp.Err() != nil || redisResp.Err() == redis.Nil {
		return nil, redisResp.Err()
	}

	if err := serializer.Unmarshal([]byte(redisResp.Val()), &result); err != nil {
		return nil, err
	}

	return result, nil
}

const (
	RefreshToken = "refresh_token"
	Token        = "token"
)

func GetKeyOfTokenUserFromRedis(token string, theType string) string {
	// 7 last char of token
	return theType + ":user:" + token[len(token)-7:]
}
