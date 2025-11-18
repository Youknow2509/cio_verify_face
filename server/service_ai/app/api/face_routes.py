"""
Face verification API routes
"""
from fastapi import APIRouter, HTTPException, Depends, Request, File, UploadFile, Form
from typing import List, Optional
from uuid import UUID
import logging
import base64

from app.models.schemas import (
    EnrollRequest, EnrollResponse,
    VerifyRequest, VerifyResponse,
    UpdateProfileRequest, FaceProfileResponse,
    ReindexRequest, ReindexResponse
)
from app.services.face_service import FaceService
from app.dependencies.auth_dependencies import get_current_user, get_current_auth, get_current_device

logger = logging.getLogger(__name__)

router = APIRouter()


def get_face_service(request: Request) -> FaceService:
    """Dependency to get face service instance"""
    return request.app.state.face_service


@router.post("/enroll", response_model=EnrollResponse, summary="Enroll new face")
async def enroll_face(
    request: EnrollRequest,
    face_service: FaceService = Depends(get_face_service),
    user: dict = Depends(get_current_user),
):
    """
    Enroll a new face profile for a user
    
    - **user_id**: UUID of the user
    - **image_base64**: Base64 encoded face image
    - **device_id**: Optional device identifier
    - **make_primary**: Set this profile as primary (default: False)
    - **metadata**: Optional metadata dict
    
    Returns:
    - **status**: "ok", "duplicate", or "failed"
    - **profile_id**: UUID of created profile (if successful)
    - **message**: Status message
    - **duplicate_profiles**: List of matching profiles (if duplicate detected)
    - **quality_score**: Image quality score
    """
    try:
        # `get_current_user` ensures request is authenticated as a user
        # and `user` contains parsed token info. Use user's id when
        # creating/updating profiles.
        result = await face_service.enroll_face(
            user_id=request.user_id,
            company_id=request.company_id,
            image_base64=request.image_base64,
            device_id=request.device_id,
            make_primary=request.make_primary,
            metadata=request.metadata
        )
        return result
    except Exception as e:
        logger.error(f"Error in enroll endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/verify", response_model=VerifyResponse, summary="Verify face")
