# docker buildx build --platform linux/amd64,linux/arm64 -t cio_verify_face_service_workforce .
# docker buildx build -t cio_verify_face_service_workforce .
# Build stage
FROM golang:1.24-trixie AS builder

# Set working directory
WORKDIR /build

# Copy shared packages first
COPY pkg/ ./pkg/

# Copy service source code
COPY service_workforce/ ./service_workforce/

# Set working directory to service
WORKDIR /build/service_workforce

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o service_workforce ./cmd/server/main.go 

# Production stage
FROM debian:latest

WORKDIR /prod

# Copy the built binary from the builder stage
COPY --from=builder /build/service_workforce/service_workforce .

# Expose the application port
EXPOSE 8080
EXPOSE 50051

# Command to run the application
CMD ["./service_workforce"]