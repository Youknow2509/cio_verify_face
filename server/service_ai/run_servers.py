"""
Main entry point for both REST (FastAPI + Uvicorn) and gRPC servers,
kèm AttendanceBatchingService (gom và flush batch attendance record).
"""
import asyncio
import logging
import sys
import threading
import signal
from typing import Optional

from app.grpc_generated import attendance_pb2

# ---- Import & dependency checks ----
try:
    from app.grpc.server.grpc_server import serve
    from app.services.face_service import FaceService
    from app.services.user_service import UserService
    from app.core.config import settings
    from app.core.logging_config import setup_logging
    from app.grpc.client.attendance_client import AttendanceClient 
    from app.services.attendance_batching_service import AttendanceBatchingService
    from app.main import build_service_session, before_flush, after_flush
except ImportError as e:
    print(f"\nImport error: {e}")
    sys.exit(1)
except Exception as e:
    print(f"\nError during import: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

setup_logging()
logger = logging.getLogger(__name__)

# ---- Infrastructure checks ----
def check_infrastructure():
    logger.info("Checking infrastructure connections...")
    # Init Redis Cache
    from app.database.redis_manager import RedisManager
    redis_cache = RedisManager()
    if redis_cache.check_connection() is False:
        logger.warning("Redis cache is disabled due to connection failure")
        exit(1)
    # Init PGManager
    from app.database.pg_manager import PGManager
    pg_manager = PGManager()
    if pg_manager.check_connection() is False:
        logger.error("PostgreSQL database connection failed during startup")
        exit(1)
    # Init MilvusManager
    from app.database.milvus_manager import MilvusManager
    milvus_manager = MilvusManager()
    if milvus_manager.check_connection() is False:
        logger.error("Milvus database connection failed during startup")
        exit(1)
    # Init ScyllaManager
    from app.database.scylladb_manager import ScyllaDBManager
    scylla_manager = ScyllaDBManager()
    if scylla_manager.check_connection() is False:
        logger.error("ScyllaDB connection failed during startup")
        exit(1)
    # Init attendance_client grpc
    from app.grpc.client.attendance_client import AttendanceClient
    attendance_client = AttendanceClient()
    if attendance_client.check_connection() is False:
        logger.error("Attendance gRPC client connection failed during startup")
        exit(1)
    # Init auth_client grpc
    from app.grpc.client.auth_client import AuthClient
    auth_client = AuthClient()
    if auth_client.check_connection() is False:
        logger.error("Auth gRPC client connection failed during startup")
        exit(1)
    
    logger.info("Initializing Face Service with batching support...")
    logger.info("All infrastructure healthy.")


# ---- Async gRPC startup ----
async def start_grpc_server():
    """
    Chạy gRPC server với FaceService và UserService đã khởi tạo.
    """
    logger.info("Initializing services for gRPC server...")

    # This logic duplicates logic from app/main.py's startup event.
    # A refactor might be needed to share service instances.
    from app.database.redis_manager import RedisManager
    from app.database.pg_manager import PGManager
    from app.database.milvus_manager import PgVectorManager
    
    redis_cache = RedisManager()
    pg_manager = PGManager()
    pgvector_manager = PgVectorManager()

    # Ensure Milvus connectivity before proceeding
    if pgvector_manager.check_connection() is False:
        logger.error("Milvus/PgVector database connection failed during gRPC startup")
        exit(1)
    
    attendance_client = AttendanceClient()
    service_session = build_service_session()
    batch_size = getattr(settings, "ATTENDANCE_BATCH_MAX_SIZE", 50)
    flush_interval = getattr(settings, "ATTENDANCE_BATCH_FLUSH_INTERVAL", 3.0)
    max_pending = getattr(settings, "ATTENDANCE_BATCH_MAX_PENDING", 100)
    
    attendance_batching_service = AttendanceBatchingService(
        client=attendance_client,
        service_session=service_session,
        max_batch_size=batch_size,
        flush_interval=flush_interval,
        max_pending_records=max_pending,
        on_before_flush=before_flush,
        on_after_flush=after_flush,
        metadata=[("x-origin", "hybrid")],
    )
    
    face_service = FaceService(batching_service=attendance_batching_service)
    # Reuse the verified Milvus manager instance
    face_service.index_manager = pgvector_manager
    user_service = UserService(redis_client=redis_cache, postgres_client=pg_manager)

    await serve(face_service, user_service)


# ---- Main Entrypoint ----
if __name__ == "__main__":
    import uvicorn

    # 1. Kiểm tra hạ tầng
    check_infrastructure()

    # 2. Khởi chạy gRPC + batching trong thread riêng
    def run_grpc_thread():
        asyncio.run(start_grpc_server())

    grpc_thread = threading.Thread(target=run_grpc_thread, name="GRPCServerThread", daemon=True)
    grpc_thread.start()

    # 3. Khởi chạy REST API (blocking ở main thread)
    workers = settings.SERVICE_WORKERS if getattr(settings, "SERVICE_WORKERS", None) is not None else 1
    logger.info("Starting REST API server on %s:%s (workers=%d reload=%s)",
                settings.SERVICE_HOST, settings.SERVICE_PORT, workers, False)

    uvicorn.run(
        "app.main:app",
        host=settings.SERVICE_HOST,
        port=settings.SERVICE_PORT,
        workers=workers,
        reload=False,
    )