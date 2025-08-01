package router

import (
	"fmt"
	"strings"
)

// Router определяет target сервис по HTTP пути
type Router struct {
}

// NewRouter создает новый Router
func NewRouter() *Router {
	return &Router{}
}

// Route определяет имя сервиса по пути запроса
func (r *Router) Route(path string) (string, error) {
	if path == "" || path == "/" {
		return "", fmt.Errorf("empty or root path not supported")
	}

	path = strings.TrimSuffix(path, "/")

	switch {
	case strings.HasPrefix(path, "/users"):
		return "users", nil
	case strings.HasPrefix(path, "/products"):
		return "products", nil
	default:
		return "", fmt.Errorf("unknown service for path: %s", path)
	}
}
