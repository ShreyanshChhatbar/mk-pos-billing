package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestRoutes(r *gin.RouterGroup, CacheTestHandler *handlers.CacheTestHandler) {
	testGroup := r.Group("/test/cache")
	testGroup.GET("/till", CacheTestHandler.GetTillCache)
	testGroup.GET("/user-auth", CacheTestHandler.GetUserAuthCache)
	testGroup.GET("/store", CacheTestHandler.GetStoreCache)
	testGroup.GET("/product", CacheTestHandler.GetProductCache)
}