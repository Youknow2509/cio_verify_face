"""
Logging configuration
"""
import logging
import sys
from pythonjsonlogger import jsonlogger

from app.core.config import settings


def setup_logging():
    """Setup logging configuration"""
    
    # Create logger
    logger = logging.getLogger()
    logger.setLevel(getattr(logging, settings.LOG_LEVEL.upper()))
        
    # Remove existing handlers
    logger.handlers = []
    
    # Create console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(getattr(logging, settings.LOG_LEVEL.upper()))
    
    # Create formatter
    if settings.LOG_FORMAT == "json":
        formatter = jsonlogger.JsonFormatter(
            '%(asctime)s %(name)s %(levelname)s %(message)s',
            timestamp=True
        )
    else:
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    
    # Create file handler
    file_handler = logging.FileHandler(settings.LOG_FILEPATH, mode="a", encoding="utf-8")
    file_handler.setLevel(getattr(logging, settings.LOG_LEVEL.upper()))
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)
    
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    return logger
