#!/bin/bash

echo "Generating gRPC code..."
protoc --go_out=. --go_opt=paths=source_relative \
                --go-grpc_out=. --go-grpc_opt=paths=source_relative \
                proto/ws.proto

echo "gRPC code generation completed."