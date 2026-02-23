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
	cacheService := cache.NewCacheService(redisClient, "pos")
	salesCfg := config.LoadSalesInvoiceConfig()

	salesInvoiceRepo := repository.NewSalesInvoiceRepository(db)
	salesInvoiceService := service.NewSalesInvoiceService(salesInvoiceRepo, cacheService, salesCfg)
	salesInvoiceHandler := handlers.NewSalesInvoiceHandler(salesInvoiceService, salesCfg)
	duplicateRequestMiddleware := middleware.NewDuplicateRequestMiddleware(cacheService, salesCfg)

	return &ServerApp{
		SalesInvoiceHandler:        salesInvoiceHandler,
		DuplicateRequestMiddleware: duplicateRequestMiddleware,
	}, nil
}
