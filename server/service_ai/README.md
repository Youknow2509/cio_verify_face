# Face Verification AI Service

This service provides face verification functionality for the attendance system using state-of-the-art face recognition models.

## 🚀 Quick Start

**New to this service? Start here:** [QUICKSTART.md](./QUICKSTART.md)

The quickstart guide will walk you through:
1. Installing dependencies
2. Verifying installation
3. Configuring the service
4. Starting the server
5. Troubleshooting common issues

## Features

- **Face Detection**: RetinaFace for accurate face detection and landmark extraction
- **Face Embedding**: ArcFace (InsightFace) for generating high-quality face embeddings
- **Liveness Detection**: Basic anti-spoofing measures
- **Vector Search**: PostgreSQL pgvector for distributed similarity search 🆕
- **Multi-profile Support**: Multiple face profiles per user
- **Soft Delete**: Retention policy for deleted profiles
- **CPU/GPU Support**: Flexible deployment options
- **ScyllaDB Integration**: High-performance authentication state tracking
- **MinIO Integration**: Object storage for face images
- **🚀 Bandwidth Optimization**: Binary upload endpoints (33% less bandwidth than base64)
- **💾 Storage Optimization**: Smart image compression and configurable retention policies
- **🔄 Microservice Ready**: Vector index synchronized across all instances 🆕

## Architecture

```
Camera → Preprocessing → Detection (RetinaFace) → Alignment → 
Embedding (ArcFace) → Compare (pgvector) → Decision → Audit Log
                                              ↓
                                        ScyllaDB (States)
                                        MinIO (Images - Optimized)
```

**Storage Architecture:**
- **PostgreSQL**: Face embeddings and profiles (core data)
- **ScyllaDB**: Authentication states, verification history, and **audit logs** (time-series data)
- **MinIO**: Face images (enrollment & successful verifications)
- **FAISS**: In-memory vector index for fast similarity search

📚 **Documentation:**
- [ScyllaDB & MinIO Integration Guide](./docs/SCYLLADB_MINIO_INTEGRATION.md)
- [Bandwidth & Storage Optimization Guide](./docs/BANDWIDTH_STORAGE_OPTIMIZATION.md) ⚡ **NEW**
- [PgVector Migration Guide](./docs/PGVECTOR_MIGRATION.md) 🆕 **IMPORTANT**

## API Endpoints

### Standard Endpoints (JSON + Base64)
- `POST /api/v1/face/enroll` - Enroll a new face profile
- `POST /api/v1/face/verify` - Verify a face (1:1 or 1:N)

### Optimized Endpoints (Multipart/Form-Data) 🚀 **NEW**
- `POST /api/v1/face/enroll/upload` - Enroll with binary upload (33% less bandwidth)
- `POST /api/v1/face/verify/upload` - Verify with binary upload (33% less bandwidth)

### Management Endpoints
- `PUT /api/v1/face/profile/{profile_id}` - Update face profile
- `DELETE /api/v1/face/profile/{profile_id}` - Delete face profile (soft)
- `POST /api/v1/face/reindex` - Rebuild FAISS index (admin)
- `GET /api/v1/face/profiles/{user_id}` - Get user's face profiles
- `GET /health` - Health check

### Audit & Monitoring Endpoints 🆕 **NEW**
- `GET /api/v1/face/audit/logs` - Get audit logs from ScyllaDB (admin)
- `GET /api/v1/face/audit/user/{user_id}` - Get user-specific audit logs (admin)

## Setup

### Installation Methods

#### Option 1: Using pyproject.toml (Recommended)

1. Copy environment file:
```bash
cp .env.example .env
```

2. Install the package in development mode:
```bash
# CPU mode (default for development)
pip install -e .

# GPU mode (production with CUDA)
pip install -e ".[gpu]"

# Development mode with all tools
pip install -e ".[dev]"

# All dependencies
pip install -e ".[all]"
```

3. Run the service:
```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8080
```

#### Option 2: Using requirements.txt

1. Copy environment file:
```bash
cp .env.example .env
```

2. Install dependencies (CPU mode for development):
```bash
pip install -r requirements-cpu.txt
# or
make install-cpu
```

3. Download models (will be done automatically on first run)

4. Run the service:
```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8080
```

### Production (GPU Mode)

For production deployment with GPU acceleration:

1. Install GPU dependencies (requires CUDA):
```bash
# Using pyproject.toml
pip install -e ".[gpu]"

# Or using requirements.txt
pip install -r requirements-gpu.txt
# or
make install-gpu
```

2. Set compute mode in `.env`:
```bash
COMPUTE_MODE=gpu
```

3. Run the service

**Note**: See [CPU_GPU_GUIDE.md](CPU_GPU_GUIDE.md) for detailed CPU/GPU setup and configuration.

### Production (Docker)

```bash
docker-compose up -d
```

## Configuration

See `.env.example` for all configuration options.

Key settings:
- `DUPLICATE_THRESHOLD`: Threshold for duplicate detection (default: 0.85)
- `VERIFY_THRESHOLD`: Threshold for verification (default: 0.50)
- `MIN_FACE_SIZE`: Minimum face size in pixels (default: 80)
- `SOFT_DELETE_RETENTION`: Days to keep soft-deleted profiles (default: 30)

## Database Schema

The service uses the following tables:
- `face_profiles`: Stores face embeddings and metadata
- `face_audit_logs`: Audit trail for all operations

## Models

The service uses pre-trained models from InsightFace:
- **Detector**: RetinaFace (R50)
- **Embedding**: ArcFace (R100)
- **Anti-spoofing**: Optional liveness detection model

Models are automatically downloaded on first run.

## Monitoring

Metrics are exposed on `/metrics` endpoint in Prometheus format:
- Request latency
- Throughput
- False Accept Rate (FAR)
- False Reject Rate (FRR)
- Index size and lag

## Security

- Embeddings are encrypted at rest
- All API calls require authentication (JWT)
- Audit logging for all operations
- Compliance with data protection regulations

## Performance

- Enrollment: < 1 second per face
- Verification (1:N): < 2 seconds for 10k faces
- Supports GPU acceleration with ONNX Runtime

## Operational Playbook

### Enrollment Failure
- Check face quality (blur, lighting)
- Ensure face size meets minimum requirements
- Check for occlusion (mask, glasses)

### Duplicate Detection
- Review matched profiles
- Admin can merge or assign to existing user
- Set duplicate threshold based on use case

### Re-indexing
- Schedule during low-traffic periods
- Monitor index lag metric
- Use canary deployment for model upgrades

### Backup
- Daily backup of embeddings and index
- Keep versioned snapshots
- Test restore procedures regularly
