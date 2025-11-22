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
    from app.main import app
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


# ---- Flush callbacks ----
def before_flush(count: int):
    logger.info("[Flush] Chuẩn bị gửi %d bản ghi ...", count)


def after_flush(count: int, status_code: Optional[int], message: Optional[str]):
    logger.info("[Flush] Hoàn tất %d status=%s message=%s", count, status_code, message)


# ---- Build service session ----
def build_service_session() -> attendance_pb2.ServiceSessionInfo:
    return attendance_pb2.ServiceSessionInfo(
        service_id=(getattr(settings, "SERVICE_ID", None) or "face_service"),
        service_name=(getattr(settings, "SERVICE_NAME", None) or "FaceService"),
        client_ip="",
        client_agent="FaceServiceAgent",
    )


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

    # Gán vào app.state để tái sử dụng
    app.state.attendance_client = attendance_client
    
    # Initialize batching service và face service ngay sau khi check infrastructure
    service_session = build_service_session()
    
    batch_size = settings.ATTENDANCE_BATCH_MAX_SIZE if getattr(settings, "ATTENDANCE_BATCH_MAX_SIZE", None) is not None else 50
    flush_interval = settings.ATTENDANCE_BATCH_FLUSH_INTERVAL if getattr(settings, "ATTENDANCE_BATCH_FLUSH_INTERVAL", None) is not None else 3.0
    max_pending = settings.ATTENDANCE_BATCH_MAX_PENDING if getattr(settings, "ATTENDANCE_BATCH_MAX_PENDING", None) is not None else 100

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
    app.state.attendance_batching_service = attendance_batching_service
    
    logger.info("Initializing Face Service with batching support...")
    face_service = FaceService(batching_service=attendance_batching_service)
    app.state.face_service = face_service
    logger.info("Face Service ready.")
    
    logger.info("All infrastructure healthy.")


# ---- Async gRPC startup ----
async def start_grpc_server():
    """
    Chạy gRPC server với FaceService đã khởi tạo.
    """
    face_service = app.state.face_service
    grpc_port = getattr(settings, "GRPC_PORT", 50051)
    logger.info("Starting gRPC server on port %d", grpc_port)
    
    asyncio.create_task(batching_heartbeat())
    await serve_grpc(face_service, grpc_port)

    logger.info("gRPC server terminated - shutting down batching service...")
    app.state.attendance_batching_service.close(flush_final=True)
    logger.info("Batching service closed.")

async def batching_heartbeat():
    """
    Task nền: mỗi N giây log thống kê của batching service (tuỳ chọn).
    """
    while True:
        await asyncio.sleep(15)
        svc = getattr(app.state, "attendance_batching_service", None)
        if svc:
            stats = svc.stats()
            logger.debug("[BatchHeartbeat] %s", stats)


# ---- Graceful shutdown handler (for SIGTERM/SIGINT) ----
def install_signal_handlers():
    def handle_sig(sig, frame):
        logger.warning("Received signal %s - initiating shutdown...", sig)
        # Flush + close batching service
        batching = getattr(app.state, "attendance_batching_service", None)
        if batching:
            try:
                batching.close(flush_final=True)
            except Exception:
                logger.exception("Error closing batching service")
        # Close attendance client
        attendance_client = getattr(app.state, "attendance_client", None)
        if attendance_client:
            try:
                attendance_client.close()
            except Exception:
                logger.exception("Error closing attendance client")
        # FaceService cleanup nếu cần
        face_service = getattr(app.state, "face_service", None)
        if hasattr(face_service, "close"):
            try:
                face_service.close()
            except Exception:
                logger.exception("Error closing face service")
        # Thoát tiến trình
        sys.exit(0)

    signal.signal(signal.SIGINT, handle_sig)
    signal.signal(signal.SIGTERM, handle_sig)


# ---- Main Entrypoint ----
if __name__ == "__main__":
    import uvicorn

    install_signal_handlers()

    # 1. Kiểm tra hạ tầng
    check_infrastructure()

    # 2. Khởi chạy gRPC + batching trong thread riêng
    def run_grpc_thread():
        asyncio.run(start_grpc_server())

    grpc_thread = threading.Thread(target=run_grpc_thread, name="GRPCServerThread", daemon=True)
    grpc_thread.start()

    # 3. Khởi chạy REST API (blocking ở main thread)
    # workers = settings.SERVICE_WORKERS if getattr(settings, "SERVICE_WORKERS", None) is not None else 1
    logger.info("Starting REST API server on %s:%s (workers=%d reload=%s)",
                settings.SERVICE_HOST, settings.SERVICE_PORT, 1, False)

    uvicorn.run(
        "app.main:app",
        host=settings.SERVICE_HOST,
        port=settings.SERVICE_PORT,
        workers=1,
        reload=False,
    )