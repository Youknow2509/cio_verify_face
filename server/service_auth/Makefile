# Makefile for gRPC Auth Service

.PHONY: proto build run test clean grpc-test docker

# Variables
PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT_DIR = .

# Default target
all: proto build

# Generate protobuf files
proto:
	@echo "Generating gRPC code from proto files..."
	@chmod +x proto-gen.sh
	@export PATH=$(PATH):$(shell go env GOPATH)/bin && ./proto-gen.sh

# Install protobuf tools
install-proto-tools:
	@echo "Installing protobuf tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build the application
build:
	@echo "Building application..."
	go build -o bin/auth-service cmd/server/main.go

# Run the application
run:
	@echo "Starting auth service..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Test gRPC endpoints with grpcurl
grpc-test:
	@echo "Testing gRPC endpoints..."
	@echo "Listing services:"
	grpcurl -plaintext localhost:50051 list
	@echo "\nTesting ValidateToken:"
	grpcurl -plaintext -d '{"token": "test-token"}' localhost:50051 auth.AuthService/ValidateToken

# Run example gRPC client
example-client:
	@echo "Running example gRPC client..."
	go run examples/grpc_client_example.go

# Clean generated files
clean:
	@echo "Cleaning up..."
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f bin/auth-service

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Lint Go code
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t auth-service:latest .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 50051:50051 auth-service:latest

# Development setup
setup: install-proto-tools proto tidy
	@echo "Development setup complete!"

# Watch for changes and rebuild (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Benchmark gRPC
benchmark:
	@echo "Running gRPC benchmarks..."
	go test -bench=. ./internal/interfaces/grpc/...

# Generate swagger docs
swagger:
	@echo "Generating swagger documentation..."
	./swagger-gen.sh

# Generate SQL code
sqlc:
	@echo "Generating SQL code..."
	./sqlc.sh

# Full setup and run
start: setup build run

# Help
help:
	@echo "Available commands:"
	@echo "  proto           - Generate gRPC code from proto files"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  grpc-test       - Test gRPC endpoints with grpcurl"
	@echo "  example-client  - Run example gRPC client"
	@echo "  clean           - Clean generated files"
	@echo "  fmt             - Format Go code"
	@echo "  lint            - Lint Go code"
	@echo "  tidy            - Tidy dependencies"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  setup           - Install tools and generate proto files"
	@echo "  dev             - Start development server with hot reload"
	@echo "  benchmark       - Run gRPC benchmarks"
	@echo "  swagger         - Generate swagger documentation"
	@echo "  sqlc            - Generate SQL code"
	@echo "  start           - Full setup and run"
	@echo "  help            - Show this help message"