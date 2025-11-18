"""
Main face service that orchestrates all face verification operations
"""
import logging
import base64
import numpy as np
import json
import ast
import cv2
from typing import Optional, List, Dict, Tuple
from uuid import UUID, uuid4
from datetime import datetime, timedelta
from sqlalchemy import create_engine, and_
from sqlalchemy.orm import sessionmaker, Session

from app.core.config import settings
from app.services.face_detector import FaceDetector
from app.services.face_embedding import FaceEmbedding
from app.services.pgvector_manager import PgVectorManager
from app.services.liveness_detector import LivenessDetector
from app.services.scylladb_manager import ScyllaDBManager
from app.services.minio_manager import MinIOManager
from app.utils.image_optimizer import ImageOptimizer
from app.models.database import Base, FaceProfile, FaceAuditLog
from app.models.schemas import (
    EnrollResponse, VerifyResponse, VerifyMatch,
    FaceProfileResponse, ReindexResponse
)

logger = logging.getLogger(__name__)


class FaceService:
    """Main face verification service"""
    
    def __init__(self):
        """Initialize face service"""
        logger.info("Initializing FaceService...")
        
        # Initialize components
        self.detector = FaceDetector()
        self.embedding_gen = FaceEmbedding(self.detector)
        self.index_manager = PgVectorManager()
        self.liveness_detector = LivenessDetector()
        
        # Initialize ScyllaDB and MinIO
        try:
            self.scylladb = ScyllaDBManager()
            logger.info("ScyllaDB manager initialized")
        except Exception as e:
            logger.warning(f"Failed to initialize ScyllaDB: {e}. Continuing without it.")
            self.scylladb = None
        
        try:
            self.minio = MinIOManager()
            logger.info("MinIO manager initialized")
        except Exception as e:
            logger.warning(f"Failed to initialize MinIO: {e}. Continuing without it.")
            self.minio = None
        
        # Initialize database
        self.engine = create_engine(settings.DATABASE_URL)
        Base.metadata.create_all(self.engine)
        self.SessionLocal = sessionmaker(bind=self.engine)
        
        # Load existing embeddings into index
        self._load_embeddings_to_index()
        
        logger.info("FaceService initialized successfully")
    
    def _load_embeddings_to_index(self):
        """Load all active embeddings into FAISS index"""
        try:
            db = self.SessionLocal()
            profiles = db.query(FaceProfile).filter(
                FaceProfile.deleted_at.is_(None)
            ).all()
            
            if profiles:
                embeddings = []
                for profile in profiles:
                    embeddings.append((
                        str(profile.profile_id),
                        str(profile.user_id),
                        self._to_numpy_embedding(profile.embedding),
                        profile.is_primary
                    ))
                
                self.index_manager.rebuild_index(embeddings)
                logger.info(f"Loaded {len(embeddings)} embeddings into index")
            
            db.close()
        except Exception as e:
            logger.error(f"Error loading embeddings: {e}")

    def _to_numpy_embedding(self, embedding_value) -> np.ndarray:
        """Convert stored embedding (list, numpy array, or string) to numpy float32 array.

        The database `Vector` type may return a Python-style string like
        "[0.1, 0.2, ...]". This helper accepts `list`, `tuple`, `np.ndarray`, or
        `str` and returns `np.ndarray(dtype=float32)`.
        """
        if embedding_value is None:
            return None

        # Already a sequence
        if isinstance(embedding_value, (list, tuple, np.ndarray)):
            return np.array(embedding_value, dtype=np.float32)

        # If stored as a string, try JSON then Python literal
        if isinstance(embedding_value, str):
            # Try JSON first (valid when stored as JSON array)
            try:
                parsed = json.loads(embedding_value)
                return np.array(parsed, dtype=np.float32)
            except Exception:
                pass

            # Fallback to ast.literal_eval for Python-style lists
            try:
                parsed = ast.literal_eval(embedding_value)
                return np.array(parsed, dtype=np.float32)
            except Exception as e:
                raise ValueError(f"Unable to parse embedding value: {e}")

        raise ValueError(f"Unsupported embedding type: {type(embedding_value)}")
    
    def _decode_image(self, image_base64: str) -> Optional[np.ndarray]:
        """Decode base64 image to numpy array"""
        try:
            img_data = base64.b64decode(image_base64)
            nparr = np.frombuffer(img_data, np.uint8)
            image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
            return image
        except Exception as e:
            logger.error(f"Error decoding image: {e}")
            return None
    
    def _log_audit(self, db: Session, operation: str, status: str, 
                   company_id: Optional[UUID] = None,
                   profile_id: Optional[UUID] = None,
                   user_id: Optional[UUID] = None,
                   device_id: Optional[str] = None,
                   similarity_score: Optional[float] = None,
                   liveness_score: Optional[float] = None,
                   quality_score: Optional[float] = None,
                   metadata: Optional[Dict] = None,
                   error_message: Optional[str] = None):
        """
        Log audit entry to ScyllaDB (preferred) with PostgreSQL fallback
        
        ScyllaDB provides better performance for time-series audit data
        """
        log_id = uuid4()
        
        # Try ScyllaDB first (preferred for audit logs)
        if self.scylladb:
            try:
                # Convert metadata to string-only dict for ScyllaDB
                scylla_metadata = {}
                if metadata:
                    scylla_metadata = {str(k): str(v) for k, v in metadata.items()}
                
                self.scylladb.save_audit_log(
                    log_id=log_id,
                    company_id=company_id if company_id is not None else UUID(int=0),
                    operation=operation,
                    status=status,
                    profile_id=profile_id,
                    user_id=user_id,
                    device_id=device_id,
                    similarity_score=similarity_score,
                    liveness_score=liveness_score,
                    quality_score=quality_score,
                    metadata=scylla_metadata,
                    error_message=error_message
                )
                logger.debug(f"Audit log saved to ScyllaDB: {log_id}")
                return  # Success - no need for PostgreSQL fallback
            except Exception as e:
                logger.warning(f"Failed to log to ScyllaDB, falling back to PostgreSQL: {e}")
        
        # Fallback to PostgreSQL if ScyllaDB not available or failed
        try:
            audit_log = FaceAuditLog(
                log_id=log_id,
                profile_id=profile_id,
                user_id=user_id,
                operation=operation,
                status=status,
                device_id=device_id,
                similarity_score=similarity_score,
                liveness_score=liveness_score,
                quality_score=quality_score,
                metadata=metadata or {},
                error_message=error_message
            )
            db.add(audit_log)
            db.commit()
            logger.debug(f"Audit log saved to PostgreSQL (fallback): {log_id}")
        except Exception as e:
            logger.error(f"Error logging audit to PostgreSQL: {e}")
            db.rollback()
    
    async def enroll_face(self, user_id: UUID, company_id: UUID, 
                         image_base64: str,
                         device_id: Optional[str] = None,
                         make_primary: bool = False,
                         metadata: Optional[Dict] = None) -> EnrollResponse:
        """
        Enroll a new face for a user
        
        Args:
            user_id: User ID
            company_id: Company ID
            image_base64: Base64 encoded image
            device_id: Device ID
            make_primary: Set as primary profile
            metadata: Additional metadata
            
        Returns:
            EnrollResponse
        """
        db = self.SessionLocal()
        
        try:
            # Ensure partition exists for company
            if not self.index_manager.ensure_company_partition(str(company_id)):
                return EnrollResponse(
                    status="failed", 
                    message=f"Failed to create partition for company {company_id}"
                )
            
            # Decode image
            image = self._decode_image(image_base64)
            if image is None:
                self._log_audit(db, "enroll", "failed", company_id=company_id, user_id=user_id, 
                              error_message="Failed to decode image")
                return EnrollResponse(status="failed", message="Invalid image format")
            
            # Liveness check
            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    self._log_audit(db, "enroll", "failed", company_id=company_id, user_id=user_id,
                                  liveness_score=liveness_score,
                                  error_message="Liveness check failed")
                    return EnrollResponse(
                        status="failed", 
                        message="Liveness check failed. Please ensure you are in front of the camera."
                    )
            
            # Generate embedding with quality
            embedding, quality, face = self.embedding_gen.get_embedding_with_quality(image)
            
            if embedding is None:
                self._log_audit(db, "enroll", "failed", company_id=company_id, user_id=user_id,
                              error_message="No face detected")
                return EnrollResponse(status="failed", message="No face detected in image")
            
            # Quality check
            if quality < 0.3:
                self._log_audit(db, "enroll", "failed", company_id=company_id, user_id=user_id,
                              quality_score=quality,
                              error_message="Low image quality")
                return EnrollResponse(
                    status="failed",
                    message=f"Image quality too low ({quality:.2f}). Please ensure good lighting and focus.",
                    quality_score=quality
                )
            
            # Duplicate detection
            matches = self.index_manager.search(str(company_id), embedding, k=5)
            if matches:
                best_match = matches[0]
                if best_match['similarity'] > settings.DUPLICATE_THRESHOLD:
                    # Check if it's the same user (OK) or different user (duplicate)
                    if best_match['user_id'] != str(user_id):
                        self._log_audit(db, "enroll", "duplicate", company_id=company_id, user_id=user_id,
                                      similarity_score=best_match['similarity'],
                                      metadata={'matched_user': best_match['user_id']})
                        return EnrollResponse(
                            status="duplicate",
                            message="Face already enrolled for another user",
                            duplicate_profiles=[{
                                'user_id': best_match['user_id'],
                                'similarity': best_match['similarity']
                            }]
                        )
            
            # Create face profile
            profile = FaceProfile(
                user_id=user_id,
                company_id=company_id,
                embedding=embedding.tolist(),
                embedding_version=settings.FACE_EMBEDDING_MODEL,
                is_primary=make_primary,
                quality_score=quality,
                meta_data=metadata or {},
                indexed=False
            )
            
            db.add(profile)
            db.flush()
            
            # Save face image to MinIO
            image_path = None
            if self.minio and settings.IMAGE_STORE_ENROLLMENTS:
                try:
                    # Optimize image before storage (resize + compress)
                    optimized_bytes = ImageOptimizer.optimize_for_storage(image)
                    
                    image_path = self.minio.upload_face_image(
                        optimized_bytes,
                        user_id,
                        profile.profile_id
                    )
                    
                    if image_path:
                        profile.enroll_image_path = image_path
                        logger.info(f"Saved optimized face image to MinIO: {image_path}")
                except Exception as e:
                    logger.error(f"Error saving to MinIO: {e}")
            
            # If make_primary, unset other primary profiles
            if make_primary:
                db.query(FaceProfile).filter(
                    and_(
                        FaceProfile.user_id == user_id,
                        FaceProfile.company_id == company_id,
                        FaceProfile.profile_id != profile.profile_id,
                        FaceProfile.deleted_at.is_(None)
                    )
                ).update({'is_primary': False})
            
            # Add to index
            try:
                self.index_manager.add_embedding(
                    str(profile.profile_id),
                    str(company_id),
                    str(user_id),
                    embedding,
                    make_primary
                )
                profile.indexed = True
                self.index_manager.save_index()
            except Exception as e:
                logger.error(f"Error adding to index: {e}")
                db.rollback()
                return EnrollResponse(status="failed", message="Failed to add to index")
            
            db.commit()
            
            # Log success
            self._log_audit(db, "enroll", "success", 
                          company_id=company_id,
                          profile_id=profile.profile_id,
                          user_id=user_id,
                          device_id=device_id,
                          quality_score=quality,
                          metadata=metadata)
            
            # Save enrollment state to ScyllaDB
            if self.scylladb:
                try:
                    # Prepare metadata for ScyllaDB (must be strings)
                    scylla_metadata = {}
                    if metadata:
                        scylla_metadata = {str(k): str(v) for k, v in metadata.items()}
                    
                    self.scylladb.save_enrollment_state(
                        profile_id=profile.profile_id,
                        company_id=company_id,
                        user_id=user_id,
                        device_id=device_id,
                        status="ok",
                        quality_score=quality,
                        metadata=scylla_metadata,
                        image_path=image_path
                    )
                except Exception as e:
                    logger.error(f"Error saving to ScyllaDB: {e}")
            
            return EnrollResponse(
                status="ok",
                profile_id=profile.profile_id,
                message="Face enrolled successfully",
                quality_score=quality
            )
            
        except Exception as e:
            logger.error(f"Error in enroll_face: {e}")
            db.rollback()
            self._log_audit(db, "enroll", "failed", user_id=user_id,
                          error_message=str(e))
            return EnrollResponse(status="failed", message=str(e))
        finally:
            db.close()
    
    async def verify_face(self, image_base64: str,
                         company_id: UUID,
                         user_id: Optional[UUID] = None,
                         device_id: Optional[str] = None,
                         search_mode: str = "1:N",
                         top_k: int = 5) -> VerifyResponse:
        """
        Verify a face
        
        Args:
            image_base64: Base64 encoded image
            user_id: User ID for 1:1 verification
            device_id: Device ID
            search_mode: "1:1" or "1:N"
            top_k: Number of top matches
            
        Returns:
            VerifyResponse
        """
        db = self.SessionLocal()
        
        try:
            # Decode image
            image = self._decode_image(image_base64)
            if image is None:
                return VerifyResponse(
                    status="failed",
                    verified=False,
                    message="Invalid image format"
                )
            
            # Liveness check
            liveness_score = None
            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    self._log_audit(db, "verify", "failed",
                                  company_id=company_id,
                                  liveness_score=liveness_score,
                                  error_message="Liveness check failed")
                    return VerifyResponse(
                        status="failed",
                        verified=False,
                        message="Liveness check failed",
                        liveness_score=liveness_score
                    )
            
            # Generate embedding
            embedding = self.embedding_gen.get_embedding(image)
            
            if embedding is None:
                return VerifyResponse(
                    status="failed",
                    verified=False,
                    message="No face detected"
                )
            
            # Search for matches
            matches = self.index_manager.search(str(company_id), embedding, k=top_k)
            
            if not matches:
                self._log_audit(db, "verify", "no_match", company_id=company_id, device_id=device_id)
                return VerifyResponse(
                    status="no_match",
                    verified=False,
                    message="No matching face found",
                    liveness_score=liveness_score
                )
            
            # Filter by user_id if 1:1 mode
            if search_mode == "1:1" and user_id:
                matches = [m for m in matches if m['user_id'] == str(user_id)]
                if not matches:
                    self._log_audit(db, "verify", "no_match",
                                  company_id=company_id,
                                  user_id=user_id,
                                  device_id=device_id)
                    return VerifyResponse(
                        status="no_match",
                        verified=False,
                        message="Face does not match user",
                        liveness_score=liveness_score
                    )
            
            # Check threshold
            best_match = matches[0]
            verified = best_match['similarity'] >= settings.VERIFY_THRESHOLD
            
            # Format matches
            match_results = []
            for match in matches:
                match_results.append(VerifyMatch(
                    user_id=UUID(match['user_id']),
                    profile_id=UUID(match['profile_id']),
                    similarity=match['similarity'],
                    confidence=match['similarity'],
                    is_primary=match['is_primary']
                ))
            
            # Log result
            status = "match" if verified else "no_match"
            matched_user_id = UUID(best_match['user_id']) if verified else None
            matched_profile_id = UUID(best_match['profile_id']) if verified else None
            
            self._log_audit(db, "verify", status,
                          company_id=company_id,
                          profile_id=matched_profile_id,
                          user_id=matched_user_id,
                          device_id=device_id,
                          similarity_score=best_match['similarity'],
                          liveness_score=liveness_score)
            
            # Save verification image to MinIO (smart storage policy)
            image_path = None
            verification_id = uuid4()
            
            # Determine if we should store this verification image
            should_store = False
            if self.minio:
                if verified and settings.IMAGE_STORE_VERIFICATIONS:
                    # Store successful verifications if enabled
                    should_store = True
                elif not verified and settings.IMAGE_STORE_FAILED_VERIFICATIONS:
                    # Store failed verifications if enabled (useful for debugging)
                    should_store = True
            
            if should_store:
                try:
                    # Optimize image before storage (resize + compress)
                    optimized_bytes = ImageOptimizer.optimize_for_storage(image)
                    
                    image_path = self.minio.upload_verification_image(
                        optimized_bytes,
                        verification_id,
                        matched_user_id
                    )
                    
                    if image_path:
                        logger.info(f"Saved optimized verification image to MinIO: {image_path}")
                except Exception as e:
                    logger.error(f"Error saving verification to MinIO: {e}")
            
            # Always save verification state to ScyllaDB (lightweight)
            if self.scylladb:
                try:
                    self.scylladb.save_verification_state(
                        verification_id=verification_id,
                        company_id=company_id,
                        user_id=matched_user_id,
                        profile_id=matched_profile_id,
                        device_id=device_id,
                        status=status,
                        verified=verified,
                        similarity_score=best_match['similarity'],
                        liveness_score=liveness_score,
                        metadata={},
                        image_path=image_path
                    )
                except Exception as e:
                    logger.error(f"Error saving to ScyllaDB: {e}")
            
            return VerifyResponse(
                status=status,
                verified=verified,
                matches=match_results,
                best_match=match_results[0] if verified else None,
                message="Face verified successfully" if verified else "Face does not match",
                liveness_score=liveness_score
            )
            
        except Exception as e:
            logger.error(f"Error in verify_face: {e}")
            self._log_audit(db, "verify", "failed",
                          company_id=company_id,
                          device_id=device_id,
                          error_message=str(e))
            return VerifyResponse(
                status="failed",
                verified=False,
                message=str(e)
            )
        finally:
            db.close()
    
    async def update_profile(self, profile_id: UUID,
                            company_id: UUID,
                            image_base64: Optional[str] = None,
                            make_primary: Optional[bool] = None,
                            metadata: Optional[Dict] = None) -> Dict:
        """Update face profile"""
        db = self.SessionLocal()
        
        try:
            profile = db.query(FaceProfile).filter(
                FaceProfile.profile_id == profile_id,
                FaceProfile.company_id == company_id,
                FaceProfile.deleted_at.is_(None)
            ).first()
            
            if not profile:
                return {"status": "failed", "message": "Profile not found"}
            
            # Update image and embedding if provided
            if image_base64:
                image = self._decode_image(image_base64)
                if image is None:
                    return {"status": "failed", "message": "Invalid image format"}
                
                embedding, quality, _ = self.embedding_gen.get_embedding_with_quality(image)
                if embedding is None:
                    return {"status": "failed", "message": "No face detected"}
                
                # Update embedding
                profile.embedding = embedding.tolist()
                profile.quality_score = quality
                profile.indexed = False
                
                # Update index
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                self.index_manager.add_embedding(
                    str(profile_id),
                    str(company_id),
                    str(profile.user_id),
                    embedding,
                    profile.is_primary
                )
            
            # Update primary status
            if make_primary is not None and make_primary:
                # Unset other primary profiles
                db.query(FaceProfile).filter(
                    and_(
                        FaceProfile.user_id == profile.user_id,
                        FaceProfile.company_id == company_id,
                        FaceProfile.profile_id != profile_id,
                        FaceProfile.deleted_at.is_(None)
                    )
                ).update({'is_primary': False})
                profile.is_primary = True
            
            # Update metadata
            if metadata:
                profile.meta_data.update(metadata)
            
            profile.updated_at = datetime.utcnow()
            db.commit()
            
            if image_base64:
                self.index_manager.save_index()
            
            self._log_audit(db, "update", "success",
                          company_id=company_id,
                          profile_id=profile_id,
                          user_id=profile.user_id)
            
            return {"status": "ok", "message": "Profile updated successfully"}
            
        except Exception as e:
            logger.error(f"Error updating profile: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e)}
        finally:
            db.close()
    
    async def delete_profile(self, profile_id: UUID, company_id: UUID, hard_delete: bool = False) -> Dict:
        """Delete face profile (soft or hard)"""
        db = self.SessionLocal()
        
        try:
            profile = db.query(FaceProfile).filter(
                FaceProfile.profile_id == profile_id,
                FaceProfile.company_id == company_id
            ).first()
            
            if not profile:
                return {"status": "failed", "message": "Profile not found"}
            
            if hard_delete:
                # Hard delete
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                db.delete(profile)
                message = "Profile permanently deleted"
            else:
                # Soft delete
                profile.deleted_at = datetime.utcnow()
                profile.is_primary = False
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                message = "Profile soft deleted"
            
            db.commit()
            self.index_manager.save_index()
            
            self._log_audit(db, "delete", "success",
                          profile_id=profile_id,
                          user_id=profile.user_id,
                          company_id=company_id,
                          metadata={"hard_delete": hard_delete})
            
            return {"status": "ok", "message": message}
            
        except Exception as e:
            logger.error(f"Error deleting profile: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e)}
        finally:
            db.close()
    
    async def get_user_profiles(self, user_id: UUID, company_id: UUID) -> List[FaceProfileResponse]:
        """Get all face profiles for a user"""
        db = self.SessionLocal()
        
        try:
            profiles = db.query(FaceProfile).filter(
                FaceProfile.user_id == user_id,
                FaceProfile.company_id == company_id,
                FaceProfile.deleted_at.is_(None)
            ).all()
            
            return [
                FaceProfileResponse(
                    profile_id=p.profile_id,
                    user_id=p.user_id,
                    company_id=p.company_id,
                    embedding_version=p.embedding_version,
                    is_primary=p.is_primary,
                    created_at=p.created_at,
                    updated_at=p.updated_at,
                    deleted_at=p.deleted_at,
                    metadata=p.meta_data,
                    quality_score=p.quality_score
                )
                for p in profiles
            ]
        finally:
            db.close()
    
    async def reindex(self, force: bool = False) -> ReindexResponse:
        """Rebuild vector index"""
        start_time = datetime.utcnow()
        db = self.SessionLocal()
        
        try:
            # Check if recent rebuild
            if not force and self.index_manager.last_rebuild:
                time_since_rebuild = datetime.utcnow() - self.index_manager.last_rebuild
                if time_since_rebuild.total_seconds() < settings.VECTOR_INDEX_REBUILD_INTERVAL:
                    return ReindexResponse(
                        status="skipped",
                        message=f"Recent rebuild exists ({time_since_rebuild.total_seconds():.0f}s ago)",
                        profiles_indexed=0,
                        duration_seconds=0
                    )
            
            # Get all active profiles
            profiles = db.query(FaceProfile).filter(
                FaceProfile.deleted_at.is_(None)
            ).all()
            
            # Rebuild index
            embeddings = []
            for profile in profiles:
                embeddings.append((
                    str(profile.profile_id),
                    str(profile.user_id),
                    self._to_numpy_embedding(profile.embedding),
                    profile.is_primary
                ))
            
            self.index_manager.rebuild_index(embeddings)
            
            # Update indexed flag
            db.query(FaceProfile).filter(
                FaceProfile.deleted_at.is_(None)
            ).update({'indexed': True})
            db.commit()
            
            duration = (datetime.utcnow() - start_time).total_seconds()
            
            return ReindexResponse(
                status="ok",
                message="Index rebuilt successfully",
                profiles_indexed=len(profiles),
                duration_seconds=duration
            )
            
        except Exception as e:
            logger.error(f"Error reindexing: {e}")
            db.rollback()
            return ReindexResponse(
                status="failed",
                message=str(e),
                profiles_indexed=0,
                duration_seconds=0
            )
        finally:
            db.close()
    
    async def cleanup_soft_deleted(self):
        """Cleanup soft-deleted profiles past retention period"""
        db = self.SessionLocal()
        
        try:
            cutoff_date = datetime.utcnow() - timedelta(days=settings.SOFT_DELETE_RETENTION)
            
            profiles = db.query(FaceProfile).filter(
                FaceProfile.deleted_at.isnot(None),
                FaceProfile.deleted_at < cutoff_date
            ).all()
            
            for profile in profiles:
                db.delete(profile)
            
            db.commit()
            logger.info(f"Cleaned up {len(profiles)} soft-deleted profiles")
            
        except Exception as e:
            logger.error(f"Error cleaning up soft-deleted profiles: {e}")
            db.rollback()
        finally:
            db.close()
