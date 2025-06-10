package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger instance
func NewLogger() (*zap.Logger, error) {
	// Get environment
	env := os.Getenv("SERVER_MODE")
	
	var config zap.Config
	
	if env == "production" {
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	// Common configurations
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	
	// Add common fields
	config.InitialFields = map[string]interface{}{
		"service": "oms-service",
	}
	
	// Build logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	
	// Replace global logger
	zap.ReplaceGlobals(logger)
	
	return logger, nil
}

// With creates a child logger with additional fields
func With(logger *zap.Logger, fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// WithContext creates a logger with context values
func WithContext(logger *zap.Logger, ctx map[string]interface{}) *zap.Logger {
	fields := make([]zap.Field, 0, len(ctx))
	for k, v := range ctx {
		fields = append(fields, zap.Any(k, v))
	}
	return logger.With(fields...)
}