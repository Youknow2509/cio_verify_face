## Multi-arch CPU image for Face Verification Service
## Usage:
##   docker build -t face-service:cpu .
##   docker run --rm -p 8080:8080 -p 50051:50051 \
##     -e DATABASE_URL=postgresql://... \
##     -e COMPUTE_MODE=cpu \
##     face-service:cpu

# Builder stage: compile native deps (insightface) with toolchain
FROM python:3.12-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    g++ \
    libgl1 \
    libglib2.0-0 \
    libgomp1 \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ARG REQUIREMENTS_FILE=requirements-cpu.txt
COPY requirements-base.txt ${REQUIREMENTS_FILE} ./
RUN python -m pip install --upgrade pip \
    && pip install --no-cache-dir -r ${REQUIREMENTS_FILE}

# Runtime stage: slim image without compiler toolchain
FROM python:3.12-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# System runtime libs for OpenCV, ONNXRuntime, etc.
RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    libgomp1 \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Bring installed site-packages and binaries from builder
COPY --from=builder /usr/local /usr/local

# Copy application source
COPY . .

# Regenerate gRPC stubs to match installed protobuf runtime
RUN bash generate_grpc.sh

# Create non-root user and set permissions for writable dirs
RUN useradd -m appuser \
    && mkdir -p /app/data /app/models /app/data/indexes \
    && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080 50051

# Default environment (override at runtime)
ENV SERVICE_HOST=0.0.0.0 \
    SERVICE_PORT=8080 \
    GRPC_PORT=50051 \
    COMPUTE_MODE=cpu \
    STORAGE_PATH=/app/data \
    MODEL_PATH=/app/models \
    INDEX_PATH=/app/data/indexes

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=5 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["python", "-u", "run_servers.py"]

