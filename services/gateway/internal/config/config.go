package config

import "os"

// Config содержит настройки Gateway
type Config struct {
	Port               string
	LogLevel           string
	UsersServiceURL    string
	ProductsServiceURL string
}

// LoadFromEnv загружает конфигурацию из environment variables
func LoadFromEnv() Config {
	return Config{
		Port:               getEnvOrDefault("PORT", "8080"),
		LogLevel:           getEnvOrDefault("LOG_LEVEL", "info"),
		UsersServiceURL:    getEnvOrDefault("USERS_SERVICE_URL", "http://localhost:8081"),
		ProductsServiceURL: getEnvOrDefault("PRODUCTS_SERVICE_URL", "http://localhost:8082"),
	}
}

// getEnvOrDefault возвращает значение env переменной или default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
