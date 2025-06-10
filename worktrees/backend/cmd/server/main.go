package main

import (
	"log"

	"github.com/openfoundry/oms/internal/config"
	"github.com/openfoundry/oms/internal/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger := initializeLogger()
	defer logger.Sync()

	cfg := loadConfiguration(logger)
	deps := initializeDependencies(cfg, logger)
	defer deps.Close()

	server := createHTTPServer(cfg, deps.Services, logger)
	startServer(server, logger)

	waitForShutdownSignal()
	shutdownServer(server, logger)
}

// initializeLogger initializes the application logger
func initializeLogger() *zap.Logger {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Logger initialized")
	return logger
}

// loadConfiguration loads application configuration
func loadConfiguration(logger *zap.Logger) *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded")
	return cfg
}

// initializeDependencies initializes all application dependencies
func initializeDependencies(cfg *config.Config, logger *zap.Logger) *Dependencies {
	logger.Info("Initializing application dependencies...")

	db, err := initializeDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	redisCache, err := initializeCache(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Redis cache", zap.Error(err))
	}

	kafkaPublisher := initializeMessaging(cfg, logger)

	deps := &Dependencies{
		DB:             db,
		RedisCache:     redisCache,
		KafkaPublisher: kafkaPublisher,
	}

	deps.Services = initializeServices(deps, logger)

	logger.Info("All dependencies initialized successfully")
	return deps
}