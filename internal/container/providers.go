package container

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/internal/domain/repository"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/internal/infrastructure/database"
	"mk-pos-billing/internal/service"
)

// ProvideCacheService constructs the cache.Service correctly merging prefixes.
func ProvideCacheService(redisClient *redis.Client, cfg config.CacheConfig, redisCfg config.RedisConfig) *cache.Service {
	return cache.NewCacheService(redisClient, redisCfg.Prefix+cfg.Prefix)
}

// ConfigProviderSet groups configuration providers
var ConfigProviderSet = wire.NewSet(
	config.LoadDatabaseConfig,
	config.LoadRedisConfig,
	config.LoadCacheConfig,
	config.LoadSalesInvoiceConfig,
	config.LoadCachePrefixes,
	config.LoadCacheTags,
	config.LoadCacheExpiry,
)

// InfrastructureProviderSet groups DB, Redis, and cache service providers
var InfrastructureProviderSet = wire.NewSet(
	database.ProvideDB,
	database.ProvideRedisClient,
	ProvideCacheService,
)

// CacheProviderSet groups domain-specific cache providers
var CacheProviderSet = wire.NewSet(
	domaincache.NewDeviceCache,
	domaincache.NewProductCache,
	domaincache.NewStoreCache,
	domaincache.NewTillCache,
	domaincache.NewUserAuthCache,
	domaincache.NewUserCache,
)

// RepositoryProviderSet groups all repositories
var RepositoryProviderSet = wire.NewSet(
	repository.NewSalesInvoiceRepository,
	repository.NewSalesInvoiceDraftRepository,
	repository.NewStoreInventoryRepository,
	repository.NewMasterDataRepository,
	repository.NewProductRepository,
	repository.NewTransactionManager,
)

// ServiceProviderSet groups business logic services
var ServiceProviderSet = wire.NewSet(
	service.NewSalesInvoiceService,
	service.NewCacheMasterService,
)

// HandlerProviderSet groups all API handlers
var HandlerProviderSet = wire.NewSet(
	handlers.NewSalesInvoiceHandler,
	handlers.NewCacheTestHandler,
)

// MiddlewareProviderSet groups all API middlewares
var MiddlewareProviderSet = wire.NewSet(
	middleware.NewDuplicateRequestMiddleware,
	middleware.NewDeviceTokenValidateMiddleware,
	middleware.NewPOSAuthTokenValidateMiddleware,
	middleware.NewMapTillMiddleware,
	middleware.NewSPOSCheckPermissionsMiddleware,
	middleware.NewCheckTillStatusMiddleware,
)

// ApplicationSet is a convenience set for the entire application dependency graph
var ApplicationSet = wire.NewSet(
	ConfigProviderSet,
	InfrastructureProviderSet,
	CacheProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
	MiddlewareProviderSet,
)
