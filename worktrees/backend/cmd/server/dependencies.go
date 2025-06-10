package main

import (
	"database/sql"
	"time"

	"github.com/openfoundry/oms/internal/config"
	"github.com/openfoundry/oms/internal/domain/service"
	"github.com/openfoundry/oms/internal/infrastructure/cache"
	"github.com/openfoundry/oms/internal/infrastructure/database"
	"github.com/openfoundry/oms/internal/infrastructure/messaging"
	"github.com/openfoundry/oms/internal/infrastructure/persistence/postgres"
	"github.com/openfoundry/oms/internal/interfaces/rest"
	"go.uber.org/zap"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	DB             *sql.DB
	RedisCache     *cache.RedisCache
	KafkaPublisher *messaging.KafkaPublisher
	Services       *rest.Services
}

// Close closes all dependencies
func (d *Dependencies) Close() error {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.RedisCache != nil {
		d.RedisCache.Close()
	}
	if d.KafkaPublisher != nil {
		d.KafkaPublisher.Close()
	}
	return nil
}

// initializeDatabase initializes database connection and runs migrations
func initializeDatabase(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	if err := database.RunMigrations(db, cfg.Database.MigrationsPath); err != nil {
		db.Close()
		return nil, err
	}

	logger.Info("Database initialized and migrations completed")
	return db, nil
}

// initializeCache initializes Redis cache
func initializeCache(cfg *config.Config, logger *zap.Logger) (*cache.RedisCache, error) {
	redisCache, err := cache.NewRedisCache(cache.RedisConfig{
		Addr:     cfg.Cache.RedisAddr,
		Password: cfg.Cache.RedisPassword,
		DB:       cfg.Cache.RedisDB,
		TTL:      time.Duration(cfg.Cache.TTL) * time.Second,
		Logger:   logger,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("Redis cache initialized")
	return redisCache, nil
}

// initializeMessaging initializes Kafka publisher
func initializeMessaging(cfg *config.Config, logger *zap.Logger) *messaging.KafkaPublisher {
	kafkaPublisher := messaging.NewKafkaPublisher(
		cfg.EventBus.KafkaBrokers,
		cfg.EventBus.KafkaTopic,
		logger,
	)

	logger.Info("Kafka publisher initialized")
	return kafkaPublisher
}

// initializeServices initializes all business services
func initializeServices(deps *Dependencies, logger *zap.Logger) *rest.Services {
	// Initialize repositories
	objectTypeRepo := postgres.NewObjectTypeRepository(deps.DB, logger)
	linkTypeRepo := postgres.NewLinkTypeRepository(deps.DB, logger)

	// Initialize caches
	objectTypeCache := cache.NewObjectTypeCache(deps.RedisCache)
	linkTypeCache := cache.NewLinkTypeCache(deps.RedisCache)

	// Initialize event publishers
	objectTypeEventPublisher := messaging.NewObjectTypeEventPublisher(deps.KafkaPublisher)
	linkTypeEventPublisher := messaging.NewLinkTypeEventPublisher(deps.KafkaPublisher)

	// Initialize services
	objectTypeService, err := service.NewObjectTypeService(service.ObjectTypeServiceConfig{
		Repository:     objectTypeRepo,
		Cache:          objectTypeCache,
		EventPublisher: objectTypeEventPublisher,
		Logger:         logger,
	})
	if err != nil {
		logger.Fatal("Failed to create object type service", zap.Error(err))
	}

	linkTypeService, err := service.NewLinkTypeService(service.LinkTypeServiceConfig{
		Repository:     linkTypeRepo,
		ObjectTypeRepo: objectTypeRepo,
		Cache:          linkTypeCache,
		EventPublisher: linkTypeEventPublisher,
		Logger:         logger,
	})
	if err != nil {
		logger.Fatal("Failed to create link type service", zap.Error(err))
	}

	logger.Info("All services initialized")
	return &rest.Services{
		ObjectTypeService: objectTypeService,
		LinkTypeService:   linkTypeService,
	}
}