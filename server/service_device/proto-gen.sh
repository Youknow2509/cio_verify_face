#!/bin/bash

# Generate gRPC code from proto files

# Ensure protoc plugins from GOPATH are available
export PATH="$(go env GOPATH)/bin:$PATH"

# Create proto/pb directory if not exists
mkdir -p proto/pb

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/auth.proto \
    proto/face_service.proto

echo "gRPC code generation completed!"