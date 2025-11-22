"""
Face verification API routes
"""
from fastapi import APIRouter, HTTPException, Depends, Request, File, UploadFile, Form
from typing import List, Optional
from uuid import UUID
import logging
import base64

from app.auth.auth import is_manager_or_admin
from app.auth.auth_bearer import JWTBearer
from app.models.schemas import (
    EnrollRequest, EnrollResponse, SessionUser,
    VerifyRequest, VerifyResponse,
    UpdateProfileRequest, FaceProfileResponse,
    ReindexRequest, ReindexResponse
)
from app.services.face_service import FaceService

logger = logging.getLogger(__name__)

router = APIRouter()

def get_face_service(request: Request) -> FaceService:
    """Dependency to get face service instance"""
    return request.app.state.face_service

@router.put("/profile/{profile_id}/upload", summary="Update face profile")
async def update_profile_upload(
    profile_id: UUID,
    company_id: UUID = Form(..., description="Company ID"),
    image: Optional[UploadFile] = File(None, description="Face image file (JPEG/PNG)"),
    make_primary: Optional[bool] = Form(None, description="Set as primary profile"),
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Update an existing face profile using multipart/form-data.
    **Usage:**
    ```
    PUT /api/v1/face/profile/{profile_id}/upload
    Content-Type: multipart/form-data
    Headers:
    - Authorization: Bearer <token>
    Form fields:
    - company_id: "550e8400-e29b-41d4-a716-446655440000" (required)
    - image: [binary file] (optional)
    - make_primary: true (optional)
    ```
    Returns status and message
    """
    # Check permissions
    if is_manager_or_admin(company_id=company_id, token_payload=token_payload) is False:
        raise HTTPException(
            status_code=403, 
            detail="Insufficient permissions to update profile for this company"
        )
    try:
        # Convert image to base64 if provided
        image_base64 = None
        metadata = {
            "session_user": token_payload.dict() if token_payload else None
        }
        if image:
            image_bytes = await image.read()
            image_base64 = base64.b64encode(image_bytes).decode('utf-8')
            metadata.update({
                "upload_type": "binary",
                "filename": image.filename
            })
        result = await face_service.update_profile(
            profile_id=profile_id,
            company_id=company_id,
            image_base64=image_base64,
            make_primary=make_primary,
            metadata=metadata
        )
        if result["status"] == "failed":
            raise HTTPException(status_code=400, detail=result["message"])
        return result
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error in update_profile_upload endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.delete("/profile/{profile_id}", summary="Delete face profile")
async def delete_profile(
    profile_id: UUID,
    company_id: UUID,
    hard_delete: bool = False,
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Delete a face profile (soft delete by default)
    
    - **profile_id**: UUID of the profile to delete
    - **company_id**: UUID of the company
    - **hard_delete**: If true, permanently delete (default: False)
    """
    # Check permissions
    if is_manager_or_admin(company_id=company_id, token_payload=token_payload) is False:
        raise HTTPException(status_code=403, detail="Insufficient permissions to delete profile for this company")
    # Perform deletion
    try:
        result = await face_service.delete_profile(
            profile_id=profile_id,
            company_id=company_id,
            metadata={
                "session_user": token_payload.dict() if token_payload else None    
            },
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
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Get all face profiles for a user
    - **user_id**: UUID of the user
    - **company_id**: UUID of the company
    - **Authentication**: Bearer <token>
    
    Returns list of face profiles
    """
    # Check permissions
    if is_manager_or_admin(company_id=company_id, token_payload=token_payload) is False:
        raise HTTPException(status_code=403, detail="Insufficient permissions to view profiles for this company")
    # TODO: Check user_id in company_id
    # Fetch profiles
    try:
        profiles = await face_service.get_user_profiles(user_id, company_id)
        return profiles
    except Exception as e:
        logger.error(f"Error in get_user_profiles endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/enroll/upload", response_model=EnrollResponse,
            summary="Enroll new face")
async def enroll_face_upload(
    image: UploadFile = File(..., description="Face image file (JPEG/PNG)"),
    user_id: UUID = Form(..., description="User ID to enroll face for"),
    company_id: UUID = Form(..., description="Company ID"),
    device_id: Optional[str] = Form(None, description="Device ID"),
    make_primary: bool = Form(False, description="Set as primary profile"),
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Enroll a new face using multipart/form-data.
    **Usage:**
    ```
    POST /api/v1/face/enroll/upload
    Content-Type: multipart/form-data
    Headers:
    - Authorization: Bearer <token>
    Form fields:
    - image: [binary file]
    - user_id: "550e8400-e29b-41d4-a716-446655440000"
    - device_id: "device_001" (optional)
    - make_primary: true (optional, default: false)
    ```
    """
    # Check permissions
    if is_manager_or_admin(company_id=company_id, token_payload=token_payload) is False:
        raise HTTPException(status_code=403, detail="Insufficient permissions to enroll face for this company")
    # TODO: Check user_id in company_id
    # Handle file upload 
    try:
        # Read image file
        image_bytes = await image.read()
        
        # Convert to base64 for internal processing
        image_base64 = base64.b64encode(image_bytes).decode('utf-8')
        
        result = await face_service.enroll_face(
            user_id=user_id,
            company_id=company_id,
            image_base64=image_base64,
            device_id=device_id,
            make_primary=make_primary,
            metadata={
                "upload_type": "binary", 
                "filename": image.filename,
                "session_user": token_payload.dict() if token_payload else None
            }
        )
        return result
    except Exception as e:
        logger.error(f"Error in enroll_upload endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/verify/upload", response_model=VerifyResponse,
            summary="Verify face")
async def verify_face_upload(
    image: UploadFile = File(..., description="Face image file (JPEG/PNG)"),
    company_id: UUID = Form(..., description="Company ID"),
    user_id: Optional[UUID] = Form(None, description="User ID for 1:1 verification"),
    device_id: Optional[str] = Form(None, description="Device ID"),
    search_mode: str = Form("1:N", description="Search mode: 1:1 or 1:N"),
    top_k: int = Form(5, description="Number of top matches"),
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Verify a face using multipart/form-data.
    **Usage:**
    ```
    POST /api/v1/face/verify/upload
    Content-Type: multipart/form-data
    Authentication: Bearer <token>
    Form fields:
    - image: [binary file]
    - user_id: "550e8400-..." (optional, for 1:1 mode)
    - device_id: "device_001" (optional)
    - search_mode: "1:N" (optional, default: "1:N")
    - top_k: 5 (optional, default: 5)
    ```
    """
    # Ensure authenticated
    if token_payload.user_id is None:
        raise HTTPException(status_code=403, detail="Authentication required for face verification")
    # Validate uuid
    if not isinstance(company_id, UUID):
        raise HTTPException(status_code=400, detail="Invalid company_id format")
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

@router.post("/cleanup/profiles", summary="Cleanup face profiles")
async def cleanup_face_profiles(
    company_id: UUID,
    face_service: FaceService = Depends(get_face_service),
    token_payload: SessionUser = Depends(JWTBearer()),
):
    """
    Cleanup unused or invalid face profiles for a company.
    
    - **company_id**: UUID of the company
    - **Authentication**: Bearer <token>
    
    Returns status and number of profiles cleaned up.
    """
    # Check permissions
    if is_manager_or_admin(company_id=company_id, token_payload=token_payload) is False:
        raise HTTPException(status_code=403, detail="Insufficient permissions to cleanup profiles for this company")
    try:
        result = await face_service.cleanup_profiles_for_company(company_id=company_id, metadata={"session_user": token_payload.dict() if token_payload else None})
        return result
    except Exception as e:
        logger.error(f"Error in cleanup_face_profiles endpoint: {e}")
        raise HTTPException(status_code=500, detail=str(e))