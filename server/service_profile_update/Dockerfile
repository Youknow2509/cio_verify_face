# docker buildx build --platform linux/amd64,linux/arm64 -t cio_verify_face_service_profile_update .
# docker buildx build -t cio_verify_face_service_profile_update .
# Build stage
FROM golang:1.24-trixie AS builder

# Set working directory
WORKDIR /build

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o service_profile_update ./cmd/main.go 

# Production stage
FROM debian:latest

WORKDIR /prod

# Copy the built binary from the builder stage
COPY --from=builder /build/service_profile_update .

# Expose the application port
EXPOSE 8080
EXPOSE 50051

# Command to run the application
CMD ["./service_profile_update"]