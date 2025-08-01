package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "all env vars set",
			envVars: map[string]string{
				"PORT":                 "8080",
				"LOG_LEVEL":            "info",
				"USERS_SERVICE_URL":    "http://users:8081",
				"PRODUCTS_SERVICE_URL": "http://products:8082",
			},
			expected: Config{
				Port:               "8080",
				LogLevel:           "info",
				UsersServiceURL:    "http://users:8081",
				ProductsServiceURL: "http://products:8082",
			},
		},
		{
			name:    "default values when env var missing",
			envVars: map[string]string{},
			expected: Config{
				Port:               "8080",
				LogLevel:           "info",
				UsersServiceURL:    "http://localhost:8081",
				ProductsServiceURL: "http://localhost:8082",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			config := LoadFromEnv()

			// Assert
			if config.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", config.Port, tt.expected.Port)
			}
			if config.LogLevel != tt.expected.LogLevel {
				t.Errorf("LogLevel = %v, want %v", config.LogLevel, tt.expected.LogLevel)
			}
			if config.UsersServiceURL != tt.expected.UsersServiceURL {
				t.Errorf("UsersServiceURL = %v, want %v", config.UsersServiceURL, tt.expected.UsersServiceURL)
			}
			if config.ProductsServiceURL != tt.expected.ProductsServiceURL {
				t.Errorf("ProductsServiceURL = %v, want %v", config.ProductsServiceURL, tt.expected.ProductsServiceURL)
			}
		})
	}
}
