package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	"services/gateway/internal/config"
	"services/gateway/internal/router"

	"go.uber.org/zap"
)

type Gateway struct {
	config *config.Config
	router *router.Router
	logger *zap.Logger
}

type GatewayResponse struct {
	Message       string `json:"message"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	TargetService string `json:"target_service"`
	Timestamp     string `json:"timestamp"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Path  string `json:"path"`
	Code  int    `json:"code"`
}

func NewGateway(cfg *config.Config, r *router.Router, logger *zap.Logger) *Gateway {
	return &Gateway{
		config: cfg,
		router: r,
		logger: logger,
	}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.logger.Info("Gateway request received",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)

	targetService, err := g.router.Route(r.URL.Path)
	if err != nil {
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

	g.logger.Info("Request routed successfully",
		zap.String("path", r.URL.Path),
		zap.String("target_service", targetService),
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
