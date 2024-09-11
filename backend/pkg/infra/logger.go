package infra

import (
	"os"

	"go.uber.org/zap"
)

// DevelopmentConfig returns a zap configuration tailored for development.
func DevelopmentConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		Encoding:          "console", // encodes output in a human-readable format directly to the console
		DisableStacktrace: true,
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(), // encoder config for development
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// ProductionConfig returns a zap configuration tailored for production.
func ProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json", // encodes output in structured JSON format
		DisableStacktrace: true,
		EncoderConfig:     zap.NewProductionEncoderConfig(), // encoder config for production
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// InitializeZapLogger initializes a zap logger based on the APP_ENV environment variable.
func InitializeZapLogger() {
	var config zap.Config
	switch env := os.Getenv("APP_ENV"); env {
	case "production":
		config = ProductionConfig()
	default:
		config = DevelopmentConfig() // default to development configuration if APP_ENV is not set or recognized
	}
	logger := zap.Must(config.Build())
	// replace global logger instance with the custom configured logger
	zap.ReplaceGlobals(logger)

	// ensure logger is flushed on exit
	defer logger.Sync()
}
