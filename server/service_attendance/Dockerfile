# docker buildx build --platform linux/amd64,linux/arm64 -t someone2509/cio_verify_face_service_attendance:latest .
# docker buildx build -t cio_verify_face_service_attendance .
# Build stage
FROM golang:1.24-trixie AS builder

# Set working directory
WORKDIR /build

# Copy shared packages first
COPY pkg/ ./pkg/

# Copy service source code
COPY service_attendance/ ./service_attendance/

# Set working directory to service
WORKDIR /build/service_attendance

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o service_attendance ./cmd/server/main.go 

# Production stage
FROM debian:latest

WORKDIR /prod

# Copy the built binary from the builder stage
COPY --from=builder /build/service_attendance/service_attendance .

# Expose the application port
EXPOSE 8080
EXPOSE 50051

# Command to run the application
CMD ["./service_attendance"]