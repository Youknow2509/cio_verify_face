"""
PostgreSQL Database Manager Module (backed by MilvusManager for vector operations)
"""

import logging
from app.database.milvus_manager import MilvusManager

logger = logging.getLogger(__name__)


class PGManager(MilvusManager):
    """Compatibility wrapper that reuses MilvusManager implementation."""
    pass
from sqlalchemy import text
