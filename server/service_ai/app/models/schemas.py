"""
Pydantic models for API requests and responses
"""
from pydantic import BaseModel, Field, validator
from typing import Optional, List, Dict, Any, Union
from datetime import datetime
from uuid import UUID
import base64

class SessionUser(BaseModel):
    """Model representing authenticated user from JWT"""
    user_id: UUID = Field(..., description="User ID")
    role: int = Field(..., description="User role level")
    company_id: UUID = Field(..., description="Company ID")
    session_id: UUID = Field(..., description="Session ID")
    exprires_at: datetime = Field(..., description="Token expiration time")

class EnrollRequest(BaseModel):
    """Request model for face enrollment"""
    user_id: UUID = Field(..., description="User ID to enroll face for")
    company_id: UUID = Field(..., description="Company ID for data partitioning")
    image_base64: str = Field(..., description="Base64 encoded image")
    device_id: Optional[str] = Field(None, description="Device ID that captured the image")
    make_primary: bool = Field(False, description="Set this profile as primary")
    metadata: Optional[Dict[str, Any]] = Field(default_factory=dict, description="Additional metadata")
    
    @validator('image_base64')
    def validate_base64(cls, v):
        """Validate base64 string"""
        try:
            base64.b64decode(v)
            return v
        except Exception:
            raise ValueError("Invalid base64 string")

class EnrollResponse(BaseModel):
    """Response model for face enrollment"""
    status: str = Field(..., description="Status: ok, duplicate, failed")
    profile_id: Optional[UUID] = Field(None, description="Created profile ID")
    message: Optional[str] = Field(None, description="Status message")
    duplicate_profiles: Optional[List[Dict[str, Any]]] = Field(None, description="Matched profiles if duplicate detected")
    quality_score: Optional[float] = Field(None, description="Image quality score")

class VerifyRequest(BaseModel):
    """Request model for face verification"""
    image_base64: str = Field(..., description="Base64 encoded image")
    company_id: UUID = Field(..., description="Company ID for data partitioning")
    user_id: Optional[UUID] = Field(None, description="User ID for 1:1 verification")
    device_id: Optional[str] = Field(None, description="Device ID")
    search_mode: str = Field("1:N", description="Search mode: 1:1 or 1:N")
    top_k: int = Field(5, description="Number of top matches to return")
    
    @validator('image_base64')
    def validate_base64(cls, v):
        """Validate base64 string"""
        try:
            base64.b64decode(v)
            return v
        except Exception:
            raise ValueError("Invalid base64 string")
    
    @validator('search_mode')
    def validate_search_mode(cls, v):
        """Validate search mode"""
        if v not in ["1:1", "1:N"]:
            raise ValueError("search_mode must be '1:1' or '1:N'")
        return v

class VerifyMatch(BaseModel):
    """Single match result"""
    user_id: UUID
    profile_id: UUID
    similarity: float
    confidence: float
    is_primary: bool

class VerifyResponse(BaseModel):
    """Response model for face verification"""
    status: str = Field(..., description="Status: match, no_match, failed")
    verified: bool = Field(..., description="Whether verification passed")
    matches: List[VerifyMatch] = Field(default_factory=list, description="List of matches")
    best_match: Optional[VerifyMatch] = Field(None, description="Best match if found")
    message: Optional[str] = Field(None, description="Status message")
    liveness_score: Optional[float] = Field(None, description="Liveness detection score")

class UpdateProfileRequest(BaseModel):
    """Request model for updating face profile"""
    image_base64: Optional[str] = Field(None, description="New base64 encoded image")
    make_primary: Optional[bool] = Field(None, description="Set as primary profile")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Update metadata")
    
    @validator('image_base64')
    def validate_base64(cls, v):
        """Validate base64 string"""
        if v is None:
            return v
        try:
            base64.b64decode(v)
            return v
        except Exception:
            raise ValueError("Invalid base64 string")

class FaceProfileResponse(BaseModel):
    """Response model for face profile"""
    profile_id: UUID
    user_id: UUID
    company_id: UUID
    embedding_version: str
    is_primary: bool
    created_at: datetime
    updated_at: datetime
    deleted_at: Optional[datetime]
    metadata: Dict[str, Any]
    quality_score: Optional[float]

class ReindexRequest(BaseModel):
    """Request model for reindexing"""
    force: bool = Field(False, description="Force rebuild even if recent rebuild exists")
    embedding_version: Optional[str] = Field(None, description="Specific embedding version to reindex")

class ReindexResponse(BaseModel):
    """Response model for reindex operation"""
    status: str
    message: str
    profiles_indexed: int
    duration_seconds: float

class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    service: str
    version: str
    environment: str
    models_loaded: bool
    index_size: int

class ImageOptimizationInfo(BaseModel):
    """Information about image optimization"""
    original_size: int
    optimized_size: int
    reduction_percent: float
    format: str