async def verify_face(
    request: VerifyRequest,
    face_service: FaceService = Depends(get_face_service),
    auth: dict = Depends(get_current_auth),
):
    """
    Verify a face against enrolled profiles
    
    - **image_base64**: Base64 encoded face image
    - **user_id**: Optional user ID for 1:1 verification
    - **device_id**: Optional device identifier
    - **search_mode**: "1:1" or "1:N" (default: "1:N")
    - **top_k**: Number of top matches to return (default: 5)
    
    Returns:
    - **status**: "match", "no_match", or "failed"
    - **verified**: Boolean indicating if verification passed
    - **matches**: List of matching profiles
    - **best_match**: Best matching profile (if verified)
    - **message**: Status message
    - **liveness_score**: Liveness detection score
    """
    try:
        # `get_current_auth` enforces the request is authenticated as
        # either a device or a user; routes can examine `auth['auth_type']`.
        result = await face_service.verify_face(
            image_base64=request.image_base64,
            company_id=request.company_id,
            user_id=request.user_id,
            device_id=request.device_id,
            search_mode=request.search_mode,
            top_k=request.top_k
        )
        return result
    except Exception as e:
        logger.error(f"Error in verify endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.put("/profile/{profile_id}", summary="Update face profile")
async def update_profile(
    profile_id: UUID,
    company_id: UUID,
    request: UpdateProfileRequest,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Update an existing face profile
    
    - **profile_id**: UUID of the profile to update
    - **company_id**: UUID of the company
    - **image_base64**: Optional new face image
    - **make_primary**: Optional flag to set as primary
    - **metadata**: Optional metadata to update
    
    Returns status and message
    """
    try:
        result = await face_service.update_profile(
            profile_id=profile_id,
            company_id=company_id,
            image_base64=request.image_base64,
            make_primary=request.make_primary,
            metadata=request.metadata
        )
        
        if result["status"] == "failed":
            raise HTTPException(status_code=400, detail=result["message"])
        
        return result
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in update_profile endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/profile/{profile_id}", summary="Delete face profile")
async def delete_profile(
    profile_id: UUID,
    company_id: UUID,
    hard_delete: bool = False,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Delete a face profile (soft delete by default)
    
    - **profile_id**: UUID of the profile to delete
    - **company_id**: UUID of the company
    - **hard_delete**: If true, permanently delete (default: False)
    
    By default, profiles are soft-deleted and retained for the configured
    retention period before permanent deletion.
    """
    try:
        result = await face_service.delete_profile(
            profile_id=profile_id,
            company_id=company_id,
            hard_delete=hard_delete
        )
        
        if result["status"] == "failed":
            raise HTTPException(status_code=404, detail=result["message"])
        
        return result
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in delete_profile endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/profiles/{user_id}", response_model=List[FaceProfileResponse], 
           summary="Get user face profiles")
async def get_user_profiles(
    user_id: UUID,
    company_id: UUID,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Get all face profiles for a user
    
    - **user_id**: UUID of the user
    - **company_id**: UUID of the company
    
    Returns list of face profiles
    """
    try:
        profiles = await face_service.get_user_profiles(user_id, company_id)
        return profiles
    except Exception as e:
        logger.error(f"Error in get_user_profiles endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/reindex", response_model=ReindexResponse, 
            summary="Rebuild FAISS index (admin)")
async def reindex(
    request: ReindexRequest,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Rebuild the FAISS index from database
    
    - **force**: Force rebuild even if recent rebuild exists (default: False)
    - **embedding_version**: Optional specific version to reindex
    
    This operation may take time depending on the number of profiles.
    Should be run during low-traffic periods.
    
    Returns:
    - **status**: "ok", "skipped", or "failed"
    - **message**: Status message
    - **profiles_indexed**: Number of profiles indexed
    - **duration_seconds**: Time taken for reindexing
    """
    try:
        result = await face_service.reindex(force=request.force)
        return result
    except Exception as e:
        logger.error(f"Error in reindex endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/cleanup", summary="Cleanup soft-deleted profiles (admin)")
async def cleanup_soft_deleted(
    face_service: FaceService = Depends(get_face_service)
):
    """
    Cleanup soft-deleted profiles past retention period
    
    This removes profiles that have been soft-deleted and are past the
    configured retention period (default: 30 days).
    """
    try:
        await face_service.cleanup_soft_deleted()
        return {"status": "ok", "message": "Cleanup completed"}
    except Exception as e:
        logger.error(f"Error in cleanup endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ==================== OPTIMIZED ENDPOINTS (Multipart/Form-Data) ====================
# These endpoints accept binary image files instead of base64, reducing bandwidth by ~33%

@router.post("/enroll/upload", response_model=EnrollResponse, 
            summary="Enroll new face (optimized - binary upload)")
async def enroll_face_upload(
    image: UploadFile = File(..., description="Face image file (JPEG/PNG)"),
    user_id: UUID = Form(..., description="User ID to enroll face for"),
    company_id: UUID = Form(..., description="Company ID"),
    device_id: Optional[str] = Form(None, description="Device ID"),
    make_primary: bool = Form(False, description="Set as primary profile"),
    face_service: FaceService = Depends(get_face_service)
):
    """
    **OPTIMIZED ENDPOINT**: Enroll a new face using multipart/form-data.
    
    This endpoint is more bandwidth-efficient than the base64 version:
    - Reduces data transfer by ~33% (no base64 encoding overhead)
    - Better for mobile and low-bandwidth environments
    - Accepts JPEG, PNG, and other image formats
    
    **Usage:**
    ```
    POST /api/v1/face/enroll/upload
    Content-Type: multipart/form-data
    
    Form fields:
    - image: [binary file]
    - user_id: "550e8400-e29b-41d4-a716-446655440000"
    - device_id: "device_001" (optional)
    - make_primary: true (optional, default: false)
    ```
    
    **Advantages over /enroll:**
    - 33% less bandwidth usage
    - Faster upload on slow connections
    - Native file upload support in browsers/apps
    """
    try:
        # Read image file
        image_bytes = await image.read()
        
        # Convert to base64 for internal processing
        # (This maintains compatibility with existing service logic)
        image_base64 = base64.b64encode(image_bytes).decode('utf-8')
        
        result = await face_service.enroll_face(
            user_id=user_id,
            company_id=company_id,
            image_base64=image_base64,
            device_id=device_id,
            make_primary=make_primary,
            metadata={"upload_type": "binary", "filename": image.filename}
        )
        return result
    except Exception as e:
        logger.error(f"Error in enroll_upload endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/verify/upload", response_model=VerifyResponse,
            summary="Verify face (optimized - binary upload)")
async def verify_face_upload(
    image: UploadFile = File(..., description="Face image file (JPEG/PNG)"),
    company_id: UUID = Form(..., description="Company ID"),
    user_id: Optional[UUID] = Form(None, description="User ID for 1:1 verification"),
    device_id: Optional[str] = Form(None, description="Device ID"),
    search_mode: str = Form("1:N", description="Search mode: 1:1 or 1:N"),
    top_k: int = Form(5, description="Number of top matches"),
    face_service: FaceService = Depends(get_face_service)
):
    """
    **OPTIMIZED ENDPOINT**: Verify a face using multipart/form-data.
    
    This endpoint is more bandwidth-efficient than the base64 version:
    - Reduces data transfer by ~33% (no base64 encoding overhead)
    - Better for mobile and low-bandwidth environments
    - Accepts JPEG, PNG, and other image formats
    
    **Usage:**
    ```
    POST /api/v1/face/verify/upload
    Content-Type: multipart/form-data
    
    Form fields:
    - image: [binary file]
    - user_id: "550e8400-..." (optional, for 1:1 mode)
    - device_id: "device_001" (optional)
    - search_mode: "1:N" (optional, default: "1:N")
    - top_k: 5 (optional, default: 5)
    ```
    
    **Advantages over /verify:**
    - 33% less bandwidth usage
    - Faster verification on slow connections
    - Native file upload support
    """
    try:
        # Validate search mode
        if search_mode not in ["1:1", "1:N"]:
            raise HTTPException(status_code=400, detail="search_mode must be '1:1' or '1:N'")
        
        # Read image file
        image_bytes = await image.read()
        
        # Convert to base64 for internal processing
        image_base64 = base64.b64encode(image_bytes).decode('utf-8')
        
        result = await face_service.verify_face(
            image_base64=image_base64,
            company_id=company_id,
            user_id=user_id,
            device_id=device_id,
            search_mode=search_mode,
            top_k=top_k
        )
        return result
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in verify_upload endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ==================== AUDIT LOG ENDPOINTS ====================
# Query audit logs from ScyllaDB for monitoring and compliance

@router.get("/audit/logs", summary="Get audit logs (admin)")
async def get_audit_logs(
    company_id: UUID,
    operation: Optional[str] = None,
    limit: int = 100,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Get audit logs from ScyllaDB
    
    **Query Parameters:**
    - **operation**: Filter by operation type (enroll, verify, update, delete) - optional
    - **limit**: Maximum number of logs to return (default: 100, max: 1000)
    
    **Returns:**
    List of audit log entries with timestamps, operations, and results
    
    **Use Cases:**
    - Monitor system activity
    - Compliance and auditing
    - Debugging issues
    - Security analysis
    
    **Note**: Logs are stored in ScyllaDB for better performance on time-series queries
    """
    try:
        if limit > 1000:
            limit = 1000
        
        if not face_service.scylladb:
            raise HTTPException(status_code=503, detail="ScyllaDB not available")
        
        # Validate operation if provided
        if operation and operation not in ["enroll", "verify", "update", "delete"]:
            raise HTTPException(
                status_code=400,
                detail="operation must be one of: enroll, verify, update, delete"
            )
        
        logs = face_service.scylladb.get_audit_logs(company_id=company_id, operation=operation, limit=limit)
        
        return {
            "status": "ok",
            "count": len(logs),
            "logs": logs
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in get_audit_logs endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/audit/user/{user_id}", summary="Get user audit logs (admin)")
async def get_user_audit_logs(
    user_id: UUID,
    company_id: UUID,
    limit: int = 100,
    face_service: FaceService = Depends(get_face_service)
):
    """
    Get audit logs for a specific user from ScyllaDB
    
    **Path Parameters:**
    - **user_id**: UUID of the user
    
    **Query Parameters:**
    - **limit**: Maximum number of logs to return (default: 100, max: 1000)
    
    **Returns:**
    List of audit log entries for the specified user, ordered by timestamp (newest first)
    
    **Use Cases:**
    - Track user activity history
    - Investigate user-specific issues
    - Compliance reporting per user
    - User behavior analysis
    
    **Performance:**
    This query is optimized with ScyllaDB's user_audit_logs table for fast retrieval
    """
    try:
        if limit > 1000:
            limit = 1000
        
        if not face_service.scylladb:
            raise HTTPException(status_code=503, detail="ScyllaDB not available")
        
        logs = face_service.scylladb.get_user_audit_logs(user_id=user_id, company_id=company_id, limit=limit)
        
        return {
            "status": "ok",
            "user_id": str(user_id),
            "count": len(logs),
            "logs": logs
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in get_user_audit_logs endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))
