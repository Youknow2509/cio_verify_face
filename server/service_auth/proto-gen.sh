#!/bin/bash

# Generate gRPC code from proto files

# Create proto/pb directory if not exists
mkdir -p proto/pb

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/auth.proto

echo "gRPC code generation completed!"