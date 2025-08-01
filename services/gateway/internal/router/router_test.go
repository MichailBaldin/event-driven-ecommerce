package router

import (
	"testing"
)

func TestRouter_Route(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantService string
		wantErr     bool
	}{
		{
			name:        "users path with ID",
			path:        "/users/123",
			wantService: "users",
			wantErr:     false,
		},
		{
			name:        "users root path",
			path:        "/users",
			wantService: "users",
			wantErr:     false,
		},
		{
			name:        "users with trailing slash",
			path:        "/users/",
			wantService: "users",
			wantErr:     false,
		},
		{
			name:        "products path with action",
			path:        "/products/search",
			wantService: "products",
			wantErr:     false,
		},
		{
			name:        "products root path",
			path:        "/products",
			wantService: "products",
			wantErr:     false,
		},
		{
			name:        "unknown path",
			path:        "/unknown/path",
			wantService: "",
			wantErr:     true,
		},
		{
			name:        "root path",
			path:        "/",
			wantService: "",
			wantErr:     true,
		},
		{
			name:        "empty path",
			path:        "",
			wantService: "",
			wantErr:     true,
		},
	}

	router := NewRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := router.Route(tt.path)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Route() expected error, got nil")
				}
				return
			}

			// Check no error when not expected
			if err != nil {
				t.Errorf("Route() unexpected error: %v", err)
				return
			}

			// Check service name
			if service != tt.wantService {
				t.Errorf("Route() = %v, want %v", service, tt.wantService)
			}
		})
	}
}
