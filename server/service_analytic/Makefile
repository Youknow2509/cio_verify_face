.PHONY: help build run test clean docker-build docker-run deps

# Variables
BINARY_NAME=service_analytic
BUILD_DIR=./bin
CMD_DIR=./cmd/server
GO=go
DOCKER_IMAGE=cio_verify_face/service_analytic
DOCKER_TAG=latest

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

build: deps ## Build the application
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

run: ## Run the application
	$(GO) run $(CMD_DIR)/main.go

test: ## Run tests
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Run linter
	golangci-lint run

swagger: ## Generate Swagger documentation
	swag init -g cmd/server/main.go -o ./internal/docs

proto-gen: ## Generate protobuf files
	./proto-gen.sh

sqlc-gen: ## Generate SQLC files
	export PATH=$PATH:$(go env GOPATH)/bin
	sqlc -f ./config/sqlc.yaml generate

.DEFAULT_GOAL := help
