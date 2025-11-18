FROM python:3.10-slim

# Build argument for compute mode (cpu or gpu)
ARG COMPUTE_MODE=cpu

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements files
COPY requirements-base.txt requirements-cpu.txt requirements-gpu.txt ./

# Install Python dependencies based on compute mode
RUN if [ "$COMPUTE_MODE" = "gpu" ]; then \
        pip install --no-cache-dir -r requirements-gpu.txt; \
    else \
        pip install --no-cache-dir -r requirements-cpu.txt; \
    fi

# Copy proto files and generate gRPC code
COPY protos ./protos
COPY generate_grpc.sh .
RUN chmod +x generate_grpc.sh && ./generate_grpc.sh

# Copy application code
COPY . .

# Create directory for models and data
RUN mkdir -p /app/models /app/data/indexes /app/logs

# Expose port
EXPOSE 8080 50051

# Set default compute mode environment variable
ENV COMPUTE_MODE=${COMPUTE_MODE}

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:8080/health')" || exit 1

# Run the application
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8080", "--workers", "1"]
