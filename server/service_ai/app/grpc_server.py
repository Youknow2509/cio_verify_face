"""
gRPC server implementation for face verification service
"""
import grpc
from concurrent import futures
import logging
from typing import List
import numpy as np
from io import BytesIO
from PIL import Image
import cv2

from app.grpc_generated import face_service_pb2, face_service_pb2_grpc
from app.services.face_service import FaceService
from app.core.config import settings
from uuid import UUID

logger = logging.getLogger(__name__)


class FaceVerificationServicer(face_service_pb2_grpc.FaceVerificationServicer):
    """gRPC servicer for face verification"""
    
    def __init__(self, face_service: FaceService):
        self.face_service = face_service
    
    def _bytes_to_image(self, image_bytes: bytes) -> np.ndarray:
        """Convert image bytes to numpy array"""
        try:
            image = Image.open(BytesIO(image_bytes))
            image_array = np.array(image)
            # Convert RGB to BGR for OpenCV
            if len(image_array.shape) == 3 and image_array.shape[2] == 3:
                image_array = cv2.cvtColor(image_array, cv2.COLOR_RGB2BGR)
            return image_array
        except Exception as e:
            logger.error(f"Error converting bytes to image: {e}")
            raise
    
    def _image_to_base64(self, image: np.ndarray) -> str:
        """Convert numpy array to base64 string"""
        import base64
        # Convert BGR to RGB
        if len(image.shape) == 3 and image.shape[2] == 3:
            image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
        # Encode as JPEG
        _, buffer = cv2.imencode('.jpg', image)
        return base64.b64encode(buffer).decode('utf-8')
    
    async def EnrollFace(self, request, context):
        """Enroll a face"""
        try:
            # Convert image bytes to numpy array
            image = self._bytes_to_image(request.image_data)
            image_base64 = self._image_to_base64(image)
            
            # Call face service
            result = await self.face_service.enroll_face(
                user_id=UUID(request.user_id),
                company_id=UUID(request.company_id),
                image_base64=image_base64,
                device_id=request.device_id if request.HasField('device_id') else None,
                make_primary=request.make_primary,
                metadata=dict(request.metadata) if request.metadata else {}
            )
            
            # Build response
            response = face_service_pb2.EnrollResponse(
                status=result.status,
                message=result.message or ""
            )
            
            if result.profile_id:
                response.profile_id = str(result.profile_id)
            
            if result.quality_score is not None:
                response.quality_score = result.quality_score
            
            if result.duplicate_profiles:
                for dup in result.duplicate_profiles:
                    response.duplicate_profiles.append(
                        face_service_pb2.DuplicateProfile(
                            user_id=dup['user_id'],
                            similarity=dup['similarity']
                        )
                    )
            
            return response
            
        except Exception as e:
            logger.error(f"Error in EnrollFace: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.EnrollResponse()
    
    async def VerifyFace(self, request, context):
        """Verify a face with single image"""
        try:
            # Convert image bytes to numpy array
            image = self._bytes_to_image(request.image_data)
            image_base64 = self._image_to_base64(image)
            
            # Call face service
            result = await self.face_service.verify_face(
                image_base64=image_base64,
                company_id=UUID(request.company_id),
                user_id=UUID(request.user_id) if request.HasField('user_id') else None,
                device_id=request.device_id if request.HasField('device_id') else None,
                search_mode=request.search_mode,
                top_k=request.top_k
            )
            
            # Build response
            response = face_service_pb2.VerifyResponse(
                status=result.status,
                verified=result.verified,
                message=result.message or ""
            )
            
            if result.liveness_score is not None:
                response.liveness_score = result.liveness_score
            
            # Add matches
            for match in result.matches:
                response.matches.append(
                    face_service_pb2.Match(
                        user_id=str(match.user_id),
                        profile_id=str(match.profile_id),
                        similarity=match.similarity,
                        confidence=match.confidence,
                        is_primary=match.is_primary
                    )
                )
            
            # Add best match
            if result.best_match:
                response.best_match.CopyFrom(
                    face_service_pb2.Match(
                        user_id=str(result.best_match.user_id),
                        profile_id=str(result.best_match.profile_id),
                        similarity=result.best_match.similarity,
                        confidence=result.best_match.confidence,
                        is_primary=result.best_match.is_primary
                    )
                )
            
            return response
            
        except Exception as e:
            logger.error(f"Error in VerifyFace: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.VerifyResponse()
    
    async def VerifyFaceMultiFrame(self, request, context):
        """Verify face with multiple frames (3-5 frames for robustness)"""
        try:
            if not request.frames or len(request.frames) < 3:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("At least 3 frames required")
                return face_service_pb2.VerifyResponse()
            
            # Process each frame and get verification results
            all_results = []
            for frame_bytes in request.frames[:5]:  # Limit to 5 frames
                image = self._bytes_to_image(frame_bytes)
                image_base64 = self._image_to_base64(image)
                
                result = await self.face_service.verify_face(
                    image_base64=image_base64,
                    company_id=UUID(request.company_id),
                    user_id=UUID(request.user_id) if request.HasField('user_id') else None,
                    device_id=request.device_id if request.HasField('device_id') else None,
                    search_mode=request.search_mode,
                    top_k=request.top_k
                )
                
                if result.verified:
                    all_results.append(result)
            
            # If no frames verified, return no match
            if not all_results:
                return face_service_pb2.VerifyResponse(
                    status="no_match",
                    verified=False,
                    message="No frames could be verified"
                )
            
            # Aggregate results - use majority voting and average confidence
            # Count matches per user
            user_votes = {}
            user_confidences = {}
            
            for result in all_results:
                if result.best_match:
                    user_id = str(result.best_match.user_id)
                    if user_id not in user_votes:
                        user_votes[user_id] = 0
                        user_confidences[user_id] = []
                    
                    user_votes[user_id] += 1
                    user_confidences[user_id].append(result.best_match.confidence)
            
            # Get user with most votes
            if user_votes:
                best_user = max(user_votes.items(), key=lambda x: (x[1], np.mean(user_confidences[x[0]])))
                best_user_id = best_user[0]
                
                # Get average confidence for best user
                avg_confidence = np.mean(user_confidences[best_user_id])
                
                # Find the best match from results
                best_result = None
                for result in all_results:
                    if result.best_match and str(result.best_match.user_id) == best_user_id:
                        best_result = result
                        break
                
                if best_result:
                    # Build response with aggregated confidence
                    response = face_service_pb2.VerifyResponse(
                        status="match",
                        verified=True,
                        message=f"Verified with {user_votes[best_user_id]} out of {len(request.frames)} frames"
                    )
                    
                    # Add best match with averaged confidence
                    response.best_match.CopyFrom(
                        face_service_pb2.Match(
                            user_id=str(best_result.best_match.user_id),
                            profile_id=str(best_result.best_match.profile_id),
                            similarity=best_result.best_match.similarity,
                            confidence=float(avg_confidence),
                            is_primary=best_result.best_match.is_primary
                        )
                    )
                    
                    # Add all unique matches
                    seen_users = set()
                    for result in all_results:
                        for match in result.matches:
                            user_id = str(match.user_id)
                            if user_id not in seen_users:
                                seen_users.add(user_id)
                                response.matches.append(
                                    face_service_pb2.Match(
                                        user_id=user_id,
                                        profile_id=str(match.profile_id),
                                        similarity=match.similarity,
                                        confidence=match.confidence,
                                        is_primary=match.is_primary
                                    )
                                )
                    
                    return response
            
            # Fallback to no match
            return face_service_pb2.VerifyResponse(
                status="no_match",
                verified=False,
                message="Could not verify face across multiple frames"
            )
            
        except Exception as e:
            logger.error(f"Error in VerifyFaceMultiFrame: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.VerifyResponse()
    
    async def GetUserProfiles(self, request, context):
        """Get all face profiles for a user"""
        try:
            profiles = await self.face_service.get_user_profiles(
                UUID(request.user_id),
                UUID(request.company_id)
            )
            
            response = face_service_pb2.GetProfilesResponse()
            for profile in profiles:
                response.profiles.append(
                    face_service_pb2.FaceProfile(
                        profile_id=str(profile.profile_id),
                        user_id=str(profile.user_id),
                        company_id=str(profile.company_id),
                        embedding_version=profile.embedding_version,
                        is_primary=profile.is_primary,
                        created_at=profile.created_at.isoformat(),
                        updated_at=profile.updated_at.isoformat(),
                        deleted_at=profile.deleted_at.isoformat() if profile.deleted_at else "",
                        metadata={k: str(v) for k, v in profile.meta_data.items()},
                        quality_score=profile.quality_score if profile.quality_score else 0.0
                    )
                )
            
            return response
            
        except Exception as e:
            logger.error(f"Error in GetUserProfiles: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.GetProfilesResponse()
    
    async def UpdateProfile(self, request, context):
        """Update a face profile"""
        try:
            image_base64 = None
            if request.HasField('image_data'):
                image = self._bytes_to_image(request.image_data)
                image_base64 = self._image_to_base64(image)
            
            result = await self.face_service.update_profile(
                profile_id=UUID(request.profile_id),
                company_id=UUID(request.company_id),
                image_base64=image_base64,
                make_primary=request.make_primary if request.HasField('make_primary') else None,
                metadata=dict(request.metadata) if request.metadata else None
            )
            
            return face_service_pb2.UpdateProfileResponse(
                status=result["status"],
                message=result["message"]
            )
            
        except Exception as e:
            logger.error(f"Error in UpdateProfile: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.UpdateProfileResponse()
    
    async def DeleteProfile(self, request, context):
        """Delete a face profile"""
        try:
            result = await self.face_service.delete_profile(
                profile_id=UUID(request.profile_id),
                company_id=UUID(request.company_id),
                hard_delete=request.hard_delete
            )
            
            return face_service_pb2.DeleteProfileResponse(
                status=result["status"],
                message=result["message"]
            )
            
        except Exception as e:
            logger.error(f"Error in DeleteProfile: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return face_service_pb2.DeleteProfileResponse()


async def serve_grpc(face_service: FaceService, port: int = 50051):
    """Start gRPC server"""
    server = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))
    face_service_pb2_grpc.add_FaceVerificationServicer_to_server(
        FaceVerificationServicer(face_service), server
    )
    
    # Enable reflection for grpcurl
    from grpc_reflection.v1alpha import reflection
    SERVICE_NAMES = (
        face_service_pb2.DESCRIPTOR.services_by_name['FaceVerification'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    
    server.add_insecure_port(f'[::]:{port}')
    logger.info(f"Starting gRPC server on port {port}...")
    await server.start()
    await server.wait_for_termination()
