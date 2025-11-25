.PHONY: all build test clean lint run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=service_profile_update

help:
	@echo "Available targets:"
		@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

all: help

swagger-gen: # Generate Swagger documentation
	swagger generate spec -o ./docs/swagger.yaml --scan-models

build: # Build the Go app
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

test: # Run tests
	$(GOTEST) -v ./tests/...

test-coverage: # Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean: # Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

lint: # Run linter
	golangci-lint run ./...
	golangci-lint run ./...

deps: # Download dependencies
	$(GOMOD) download

tidy: # Tidy up dependencies 
	$(GOMOD) tidy

run: # Run the application
	$(GOCMD) run ./cmd/main.go

docker-build: # Build Docker image
	docker build -t $(BINARY_NAME):latest .

.DEFAULT_GOAL := all
