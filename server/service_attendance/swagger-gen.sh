#!/bin/bash
@echo "Generating Swagger documentation..."
export PATH=$PATH:$(go env GOPATH)/bin   
swag init -g ./cmd/server/main.go -o ./cmd/swag/docs
echo "Swagger documentation generation completed."