"""
Main FastAPI application
"""
from typing import Optional
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from prometheus_client import make_asgi_app
import logging

from app.core.config import settings
from app.core.logging_config import setup_logging
from app.api import face_routes
from app.grpc_generated import attendance_pb2
# Setup logging
setup_logging()
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Face Verification Service",
    description="AI service for face verification in attendance system",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Mount prometheus metrics
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

# Include routers
app.include_router(face_routes.router, prefix="/api/v1/face", tags=["face"])

@app.get("/health", tags=["health"])
async def health_check():
    """Health check endpoint"""
    return JSONResponse(
        status_code=200,
        content={
            "status": "healthy",
            "service": settings.SERVICE_NAME,
            "version": "1.0.0",
            "environment": settings.ENVIRONMENT
        }
    )


@app.get("/", tags=["root"])
async def root():
    """Root endpoint"""
    return {
        "service": "Face Verification Service",
        "version": "1.0.0",
        "status": "running",
        "docs": "/docs"
    }

@app.on_event("startup")
async def startup_event():
    logger.info("Starting up the Face Verification Service...")
    # Init attendance_client grpc
    from app.grpc.client.attendance_client import AttendanceClient
    attendance_client = AttendanceClient()
    app.state.attendance_client = attendance_client
    # Init auth_client grpc
    from app.grpc.client.auth_client import AuthClient
    auth_client = AuthClient()
    app.state.auth_client = auth_client
    # Initialize batching service
    service_session = build_service_session()
    batch_size = getattr(settings, "ATTENDANCE_BATCH_MAX_SIZE", 50)
    flush_interval = getattr(settings, "ATTENDANCE_BATCH_FLUSH_INTERVAL", 3.0)
    max_pending = getattr(settings, "ATTENDANCE_BATCH_MAX_PENDING", 100)
    from app.services.attendance_batching_service import AttendanceBatchingService
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
    # Init FaceService
    from app.services.face_service import FaceService
    face_service = FaceService(batching_service=attendance_batching_service)
    app.state.face_service = face_service
    
@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down the Face Verification Service...")
    batching = getattr(app.state, "attendance_batching_service", None)
    attendance_client = getattr(app.state, "attendance_client", None)
    face_service = getattr(app.state, "face_service", None)
    if batching:
        try:
            batching.close(flush_final=True)
        except Exception:
            logger.exception("Error closing batching service")
    if attendance_client:
        try:
            attendance_client.close()
        except Exception:
            logger.exception("Error closing attendance client")
    if hasattr(face_service, "close"):
        try:
            face_service.close()
        except Exception:
            logger.exception("Error closing face service")
            
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