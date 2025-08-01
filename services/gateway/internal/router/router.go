package router

import (
	"fmt"
	"strings"
)

// Router определяет target сервис по HTTP пути
type Router struct {
	// пока простая реализация без настроек
}

// NewRouter создает новый Router
func NewRouter() *Router {
	return &Router{}
}

// Route определяет имя сервиса по пути запроса
func (r *Router) Route(path string) (string, error) {
	// Убираем пустые пути
	if path == "" || path == "/" {
		return "", fmt.Errorf("empty or root path not supported")
	}

	// Убираем trailing slash для единообразия
	path = strings.TrimSuffix(path, "/")

	// Определяем сервис по префиксу пути
	switch {
	case strings.HasPrefix(path, "/users"):
		return "users", nil
	case strings.HasPrefix(path, "/products"):
		return "products", nil
	default:
		return "", fmt.Errorf("unknown service for path: %s", path)
	}
}
