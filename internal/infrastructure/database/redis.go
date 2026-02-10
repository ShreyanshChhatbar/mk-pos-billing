package database

import (
	"context"
	"mk-pos-billing/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ProvideRedisClient(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		zap.L().Error("Failed to connect to Redis", zap.Error(err))
		panic(err)
	} else {
		zap.L().Info("Connected to Redis", zap.String("addr", cfg.Addr))
	}

	return rdb
}
