package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	"services/gateway/internal/config"
	"services/gateway/internal/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"services/gateway/internal/metrics"
)

type Gateway struct {
	config  *config.Config
	router  *router.Router
	logger  *zap.Logger
	version string
	metrics *metrics.Metrics
}

type GatewayResponse struct {
	Message       string `json:"message" example:"Gateway received request"`
	Method        string `json:"method" example:"GET"`
	Path          string `json:"path" example:"/users/123"`
	TargetService string `json:"target_service" example:"users"`
	Timestamp     string `json:"timestamp" example:"2025-08-01T12:30:00Z"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"unknown service for path: /unknown/path"`
	Path  string `json:"path" example:"/unknown/path"`
	Code  int    `json:"code" example:"404"`
}

// HealthResponse структура для health check
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Service   string `json:"service" example:"gateway"`
	Version   string `json:"version" example:"v0.1.0"`
	Timestamp string `json:"timestamp" example:"2025-08-01T12:30:00Z"`
}

func NewGateway(cfg *config.Config, r *router.Router, logger *zap.Logger) *Gateway {
	return &Gateway{
		config:  cfg,
		router:  r,
		logger:  logger,
		version: "v0.1.0",
		metrics: metrics.NewMetrics(),
	}
}

// @title E-commerce Gateway API
// @version 1.0
// @description API Gateway for microservices e-commerce system. Routes requests to users and products services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://github.com/MichailBaldin/event-driven-ecommerce
// @contact.email michailbaldin@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// ServeHTTP handles HTTP requests and routes them to appropriate services
// @Summary Route request to microservice
// @Description Routes HTTP requests to users or products microservices based on path prefix
// @Tags routing
// @Accept json
// @Produce json
// @Param path path string true "Request path" example("/users/123")
// @Success 200 {object} GatewayResponse "Successfully routed request"
// @Failure 404 {object} ErrorResponse "Unknown service path"
// @Router /users/{id} [get]
// @Router /users/{id} [post]
// @Router /users/{id} [put]
// @Router /users/{id} [delete]
// @Router /products/{id} [get]
// @Router /products/{id} [post]
// @Router /products/{id} [put]
// @Router /products/{id} [delete]
// @Router /products/search [get]
// @Router /users/create [post]
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	g.metrics.IncActiveRequests()

	defer g.metrics.DecActiveRequests()

	g.logger.Info("Gateway request received",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)

	targetService, err := g.router.Route(r.URL.Path)
	if err != nil {
		duration := time.Since(start)
		g.metrics.RecordRequest(r.Method, r.URL.Path, "404", "none", duration)

		g.logger.Error("Failed to route request",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		errorResp := ErrorResponse{
			Error: err.Error(),
			Path:  r.URL.Path,
			Code:  http.StatusNotFound,
		}

		json.NewEncoder(w).Encode(errorResp)
		return
	}

	duration := time.Since(start)
	g.metrics.RecordRequest(r.Method, r.URL.Path, "200", targetService, duration)

	g.logger.Info("Request routed successfully",
		zap.String("path", r.URL.Path),
		zap.String("target_service", targetService),
		zap.Duration("duration", duration),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := GatewayResponse{
		Message:       "Gateway received request",
		Method:        r.Method,
		Path:          r.URL.Path,
		TargetService: targetService,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Health проверка состояния Gateway
// @Summary Health check
// @Description Returns the health status of the Gateway service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "Gateway is healthy"
// @Router /health [get]
func (g *Gateway) Health(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug("Health check requested")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:    "ok",
		Service:   "gateway",
		Version:   g.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Metrics возвращает Prometheus метрики
// @Summary Prometheus metrics
// @Description Returns Prometheus metrics for monitoring
// @Tags monitoring
// @Produce text
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (g *Gateway) Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
