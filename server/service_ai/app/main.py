"""
Main FastAPI application
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from prometheus_client import make_asgi_app
import logging

from app.core.config import settings
from app.core.logging_config import setup_logging
from app.api import face_routes
from app.services.face_service import FaceService
from app.services.auth_client import get_client

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
    allow_origins=["*"],  # Configure based on your needs
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Mount prometheus metrics
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

# Include routers
app.include_router(face_routes.router, prefix="/api/v1/face", tags=["face"])


@app.on_event("startup")
async def startup_event():
    """Initialize services on startup"""
    logger.info("Starting Face Verification Service...")
    logger.info(f"Environment: {settings.ENVIRONMENT}")
    logger.info(f"Service: {settings.SERVICE_NAME}")
    logger.info(f"Compute Mode: {settings.COMPUTE_MODE}")
    
    try:
        # Initialize face service (loads models)
        logger.info("Initializing face service (this may take a minute on first run)...")
        face_service = FaceService()
        app.state.face_service = face_service
        # Initialize a shared AuthClient for middleware to reuse
        try:
            app.state.auth_client = get_client()
        except Exception:
            # Don't fail startup if auth client can't be created now; middleware will attempt its own
            logger.warning("Could not initialize shared AuthClient at startup")
        logger.info("Face service initialized successfully")
    except ModuleNotFoundError as e:
        logger.error(f"Missing dependency: {e}")
        logger.error("Please install dependencies:")
        logger.error("  pip install -e .")
        logger.error("  or")
        logger.error("  pip install -r requirements-cpu.txt")
        logger.error("\nFor more help, see QUICKSTART.md")
        raise
    except Exception as e:
        logger.error(f"Failed to initialize face service: {e}")
        logger.error("Check logs above for details. For troubleshooting, see QUICKSTART.md")
        raise


@app.on_event("shutdown")
async def shutdown_event():
    """Cleanup on shutdown"""
    logger.info("Shutting down Face Verification Service...")
    if hasattr(app.state, 'face_service'):
        # Cleanup if needed
        pass
    # Close auth client if we created one on startup
    if hasattr(app.state, 'auth_client'):
        try:
            app.state.auth_client.close()
        except Exception:
            logger.exception("Error closing auth client on shutdown")


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
