package container

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"
	"mk-pos-billing/internal/domain/repository"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/internal/infrastructure/database"
	"mk-pos-billing/internal/service"
)

func InitializeServerApp() (*ServerApp, error) {
	dbCfg := config.LoadDatabaseConfig()
	db, err := database.ProvideDB(dbCfg)
	if err != nil {
		return nil, err
	}

	redisCfg := config.LoadRedisConfig()
	redisClient := database.ProvideRedisClient(redisCfg)
	cacheCfg := config.LoadCacheConfig()
	cacheService := cache.NewCacheService(redisClient, redisCfg.Prefix+cacheCfg.Prefix)
	salesCfg := config.LoadSalesInvoiceConfig()

	salesInvoiceRepo := repository.NewSalesInvoiceRepository(db)
	salesInvoiceService := service.NewSalesInvoiceService(salesInvoiceRepo, cacheService, salesCfg)
	salesInvoiceHandler := handlers.NewSalesInvoiceHandler(salesInvoiceService, salesCfg)
	duplicateRequestMiddleware := middleware.NewDuplicateRequestMiddleware(cacheService, salesCfg)
	
	cachePrefixesCfg := cacheCfg.Prefixes
	cacheTagsCfg := cacheCfg.Tags
	cacheExpiryCfg := cacheCfg.Expiry
	cacheMasterService := service.NewCacheMasterService(cacheService, cachePrefixesCfg, cacheTagsCfg, cacheExpiryCfg)
	// sposCheckPermissionsMiddleware := middleware.NewSPOSCheckPermissionsMiddleware(cacheMasterService)
	cacheTestHandler := handlers.NewCacheTestHandler(cacheMasterService)

	return &ServerApp{
		SalesInvoiceHandler:            salesInvoiceHandler,
		DuplicateRequestMiddleware:     duplicateRequestMiddleware,
		CacheMasterService:             cacheMasterService,
		// SPOSCheckPermissionsMiddleware: sposCheckPermissionsMiddleware,
		CacheTestHandler:               cacheTestHandler,
	}, nil
}
