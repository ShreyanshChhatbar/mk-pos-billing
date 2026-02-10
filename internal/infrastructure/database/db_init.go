package database

import (
	"log"
	"mk-pos-billing/internal/infrastructure/config"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func ProvideDB(cfg config.DatabaseConfig) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: buildGormLogger(),
	})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
		return nil, err
	}
	return &DB{DB: applyGlobalScopes(db)}, nil
}

func buildGormLogger() gormLogger.Interface {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	switch env {
	case "production", "prod", "staging":
		return gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  gormLogger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	default:
		return gormLogger.Default.LogMode(gormLogger.Info)
	}
}
