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
    
    async def EnrollFaceStream(self, request_iterator, context):
        """Enroll a face with streaming (optimized for bandwidth)"""
        try:
            metadata = None
            image_chunks = []
            total_received = 0
            
            # Process streaming requests
            async for request in request_iterator:
                if request.HasField('metadata'):
                    # First message contains metadata
                    metadata = request.metadata
                    logger.info(f"Received enrollment metadata for user {metadata.user_id}, "
                              f"expected size: {metadata.total_size} bytes")
                    
                elif request.HasField('image_chunk'):
                    # Subsequent messages contain image chunks
                    chunk = request.image_chunk
                    image_chunks.append(chunk)
                    total_received += len(chunk)
                    
                    # Optional: Log progress
                    if metadata and metadata.total_size > 0:
                        progress = (total_received / metadata.total_size) * 100
                        if progress % 25 < 1:  # Log at 25%, 50%, 75%, 100%
                            logger.info(f"Received {progress:.1f}% of image data")
            
            # Validate metadata received
            if not metadata:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Missing metadata in stream")
                return face_service_pb2.EnrollResponse()
            
            # Validate image chunks received
            if not image_chunks:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("No image data received")
                return face_service_pb2.EnrollResponse()
            
            # Reconstruct complete image from chunks
            image_data = b''.join(image_chunks)
            logger.info(f"Reconstructed image: {len(image_data)} bytes "
                       f"(expected: {metadata.total_size} bytes)")
            
            # Verify size matches (with some tolerance for compression)
            if metadata.total_size > 0:
                size_diff = abs(len(image_data) - metadata.total_size)
                if size_diff > 1024:  # Allow 1KB tolerance
                    logger.warning(f"Image size mismatch: received {len(image_data)}, "
                                 f"expected {metadata.total_size}")
            
            # Convert image bytes to numpy array
            image = self._bytes_to_image(image_data)
            image_base64 = self._image_to_base64(image)
            
            # Prepare metadata dict
            enroll_metadata = dict(metadata.metadata) if metadata.metadata else {}
            enroll_metadata['image_format'] = metadata.image_format
            enroll_metadata['streaming'] = 'true'
            
            # Call face service
            result = await self.face_service.enroll_face(
                user_id=UUID(metadata.user_id),
                company_id=UUID(metadata.company_id),
                image_base64=image_base64,
                device_id=metadata.device_id if metadata.HasField('device_id') else None,
                make_primary=metadata.make_primary,
                metadata=enroll_metadata
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
            
            logger.info(f"Successfully enrolled face via streaming for user {metadata.user_id}")
            return response
            
        except Exception as e:
            logger.error(f"Error in EnrollFaceStream: {e}")
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
    
    async def VerifyFaceStream(self, request_iterator, context):
        """Verify a face with streaming (optimized for bandwidth)"""
        try:
            metadata = None
            image_chunks = []
            total_received = 0
            
            # Process streaming requests
            async for request in request_iterator:
                if request.HasField('metadata'):
                    # First message contains metadata
                    metadata = request.metadata
                    logger.info(f"Received verification metadata for company {metadata.company_id}, "
                              f"expected size: {metadata.total_size} bytes")
                    
                elif request.HasField('image_chunk'):
                    # Subsequent messages contain image chunks
                    chunk = request.image_chunk
                    image_chunks.append(chunk)
                    total_received += len(chunk)
                    
                    # Optional: Log progress
                    if metadata and metadata.total_size > 0:
                        progress = (total_received / metadata.total_size) * 100
                        if progress % 25 < 1:
                            logger.info(f"Received {progress:.1f}% of image data")
            
            # Validate metadata received
            if not metadata:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Missing metadata in stream")
                return face_service_pb2.VerifyResponse()
            
            # Validate image chunks received
            if not image_chunks:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("No image data received")
                return face_service_pb2.VerifyResponse()
            
            # Reconstruct complete image from chunks
            image_data = b''.join(image_chunks)
            logger.info(f"Reconstructed image: {len(image_data)} bytes")
            
            # Convert image bytes to numpy array
            image = self._bytes_to_image(image_data)
            image_base64 = self._image_to_base64(image)
            
            # Call face service
            result = await self.face_service.verify_face(
                image_base64=image_base64,
                company_id=UUID(metadata.company_id),
                user_id=UUID(metadata.user_id) if metadata.HasField('user_id') else None,
                device_id=metadata.device_id if metadata.HasField('device_id') else None,
                search_mode=metadata.search_mode,
                top_k=metadata.top_k
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
            
            logger.info(f"Successfully verified face via streaming")
            return response
            
        except Exception as e:
            logger.error(f"Error in VerifyFaceStream: {e}")
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
    
    async def VerifyFaceMultiFrameStream(self, request_iterator, context):
        """Verify face with multiple frames using streaming (optimized for bandwidth)"""
        try:
            metadata = None
            frames = []  # List of complete frames
            current_frame_chunks = []  # Chunks for current frame
            current_frame_info = None
            
            # Process streaming requests
            async for request in request_iterator:
                if request.HasField('metadata'):
                    # First message contains metadata
                    metadata = request.metadata
                    logger.info(f"Received multi-frame verification metadata, "
                              f"frame count: {metadata.frame_count}")
                    
                elif request.HasField('frame_delimiter'):
                    # Save previous frame if exists
                    if current_frame_chunks:
                        frame_data = b''.join(current_frame_chunks)
                        frames.append(frame_data)
                        logger.info(f"Completed frame {len(frames)}: {len(frame_data)} bytes")
                        current_frame_chunks = []
                    
                    # Start new frame
                    current_frame_info = request.frame_delimiter
                    
                elif request.HasField('frame_chunk'):
                    # Add chunk to current frame
                    current_frame_chunks.append(request.frame_chunk)
            
            # Save last frame
            if current_frame_chunks:
                frame_data = b''.join(current_frame_chunks)
                frames.append(frame_data)
                logger.info(f"Completed final frame {len(frames)}: {len(frame_data)} bytes")
            
            # Validate metadata received
            if not metadata:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Missing metadata in stream")
                return face_service_pb2.VerifyResponse()
            
            # Validate frames
            if not frames or len(frames) < 3:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details(f"At least 3 frames required, received {len(frames)}")
                return face_service_pb2.VerifyResponse()
            
            # Process each frame and get verification results
            all_results = []
            for i, frame_data in enumerate(frames[:5]):  # Limit to 5 frames
                image = self._bytes_to_image(frame_data)
                image_base64 = self._image_to_base64(image)
                
                result = await self.face_service.verify_face(
                    image_base64=image_base64,
                    company_id=UUID(metadata.company_id),
                    user_id=UUID(metadata.user_id) if metadata.HasField('user_id') else None,
                    device_id=metadata.device_id if metadata.HasField('device_id') else None,
                    search_mode=metadata.search_mode,
                    top_k=metadata.top_k
                )
                
                if result.verified:
                    all_results.append(result)
                    logger.info(f"Frame {i+1} verified successfully")
            
            # If no frames verified, return no match
            if not all_results:
                return face_service_pb2.VerifyResponse(
                    status="no_match",
                    verified=False,
                    message="No frames could be verified"
                )
            
            # Aggregate results - use majority voting and average confidence
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
                avg_confidence = np.mean(user_confidences[best_user_id])
                
                # Find the best match from results
                best_result = None
                for result in all_results:
                    if result.best_match and str(result.best_match.user_id) == best_user_id:
                        best_result = result
                        break
                
                if best_result:
                    response = face_service_pb2.VerifyResponse(
                        status="match",
                        verified=True,
                        message=f"Verified with {user_votes[best_user_id]} out of {len(frames)} frames (streaming)"
                    )
                    
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
                    
                    logger.info(f"Successfully verified multi-frame via streaming")
                    return response
            
            return face_service_pb2.VerifyResponse(
                status="no_match",
                verified=False,
                message="Could not verify face across multiple frames"
            )
            
        except Exception as e:
            logger.error(f"Error in VerifyFaceMultiFrameStream: {e}")
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
    
    async def UpdateProfileStream(self, request_iterator, context):
        """Update a face profile with streaming (optimized for bandwidth)"""
        try:
            metadata = None
            image_chunks = []
            total_received = 0
            
            # Process streaming requests
            async for request in request_iterator:
                if request.HasField('metadata'):
                    # First message contains metadata
                    metadata = request.metadata
                    logger.info(f"Received update profile metadata for profile {metadata.profile_id}, "
                              f"has_image: {metadata.has_image}")
                    
                elif request.HasField('image_chunk'):
                    # Subsequent messages contain image chunks
                    chunk = request.image_chunk
                    image_chunks.append(chunk)
                    total_received += len(chunk)
                    
                    # Optional: Log progress
                    if metadata and metadata.has_image and metadata.total_size > 0:
                        progress = (total_received / metadata.total_size) * 100
                        if progress % 25 < 1:
                            logger.info(f"Received {progress:.1f}% of image data")
            
            # Validate metadata received
            if not metadata:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Missing metadata in stream")
                return face_service_pb2.UpdateProfileResponse()
            
            # Process image if provided
            image_base64 = None
            if metadata.has_image:
                if not image_chunks:
                    context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                    context.set_details("No image data received but has_image=true")
                    return face_service_pb2.UpdateProfileResponse()
                
                # Reconstruct complete image from chunks
                image_data = b''.join(image_chunks)
                logger.info(f"Reconstructed image: {len(image_data)} bytes")
                
                # Convert image bytes to numpy array
                image = self._bytes_to_image(image_data)
                image_base64 = self._image_to_base64(image)
            
            # Prepare metadata dict
            update_metadata = dict(metadata.metadata) if metadata.metadata else None
            if update_metadata and metadata.has_image:
                update_metadata['image_format'] = metadata.image_format
                update_metadata['streaming'] = 'true'
            
            # Call face service
            result = await self.face_service.update_profile(
                profile_id=UUID(metadata.profile_id),
                company_id=UUID(metadata.company_id),
                image_base64=image_base64,
                make_primary=metadata.make_primary if metadata.HasField('make_primary') else None,
                metadata=update_metadata
            )
            
            logger.info(f"Successfully updated profile via streaming: {metadata.profile_id}")
            return face_service_pb2.UpdateProfileResponse(
                status=result["status"],
                message=result["message"]
            )
            
        except Exception as e:
            logger.error(f"Error in UpdateProfileStream: {e}")
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
