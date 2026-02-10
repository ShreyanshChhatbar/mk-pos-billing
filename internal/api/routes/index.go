package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterAllRoutes(r *gin.Engine, db *database.DB, redisClient *redis.Client) {
	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Public routes
	// Add your application routes here

	// Protected routes
	// protected := r.Group("/api/v1")
	// protected.Use(authMW.Authenticate())
	// protected.GET("/protected", handlers.SomeProtectedHandler)
}
