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

	cachePrefixesCfg := cacheCfg.Prefixes
	cacheTagsCfg := cacheCfg.Tags
	cacheExpiryCfg := cacheCfg.Expiry
	cacheMasterService := service.NewCacheMasterService(cacheService, cachePrefixesCfg, cacheTagsCfg, cacheExpiryCfg)

	salesInvoiceRepo := repository.NewSalesInvoiceRepository(db)
	draftRepo := repository.NewSalesInvoiceDraftRepository(db)
	inventoryRepo := repository.NewStoreInventoryRepository(db)
	masterRepo := repository.NewMasterDataRepository(db)
	productRepo := repository.NewProductRepository(db)
	txManager := repository.NewTransactionManager(db)
	salesInvoiceService := service.NewSalesInvoiceService(
		salesInvoiceRepo,
		draftRepo,
		inventoryRepo,
		productRepo,
		masterRepo,
		cacheService,
		salesCfg,
		cacheMasterService,
		txManager,
	)
	salesInvoiceHandler := handlers.NewSalesInvoiceHandler(salesInvoiceService, salesCfg)
	duplicateRequestMiddleware := middleware.NewDuplicateRequestMiddleware(cacheService, salesCfg)

	deviceTokenValidateMiddleware := middleware.NewDeviceTokenValidateMiddleware(cacheMasterService, cachePrefixesCfg)
	posAuthTokenValidateMiddleware := middleware.NewPOSAuthTokenValidateMiddleware(cacheMasterService)
	mapTillMiddleware := middleware.NewMapTillMiddleware(cacheMasterService)
	sposCheckPermissionsMiddleware := middleware.NewSPOSCheckPermissionsMiddleware(cacheMasterService)
	checkTillStatusMiddleware := middleware.NewCheckTillStatusMiddleware(cacheMasterService)
	cacheTestHandler := handlers.NewCacheTestHandler(cacheMasterService)

	return &ServerApp{
		SalesInvoiceHandler:            salesInvoiceHandler,
		DuplicateRequestMiddleware:     duplicateRequestMiddleware,
		DeviceTokenValidateMiddleware:  deviceTokenValidateMiddleware,
		POSAuthTokenValidateMiddleware: posAuthTokenValidateMiddleware,
		MapTillMiddleware:              mapTillMiddleware,
		CacheMasterService:             cacheMasterService,
		SPOSCheckPermissionsMiddleware: sposCheckPermissionsMiddleware,
		CheckTillStatusMiddleware:      checkTillStatusMiddleware,
		CacheTestHandler:               cacheTestHandler,
	}, nil
}
