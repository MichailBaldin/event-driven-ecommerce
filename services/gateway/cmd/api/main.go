// services/gateway/cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"services/gateway/internal/config"
	"services/gateway/internal/gateway"
	"services/gateway/internal/router"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := config.LoadFromEnv()

	logger, err := createLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Gateway service",
		zap.String("port", cfg.Port),
		zap.String("log_level", cfg.LogLevel),
		zap.String("users_service_url", cfg.UsersServiceURL),
		zap.String("products_service_url", cfg.ProductsServiceURL),
	)

	r := router.NewRouter()
	gw := gateway.NewGateway(&cfg, r, logger)

	http.Handle("/", gw)

	addr := ":" + cfg.Port
	logger.Info("Gateway server starting", zap.String("address", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func createLogger(level string) (*zap.Logger, error) {
	var config zap.Config

	switch level {
	case "debug":
		config = zap.NewDevelopmentConfig()
	case "info", "warn", "error":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(level))
	default:
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
