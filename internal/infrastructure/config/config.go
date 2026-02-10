package config

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	DSN string
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DSN: os.Getenv("DB_DSN"),
	}
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadRedisConfig() RedisConfig {
	db := 0
	if val := os.Getenv("REDIS_DB"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			db = i
		}
	}
	return RedisConfig{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}

type AsynqRedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadAsynqRedisConfig() AsynqRedisConfig {
	db := 1
	if val := os.Getenv("ASYNQ_REDIS_DB"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			db = i
		}
	}
	return AsynqRedisConfig{
		Addr:     os.Getenv("ASYNQ_REDIS_HOST") + ":" + os.Getenv("ASYNQ_REDIS_PORT"),
		Password: os.Getenv("ASYNQ_REDIS_PASSWORD"),
		DB:       db,
	}
}

func envStringOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
