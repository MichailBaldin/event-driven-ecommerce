package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"services/gateway/internal/config"
	"services/gateway/internal/metrics"
	"services/gateway/internal/router"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
)

func TestGateway_ServeHTTP_UsersRequest(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetricsWithRegistry(registry)

	cfg := &config.Config{
		Port:               "8080",
		LogLevel:           "info",
		UsersServiceURL:    "http://users:8081",
		ProductsServiceURL: "http://products:8082",
	}

	r := router.NewRouter()
	logger := zap.NewNop()

	gateway := NewGatewayWithMetrics(cfg, r, logger, metrics)

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	gateway.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["method"] != "GET" {
		t.Errorf("Expected method GET, got %v", response["method"])
	}

	if response["path"] != "/users/123" {
		t.Errorf("Expected path /users/123, got %v", response["path"])
	}

	if response["target_service"] != "users" {
		t.Errorf("Expected target_service users, got %v", response["target_service"])
	}
}

func TestGateway_ServeHTTP_ProductsRequest(t *testing.T) {
	cfg := &config.Config{
		Port:               "8080",
		LogLevel:           "info",
		UsersServiceURL:    "http://users:8081",
		ProductsServiceURL: "http://products:8082",
	}

	r := router.NewRouter()
	logger := zap.NewNop()

	gateway := NewGateway(cfg, r, logger)

	req := httptest.NewRequest("POST", "/products/search", nil)
	w := httptest.NewRecorder()

	gateway.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["target_service"] != "products" {
		t.Errorf("Expected target_service products, got %v", response["target_service"])
	}
}

func TestGateway_ServeHTTP_UnknownPath(t *testing.T) {
	cfg := &config.Config{}
	r := router.NewRouter()
	logger := zap.NewNop()

	gateway := NewGateway(cfg, r, logger)

	req := httptest.NewRequest("GET", "/unknown/path", nil)
	w := httptest.NewRecorder()

	gateway.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
