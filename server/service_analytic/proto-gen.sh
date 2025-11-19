#!/bin/bash

# Generate gRPC code from proto files

PROTO_DIR="proto"
OUT_DIR="proto/pb"

# Create output directory if it doesn't exist
mkdir -p $OUT_DIR

# Generate proto files
protoc --go_out=$OUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
    $PROTO_DIR/analytic.proto

protoc --go_out=$OUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
    $PROTO_DIR/auth.proto

echo "Proto generation complete!"
