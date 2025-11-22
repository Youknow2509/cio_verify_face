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
    from app.grpc.server.grpc_server import serve_grpc
    from app.services.face_service import FaceService
    from app.core.config import settings
    from app.core.logging_config import setup_logging
    from app.grpc.client.attendance_client import AttendanceClient 
    from app.services.attendance_batching_service import AttendanceBatchingService
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
    # Postgres
    from app.database.pgvector_manager import PgVectorManager
    try:
        PgVectorManager().check_connection()
        logger.info("Postgres connection OK")
    except Exception as e:
        logger.error("Postgres connection failed: %s", e)
        sys.exit(1)

    # ScyllaDB
    from app.database.scylladb_manager import ScyllaDBManager
    try:
        ScyllaDBManager().check_connection()
        logger.info("ScyllaDB connection OK")
    except Exception as e:
        logger.error("ScyllaDB connection failed: %s", e)
        sys.exit(1)

    # Auth service
    from app.grpc.client.auth_client import AuthClient
    auth_client = AuthClient()
    try:
        auth_client.check_connection()
        logger.info("Auth service connection OK")
    except Exception as e:
        logger.error("Auth service connection failed: %s", e)
        sys.exit(1)

    # Attendance service
    attendance_client = AttendanceClient()
    try:
        attendance_client.check_connection()
        logger.info("Attendance service connection OK")
    except Exception as e:
        logger.error("Attendance service connection failed: %s", e)
        sys.exit(1)
    
    logger.info("Initializing Face Service with batching support...")
    logger.info("All infrastructure healthy.")


# ---- Async gRPC startup ----
async def start_grpc_server():
    """
    Chạy gRPC server với FaceService đã khởi tạo.
    """
    grpc_port = getattr(settings, "GRPC_PORT", 50051)
    logger.info("Starting gRPC server on port %d", grpc_port)


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