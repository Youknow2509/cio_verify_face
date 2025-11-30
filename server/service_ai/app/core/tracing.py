"""
OpenTelemetry tracing configuration for Jaeger
"""
import logging
from typing import Optional

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import Resource, SERVICE_NAME
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from app.core.config import settings

logger = logging.getLogger(__name__)

_tracer_provider: Optional[TracerProvider] = None


def init_tracing(app=None) -> Optional[TracerProvider]:
    """
    Initialize OpenTelemetry tracing with OTLP exporter for Jaeger
    
    Args:
        app: FastAPI application instance to instrument
        
    Returns:
        TracerProvider instance or None if tracing is disabled
    """
    global _tracer_provider
    
    if not settings.TRACING_ENABLED:
        logger.info("Tracing is disabled")
        return None
    
    if not settings.OTLP_ENDPOINT:
        logger.warning("OTLP_ENDPOINT not configured, tracing disabled")
        return None
    
    try:
        # Create resource with service information
        resource = Resource.create({
            SERVICE_NAME: settings.SERVICE_NAME,
            "service.version": "1.0.0",
            "environment": settings.ENVIRONMENT,
        })
        
        # Create OTLP exporter
        otlp_exporter = OTLPSpanExporter(
            endpoint=settings.OTLP_ENDPOINT,
        )
        
        # Create tracer provider with batch span processor
        _tracer_provider = TracerProvider(resource=resource)
        span_processor = BatchSpanProcessor(otlp_exporter)
        _tracer_provider.add_span_processor(span_processor)
        
        # Set as global tracer provider
        trace.set_tracer_provider(_tracer_provider)
        
        # Instrument FastAPI if app is provided
        if app is not None:
            FastAPIInstrumentor.instrument_app(app)
            logger.info(f"FastAPI instrumented with tracing")
        
        logger.info(f"Tracing initialized with endpoint: {settings.OTLP_ENDPOINT}")
        return _tracer_provider
        
    except Exception as e:
        logger.error(f"Failed to initialize tracing: {e}")
        return None


def get_tracer(name: str = None):
    """
    Get a tracer instance
    
    Args:
        name: Optional tracer name (defaults to service name)
        
    Returns:
        Tracer instance
    """
    tracer_name = name or settings.SERVICE_NAME
    return trace.get_tracer(tracer_name)


def shutdown_tracing():
    """
    Shutdown the tracer provider gracefully
    """
    global _tracer_provider
    
    if _tracer_provider:
        try:
            _tracer_provider.shutdown()
            logger.info("Tracing provider shut down")
        except Exception as e:
            logger.error(f"Error shutting down tracing: {e}")
