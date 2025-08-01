.PHONY: gateway-build gateway-dev gateway-docker gateway-test

# Local development
gateway-dev:
	cd services/gateway && go run cmd/api/main.go

# Docker commands  
gateway-docker:
	docker-compose up --build gateway

gateway-docker-logs:
	docker-compose logs -f gateway

# Testing
gateway-test:
	cd services/gateway && go test ./...

# Cleanup
clean:
	docker-compose down
	docker system prune -f

# Swagger commands
gateway-swagger-gen:
	cd services/gateway && swag init -g cmd/api/main.go -o docs

gateway-swagger-dev: gateway-swagger-gen
	cd services/gateway && PORT=8080 LOG_LEVEL=debug go run cmd/api/main.go

gateway-swagger-open:
	start http://localhost:8080/swagger/index.html