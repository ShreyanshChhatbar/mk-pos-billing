package main

import (
	"mk-pos-billing/internal/api/routes"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/internal/infrastructure/database"
	"mk-pos-billing/internal/infrastructure/logger"
	"os"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Set default timezone globally (Asia/Kolkata)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	time.Local = loc

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	appEnv := os.Getenv("APP_ENV")
	sentryDSN := os.Getenv("SENTRY_DSN")

	// Initialize logger with Sentry
	logger.Init(appEnv, sentryDSN)
	defer logger.Sync()

	// Redirect Gin logs to Zap
	gin.DefaultWriter = zap.NewStdLog(zap.L()).Writer()

	// Run web server
	runServer()
}

func runServer() {
	// Initialize infrastructure
	dbCfg := config.LoadDatabaseConfig()
	db, err := database.ProvideDB(dbCfg)
	if err != nil {
		zap.L().Fatal("Failed to initialize database", zap.Error(err))
	}

	redisCfg := config.LoadRedisConfig()
	redisClient := database.ProvideRedisClient(redisCfg)

	// Gin setup
	r := gin.New()
	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(func(c *gin.Context) {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			if env := os.Getenv("APP_ENV"); env != "" {
				hub.Scope().SetTag("app_env", env)
			}
		}
		c.Next()
	})
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Log HTTP requests to console only in development
	if logger.ConsoleLoggingEnabled() {
		r.Use(gin.LoggerWithWriter(zap.NewStdLog(zap.L()).Writer()))
	}

	// Recover from panics, log stack traces
	r.Use(gin.RecoveryWithWriter(logger.Writer{}))

	// Register API routes
	routes.RegisterAllRoutes(r, db, redisClient)

	port := os.Getenv("PORT")
	zap.L().Info("Starting web server", zap.String("port", port), zap.String("env", os.Getenv("APP_ENV")))

	// Start server
	if err := r.Run("0.0.0.0:" + port); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
