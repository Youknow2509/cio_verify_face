"""
Configuration settings for the service
"""
from pydantic_settings import BaseSettings
from typing import List, Optional
import os


class Settings(BaseSettings):
    """Application settings"""
    
    # Environment
    ENVIRONMENT: str = "development"
    SERVICE_NAME: str = "service_ai"
    SERVICE_ID: str = "service_ai_001"
    SERVICE_WORKERS: int = 4
    
    # Service
    SERVICE_NAME: str = "service_ai"
    SERVICE_HOST: str = "0.0.0.0"
    SERVICE_PORT: int = 8080
    GRPC_PORT: int = 50051
    
    # Compute mode (cpu or gpu)
    COMPUTE_MODE: str = "cpu"
    
    # Observability Configuration
    OBSERVABILITY_ENABLED: bool = True
    METRICS_PORT: int = 9090
    METRICS_PATH: str = "/metrics"
    TRACING_ENABLED: bool = True
    OTLP_ENDPOINT: str = "http://jaeger:4318/v1/traces"
    
    # GRPC client configuration
    GRPC_CLIENT_KEEPALIVE_TIME_MS: int = 120000
    GRPC_CLIENT_KEEPALIVE_TIMEOUT_MS: int = 20000
    GRPC_CLIENT_KEEPALIVE_PERMIT_WITHOUT_CALLS: int = 0
    GRPC_CLIENT_HTTP2_MAX_PINGS_WITHOUT_DATA: int = 1
    GRPC_CLIENT_HTTP2_MIN_TIME_BETWEEN_PINGS_MS: int = 60000
    GRPC_CLIENT_HTTP2_MIN_PING_INTERVAL_WITHOUT_DATA_MS: int = 60000

    # GRPC server keepalive settings
    GRPC_SERVER_KEEPALIVE_TIME_MS: int = 120000
    GRPC_SERVER_KEEPALIVE_TIMEOUT_MS: int = 20000
    GRPC_SERVER_HTTP2_MIN_TIME_BETWEEN_PINGS_MS: int = 60000
    GRPC_SERVER_KEEPALIVE_PERMIT_WITHOUT_CALLS: int = 1
    
    # GRPC auth client settings
    GRPC_AUTH_URL: str = "localhost:50051"
    GRPC_AUTH_TLS: bool = False
    GRPC_AUTH_CERT_PATH: Optional[str] = None
    GRPC_AUTH_KEY_PATH: Optional[str] = None
    
    # GRPC attendance client settings
    GRPC_ATTENDANCE_URL: str = "localhost:50052"
    GRPC_ATTENDANCE_TLS: bool = False
    GRPC_ATTENDANCE_CERT_PATH: Optional[str] = None
    GRPC_ATTENDANCE_KEY_PATH: Optional[str] = None
    
    # Attendance batching settings
    ATTENDANCE_BATCH_MAX_SIZE: Optional[int] = 10
    ATTENDANCE_BATCH_FLUSH_INTERVAL: Optional[float] = 3.0
    ATTENDANCE_BATCH_MAX_PENDING: Optional[int] = 100
    
    # Database
    DATABASE_URL: str = "postgresql://postgres:postgres@localhost:5432/cio_attendance_db"
    
    # Redis Configuration
    REDIS_TYPE: int = 1  # 1: standalone, 2: sentinel, 3: cluster
    REDIS_USE_TLS: bool = False
    REDIS_CERT_PATH: str = "./config/redis/cert.pem"
    REDIS_KEY_PATH: str = "./config/redis/key.pem"
    REDIS_PASSWORD: str = "root1234"
    REDIS_DB: int = 0
    
    # Redis Standalone
    REDIS_HOST: str = "127.0.0.1"
    REDIS_PORT: int = 6379
    
    # Redis Sentinel
    REDIS_MASTER_NAME: str = "mymaster"
    REDIS_SENTINEL_ADDRS: List[str] = ["127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"]
    
    # Redis Cluster
    REDIS_CLUSTER_ADDRS: List[str] = ["127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"]
    REDIS_ROUTE_BY_LATENCY: bool = True
    REDIS_ROUTE_RANDOMLY: bool = False
    
    # Redis Pool Configuration
    REDIS_POOL_SIZE: int = 10
    REDIS_MIN_IDLE_CONNS: int = 2
    REDIS_MAX_RETRIES: int = 3
    
    # ScyllaDB - for authentication state tracking
    SCYLLADB_HOSTS: str = "localhost"
    SCYLLADB_PORT: int = 9042
    SCYLLADB_KEYSPACE: str = "cio_verify_face"
    SCYLLADB_USERNAME: str = "cassandra"
    SCYLLADB_PASSWORD: str = "root1234"
    
    # MinIO - for face image storage
    MINIO_ENDPOINT: str = "localhost:9000"
    MINIO_ACCESS_KEY: str = "minioadmin"
    MINIO_SECRET_KEY: str = "minioadmin"
    MINIO_SECURE: bool = False
    MINIO_BUCKET_FACES: str = "face-images"
    MINIO_BUCKET_VERIFICATIONS: str = "verification-images"
    
    # Image optimization settings
    IMAGE_MAX_SIZE: int = 1920  # Max width/height in pixels
    IMAGE_QUALITY: int = 85  # JPEG quality (1-100)
    IMAGE_STORE_ENROLLMENTS: bool = True  # Always store enrollment images
    IMAGE_STORE_VERIFICATIONS: bool = False  # Only store verification images on demand
    IMAGE_STORE_FAILED_VERIFICATIONS: bool = False  # Store failed verification attempts
    
    # Model settings
    FACE_DETECTOR_MODEL: str = "retinaface_r50_v1"
    FACE_EMBEDDING_MODEL: str = "arcface_r100_v1"
    EMBEDDING_DIMENSION: int = 512
    
    # Thresholds
    QUALITY_THRESHOLD: float = 0.5       # Image quality threshold for enrollment
    DUPLICATE_THRESHOLD: float = 0.95     # Threshold to detect face already enrolled to different user (95% = likely same person)
    DUPLICATE_GAP_THRESHOLD: float = 0.08  # Minimum gap between top 2 matches to confirm duplicate (prevent ambiguous cases)
    VERIFY_THRESHOLD: float = 0.80        # Threshold to verify match in 1:1 verification (80% = high confidence match)
    MIN_FACE_SIZE: int = 80
    
    # Vector search settings (pgvector)
    # Note: pgvector index is automatically maintained by PostgreSQL
    # No need for manual rebuild intervals like FAISS
    VECTOR_INDEX_REBUILD_INTERVAL: int = 3600  # seconds (for compatibility)
    VECTOR_DB_INDEX_VERSION: int = 1

    # Milvus settings
    MILVUS_URI: Optional[str] = None  # e.g. https://your-project.api.gcp-milvus.zillizcloud.com
    MILVUS_TOKEN: Optional[str] = None  # Zilliz/Milvus Cloud token
    MILVUS_HOST: str = "localhost"
    MILVUS_PORT: int = 19530
    MILVUS_SECURE: bool = False
    MILVUS_USERNAME: Optional[str] = None
    MILVUS_PASSWORD: Optional[str] = None
    MILVUS_DB: str = "default"
    MILVUS_COLLECTION: str = "face_profiles"
    MILVUS_INDEX_TYPE: str = "IVF_FLAT"  # Options: FLAT, IVF_FLAT, IVF_SQ8, HNSW, AUTOINDEX
    MILVUS_METRIC_TYPE: str = "IP"   # Options: L2, COSINE, IP (IP=Inner Product for normalized vectors, best for face embeddings)
    MILVUS_NLIST: int = 1024
    MILVUS_NPROBE: int = 16
    
    # Liveness detection
    LIVENESS_ENABLED: bool = True
    LIVENESS_THRESHOLD: float = 0.7
    
    # Security
    ENCRYPTION_KEY: Optional[str] = None
    JWT_SECRET: Optional[str] = None
    
    # Logging
    LOG_LEVEL: str = "INFO"
    LOG_FORMAT: str = "json"
    LOG_FILEPATH: str = "./service_ai.log"
    
    # Storage
    STORAGE_PATH: str = "./data"
    MODEL_PATH: str = "~/.insightface"
    INDEX_PATH: str = "./data/indexes"
    
    # Performance
    MAX_WORKERS: int = 4
    BATCH_SIZE: int = 32
    
    # Retention policy (days)
    SOFT_DELETE_RETENTION: int = 30
    
    # FAISS settings
    FAISS_INDEX_TYPE: str = "IVF"  # Options: "Flat",
    FAISS_REBUILD_INTERVAL: int = 3600
    
    # Pydantic v2 configuration: allow extra env vars (ignore unknowns)
    model_config = {
        "env_file": ".env",
        "case_sensitive": True,
        # If extra env vars are present (e.g. AUTH_SERVICE_URL) ignore them
        "extra": "ignore",
    }


# Create settings instance
settings = Settings()

# Create directories if they don't exist (with error handling)
def _ensure_directories():
    """Create necessary directories if they don't exist"""
    try:
        os.makedirs(settings.STORAGE_PATH, exist_ok=True)
        os.makedirs(settings.MODEL_PATH, exist_ok=True)
        os.makedirs(settings.INDEX_PATH, exist_ok=True)
    except OSError as e:
        # If we can't create directories (e.g., read-only filesystem),
        # warn but don't fail - let the application handle it later
        import logging
        logging.warning(f"Could not create directories: {e}")

_ensure_directories()
