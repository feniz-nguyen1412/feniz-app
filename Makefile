.PHONY: build run clean proto docker-up docker-down

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Clean build artifacts
clean:
	rm -rf bin/

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# View Docker logs
logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Development setup
dev-setup: deps proto docker-up
	@echo "Development environment is ready!"
	@echo "HTTP API: http://localhost:8080"
	@echo "gRPC API: localhost:50051"
	@echo "Kafka UI: http://localhost:8080 (kafka-ui container)"

# Full development workflow
dev: dev-setup build run
