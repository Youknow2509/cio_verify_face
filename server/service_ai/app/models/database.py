"""
Database models
"""
from sqlalchemy import Column, String, DateTime, Boolean, Integer, ForeignKey, Text, ARRAY, Float
from sqlalchemy.dialects.postgresql import UUID, JSONB
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
from sqlalchemy.types import UserDefinedType
import uuid


class Vector(UserDefinedType):
    """Custom SQLAlchemy type for pgvector's vector type"""
    
    cache_ok = True
    
    def __init__(self, dim=512):
        self.dim = dim
    
    def get_col_spec(self, **kw):
        return f"vector({self.dim})"
    
    def bind_processor(self, dialect):
        def process(value):
            if value is None:
                return None
            # Convert list/array to PostgreSQL vector format
            if isinstance(value, (list, tuple)):
                return str(list(value))
            return str(value)
        return process
    
    def result_processor(self, dialect, coltype):
        def process(value):
            if value is None:
                return None
            # PostgreSQL returns vector as string, we keep it as is
            # The application layer will handle conversion to numpy if needed
            return value
        return process


Base = declarative_base()


class FaceProfile(Base):
    """Face profile model for storing embeddings"""
    __tablename__ = "face_profiles"
    __table_args__ = (
        {'postgresql_partition_by': 'LIST (company_id)'},
    )
    
    profile_id = Column(UUID(as_uuid=True), primary_key=True, nullable=False, default=uuid.uuid4)
    company_id = Column(UUID(as_uuid=True), primary_key=True, nullable=False)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    
    # Embedding data - using pgvector's vector type for efficient similarity search
    # Falls back to ARRAY(Float) if pgvector is not available
    embedding = Column(Vector(512), nullable=False)
    embedding_version = Column(String(50), nullable=False, index=True)
    
    # Image reference (optional)
    enroll_image_path = Column(Text, nullable=True)
    
    # Status and metadata
    is_primary = Column(Boolean, default=False, nullable=False)
    quality_score = Column(Float, nullable=True)
    meta_data = Column(JSONB, default={}, nullable=False)  # Use meta_data to avoid SQLAlchemy reserved name
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False)
    deleted_at = Column(DateTime(timezone=True), nullable=True)
    
    # Indexing status
    indexed = Column(Boolean, default=False, nullable=False)
    index_version = Column(Integer, default=0, nullable=False)
