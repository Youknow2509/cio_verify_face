"""
Main entry point for both REST and gRPC servers
"""
import asyncio
import logging
import sys

# Check dependencies before importing app modules
try:
    from app.main import app
    from app.grpc_server import serve_grpc
    from app.services.face_service import FaceService
    from app.core.config import settings
    from app.core.logging_config import setup_logging
except ImportError as e:
    print(f"\nimport errror: {e}")
    sys.exit(1)
except Exception as e:
    print(f"\nerror during import: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Setup logging
setup_logging()
logger = logging.getLogger(__name__)


async def start_servers():
    """Start both REST API and gRPC servers"""
    # Initialize face service
    logger.info("Initializing Face Service...")
    face_service = FaceService()
    app.state.face_service = face_service
    logger.info("Face service initialized successfully")
    
    # Start gRPC server in background
    grpc_port = settings.GRPC_PORT if hasattr(settings, 'GRPC_PORT') else 50051
    grpc_task = asyncio.create_task(serve_grpc(face_service, grpc_port))
    
    logger.info(f"gRPC server started on port {grpc_port}")
    logger.info(f"REST API will run on port {settings.SERVICE_PORT}")
    
    # Keep the gRPC server running
    await grpc_task


if __name__ == "__main__":
    import uvicorn
    
    # Start gRPC in background thread
    import threading
    
    def run_grpc():
        asyncio.run(start_servers())
    
    grpc_thread = threading.Thread(target=run_grpc, daemon=True)
    grpc_thread.start()
    
    # Start REST API server
    uvicorn.run(
        "app.main:app",
        host=settings.SERVICE_HOST,
        port=settings.SERVICE_PORT,
        workers=1,
        reload=False
    )
