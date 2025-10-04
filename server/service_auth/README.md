# Contact:

-   **Mail**: *lytranvinh.work@gmail.com*
-   **Github**: *https://github.com/Youknow2509*

# Auth Service

Microservice xác thực cho hệ thống CIO Verify Face, hỗ trợ cả HTTP REST API và gRPC cho inter-service communication.

## Features

### HTTP REST API

-   User authentication và authorization
-   JWT token management
-   Session management
-   Device management
-   Company permission checking

### gRPC Service

-   Inter-service authentication
-   Token validation cho các service khác
-   User information retrieval
-   Permission checking
-   Device validation
-   Batch operations

## Architecture

-   **Clean Architecture** với layers: Domain, Application, Infrastructure, Interfaces
-   **Multi-level Caching**: Local (Ristretto) + Distributed (Redis)
-   **Dual Protocol Support**: HTTP REST + gRPC
-   **Database**: PostgreSQL với SQLC
-   **Message Queue**: Kafka
-   **Object Storage**: MinIO

## Quick Start

### 1. Using Make

```bash
# Setup và chạy
make start

# Chỉ generate gRPC code
make proto

# Test gRPC endpoints
make grpc-test
```

### 2. Using Docker

```bash
# Chạy toàn bộ stack
docker-compose up -d

# Chỉ chạy auth service
docker-compose up auth-service
```

### 3. Manual Setup

```bash
# Generate gRPC code
./proto-gen.sh

# Run service
go run cmd/server/main.go
```

## Services

### HTTP API

-   **Port**: 8080
-   **Documentation**: `/swagger/index.html`
-   **Health Check**: `/health`

### gRPC Service

-   **Port**: 50051
-   **Documentation**: [GRPC_README.md](./GRPC_README.md)
-   **Reflection**: Enabled

## gRPC Methods

1. **ValidateToken** - Xác thực access token
2. **GetUserInfo** - Lấy thông tin user
3. **CheckUserPermission** - Kiểm tra quyền user
4. **CheckDeviceInCompany** - Kiểm tra device thuộc company
5. **ValidateDeviceSession** - Xác thực device session
6. **BatchValidateTokens** - Xác thực nhiều token cùng lúc

## Configuration

### Config Files

-   `config/config.yaml` - Production config
-   `config/config.dev.yaml` - Development config

### Environment Variables

```bash
CONFIG_PATH=/path/to/config.yaml
```

### gRPC Configuration

```yaml
grpc:
    network: 'tcp'
    port: 50051
    tls:
        enabled: false
        cert_file: ''
        key_file: ''
```
