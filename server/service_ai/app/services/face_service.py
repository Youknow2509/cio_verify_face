import logging
import base64
import numpy as np
import json
import ast
import cv2
import warnings
from typing import Optional, List, Dict
from uuid import UUID, uuid4
from datetime import datetime, time, timedelta
from sqlalchemy import create_engine, and_, text
from sqlalchemy.orm import sessionmaker, Session
# Core modules
from app.core.config import settings
from app.core.face_detector import FaceDetector
from app.core.face_embedding import FaceEmbedding
from app.core.liveness_detector import LivenessDetector
# Database managers
from app.database.milvus_manager import MilvusManager, VectorDBUnavailable

# Silence upstream rcond FutureWarning from insightface transform
warnings.filterwarnings(
    "ignore",
    category=FutureWarning,
    module="insightface.utils.transform"
)
from app.database.scylladb_manager import ScyllaDBManager
from app.database.minio_manager import MinIOManager
# Grpc client
from app.grpc.client.attendance_client import get_client as get_grpc_attendance_client
from app.grpc.client.attendance_client import AttendanceClient as _AttendanceClient
# Models and Schemas
from app.models.attendance_batching_models import RawAttendanceRecord
from app.models.database import Base, FaceProfile
from app.models.schemas import (
    EnrollResponse, VerifyResponse, VerifyMatch,
    FaceProfileResponse, ReindexResponse
)
# Utilities
from app.services.attendance_batching_service import AttendanceBatchingService
from app.utils.image_optimizer import ImageOptimizer
from app.utils.database import get_face_profile_partition_name

logger = logging.getLogger(__name__)


class FaceService:
    def __init__(self, batching_service: Optional[AttendanceBatchingService] = None):
        logger.info("Initializing FaceService ...")
        self.detector = FaceDetector()
        self.embedding_gen = FaceEmbedding(self.detector)
        self.index_manager = MilvusManager()
        self.liveness_detector = LivenessDetector()
        self.batching_service = batching_service
        try:
            self.scylladb = ScyllaDBManager()
        except Exception:
            logger.warning("ScyllaDB init failed; operating in degraded mode.")
            self.scylladb = None

        try:
            self.minio = MinIOManager()
        except Exception:
            logger.warning("MinIO init failed; image persistence disabled.")
            self.minio = None

        try:
            self.attendance_client = get_grpc_attendance_client(
                target=getattr(settings, 'GRPC_ATTENDANCE_URL', None),
                timeout=getattr(settings, 'GRPC_ATTENDANCE_TIMEOUT', 5.0)
            )
        except Exception:
            logger.warning("Attendance gRPC client init failed.")
            self.attendance_client = None

        self.engine = create_engine(settings.DATABASE_URL)
        Base.metadata.create_all(self.engine)
        self.SessionLocal = sessionmaker(bind=self.engine)

        self._load_embeddings_to_index()
        logger.info("FaceService ready.")

    # ---------------- Helpers ----------------
    def _sanitize_metadata(self, metadata: Optional[Dict]) -> Dict:
        """Convert metadata values to JSON-serializable types."""
        if not metadata:
            return {}

        sanitized = {}
        for key, value in metadata.items():
            if isinstance(value, UUID):
                sanitized[key] = str(value)
            elif isinstance(value, datetime):
                sanitized[key] = value.isoformat()
            elif isinstance(value, (list, tuple)):
                sanitized[key] = [str(v) if isinstance(v, (UUID, datetime)) else v for v in value]
            elif isinstance(value, dict):
                sanitized[key] = self._sanitize_metadata(value)
            else:
                sanitized[key] = value
        return sanitized
    
    def _year_month(self, dt: Optional[datetime] = None) -> str:
        """Trả về 'YYYY-MM' từ một datetime."""
        return dt.strftime("%Y-%m") if dt else datetime.utcnow().strftime("%Y-%m")

    def _ensure_company_partition(self, db: Session, company_id: UUID) -> bool:
        partition_name = get_face_profile_partition_name(company_id)
        exists = db.execute(
            text("SELECT relname FROM pg_class WHERE relname = :n"),
            {"n": partition_name}
        ).scalar()
        if exists:
            return True
        try:
            db.execute(text(f"""
                CREATE TABLE IF NOT EXISTS {partition_name}
                PARTITION OF face_profiles FOR VALUES IN ('{company_id}')
            """))
            db.commit()
            logger.info(f"Created PostgreSQL partition {partition_name}")
            return True
        except Exception as e:
            db.rollback()
            logger.error(f"Partition creation failed: {e}")
            return False

    def _load_embeddings_to_index(self):
        try:
            db = self.SessionLocal()
            profiles = db.query(FaceProfile).filter(FaceProfile.deleted_at.is_(None)).all()
            embeddings = []
            for p in profiles:
                embeddings.append((
                    str(p.profile_id),
                    str(p.user_id),
                    self._to_numpy_embedding(p.embedding),
                    p.is_primary
                ))
            if embeddings:
                self.index_manager.rebuild_index(embeddings)
                logger.info(f"Loaded {len(embeddings)} embeddings.")
        except Exception as e:
            logger.error(f"Load embeddings error: {e}")
        finally:
            db.close()

    def _to_numpy_embedding(self, value) -> np.ndarray:
        if value is None:
            return None
        if isinstance(value, (list, tuple, np.ndarray)):
            return np.array(value, dtype=np.float32)
        if isinstance(value, str):
            try:
                return np.array(json.loads(value), dtype=np.float32)
            except Exception:
                try:
                    return np.array(ast.literal_eval(value), dtype=np.float32)
                except Exception as e:
                    raise ValueError(f"Embedding parse failed: {e}")
        raise ValueError(f"Unsupported embedding type: {type(value)}")

    def _decode_image(self, image_base64: str) -> Optional[np.ndarray]:
        try:
            data = base64.b64decode(image_base64)
            arr = np.frombuffer(data, np.uint8)
            return cv2.imdecode(arr, cv2.IMREAD_COLOR)
        except Exception as e:
            logger.error(f"Image decode error: {e}")
            return None

    # ---------------- Logging ----------------
    def _log_audit(
        self,
        *,
        company_id: Optional[UUID],
        actor_id: Optional[UUID],
        action_category: str,
        action_name: str,
        resource_type: str,
        resource_id: Optional[str],
        status: str,
        details: Optional[Dict] = None,
        ip_address: Optional[str] = None,
        user_agent: Optional[str] = None,
        error_message: Optional[str] = None
    ):
        created_at = datetime.utcnow()
        if error_message:
            details = details or {}
            details["error_message"] = error_message
        details = details or {}

        # Scylla preferred
        if self.scylladb and company_id:
            try:
                self.scylladb.add_audit_log(
                    company_id=company_id,
                    actor_id=actor_id,
                    action_category=action_category,
                    action_name=action_name,
                    resource_type=resource_type,
                    resource_id=resource_id or "",
                    details={str(k): str(v) for k, v in details.items()},
                    ip_address=ip_address or "",
                    user_agent=user_agent or "",
                    status=status,
                    created_at=created_at,
                )
                return
            except Exception as e:
                logger.warning(f"Scylla audit failed -> fallback: {e}")
        else:
            logger.warning("ScyllaDB not configured; using fallback audit logging.")

    # ---------------- Enrollment ----------------
    async def enroll_face(
        self,
        user_id: UUID,
        company_id: UUID,
        image_base64: str,
        device_id: Optional[str] = None,
        make_primary: bool = False,
        metadata: Optional[Dict] = None,
        ip_address: Optional[str] = None,
        user_agent: Optional[str] = None
    ) -> EnrollResponse:
        # Get postgres session
        db = self.SessionLocal()
        try:
            if not self._ensure_company_partition(db, company_id):
                return EnrollResponse(status="failed", message="Partition creation failed")

            if not self.index_manager.ensure_company_partition(str(company_id)):
                return EnrollResponse(status="failed", message="Vector partition failed")

            image = self._decode_image(image_base64)
            # Check image decode
            if image is None:
                logger.warning("Enroll failed: invalid image format")
                return EnrollResponse(status="failed", message="Invalid image format")
            # Check image have face is real
            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    logger.warning(f"Enroll failed: liveness check failed (score={liveness_score:.3f})")
                    return EnrollResponse(status="failed", message=f"Liveness check failed with score {liveness_score:.3f}")
            # Get embedding and quality
            embedding, quality, _ = self.embedding_gen.get_embedding_with_quality(image)
            if embedding is None:
                logger.warning("Enroll failed: no face detected in image")
                return EnrollResponse(status="failed", message="No face detected in image")
            # Check quality threshold
            if quality < settings.QUALITY_THRESHOLD:
                logger.warning(f"Enroll failed: image quality too low ({quality:.2f})")
                return EnrollResponse(
                    status="failed",
                    message=f"Image quality too low ({quality:.2f})",
                    quality_score=quality
                )
            # Check duplicate enrollment face for another user
            try:
                matches = self.index_manager.search(str(company_id), embedding, k=5)
            except (VectorDBUnavailable, Exception) as e:
                logger.error(f"Vector search unavailable: {e}")
                return EnrollResponse(status="failed", message="Vector database unavailable, please retry later")
            logger.info(f"Enroll search found {len(matches)} matches for duplicate check.")
            if matches:
                best = matches[0]
                # Check if face is already enrolled to different user
                # DUPLICATE_THRESHOLD = 0.95 means 95% similarity = likely same person
                if best["similarity"] > settings.DUPLICATE_THRESHOLD and best["user_id"] != str(user_id):
                    # Additional check: require sufficient gap from second match to confirm
                    duplicate_confirmed = True
                    if len(matches) > 1:
                        second_best = matches[1]
                        gap = best["similarity"] - second_best["similarity"]
                        # If gap is too small, it's ambiguous - don't block enrollment
                        # E.g., both matches at 0.96 - unclear which is the real match
                        if gap < settings.DUPLICATE_GAP_THRESHOLD:
                            duplicate_confirmed = False
                            logger.info(
                                f"Duplicate check ambiguous: gap too small "
                                f"(best={best['similarity']:.4f}, second={second_best['similarity']:.4f}, gap={gap:.4f}, threshold={settings.DUPLICATE_GAP_THRESHOLD})"
                            )
                    
                    if duplicate_confirmed:
                        logger.warning(
                            f"Enroll failed: duplicate face detected "
                            f"(matched_user={best['user_id']}, similarity={best['similarity']:.4f}, gap_threshold={settings.DUPLICATE_GAP_THRESHOLD})"
                            f" device_id={device_id}, ip_address={ip_address}, user_agent={user_agent}"
                        )
                        return EnrollResponse(
                            status="duplicate",
                            message="Face already enrolled for another user",
                            duplicate_profiles=[{
                                "user_id": best["user_id"],
                                "similarity": best["similarity"]
                            }]
                        )
            # Create new face profile for Milvus only (skip PostgreSQL)
            profile_id = str(uuid4())
            
            # Add embedding to Milvus
            try:
                self.index_manager.add_embedding(
                    profile_id,
                    str(company_id),
                    str(user_id),
                    embedding,
                    make_primary
                )
                self.index_manager.save_index()
                logger.info(f"Added embedding to Milvus for profile {profile_id}")
            except Exception as e:
                logger.error(f"Milvus add embedding failed: {e}")
                return EnrollResponse(status="failed", message="Milvus indexing failed")
            
            # Upload image to MinIO if enabled
            image_path = None
            if self.minio and settings.IMAGE_STORE_ENROLLMENTS:
                try:
                    optimized = ImageOptimizer.optimize_for_storage(image)
                    image_path = self.minio.upload_face_image(optimized, user_id, profile_id)
                except Exception as e:
                    logger.error(f"MinIO enrollment store failed: {e}")
            # Audit log
            try:
                self.scylladb.add_face_enrollment_log(
                    company_id=company_id,
                    employee_id=user_id,
                    action_type="enroll",
                    status="success",
                    image_url=image_path,
                    metadata=metadata or {},
                    created_at=datetime.utcnow()
                )
            except Exception as e:
                logger.error(f"Audit log failed: {e}")
                
            logger.info(f"Face enrollment successful for user {user_id} with profile ID {profile_id}")
            return EnrollResponse(
                status="ok",
                profile_id=profile_id,
                message="Face enrolled successfully",
                quality_score=quality
            )

        except Exception as e:
            logger.error(f"Enroll error: {e}")
            db.rollback()
            return EnrollResponse(status="failed", message=str(e))
        finally:
            db.close()

    # ---------------- Verification ----------------
    async def verify_face(
        self,
        image_base64: str,
        company_id: UUID,
        user_id: Optional[UUID] = None,
        device_id: Optional[str] = None,
        search_mode: str = "1:N",
        top_k: int = 5,
        record_attendance: bool = True,
        location_coordinates: Optional[str] = None,
        ip_address: Optional[str] = None,
        user_agent: Optional[str] = None
    ) -> VerifyResponse:
        db = self.SessionLocal()
        try:
            image = self._decode_image(image_base64)
            if image is None:
                return VerifyResponse(status="failed", verified=False, message="Invalid image format")
            # Liveness check
            liveness_score = None
            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    logger.warning("Verification failed: liveness check failed")
                    return VerifyResponse(
                        status="failed", verified=False,
                        message="Liveness check failed", liveness_score=liveness_score
                    )
            # Get embedding from image
            embedding = self.embedding_gen.get_embedding(image)
            if embedding is None:
                return VerifyResponse(status="failed", verified=False, message="No face detected")
            # Search for matches
            try:
                matches = self.index_manager.search(str(company_id), embedding, k=top_k)
            except (VectorDBUnavailable, Exception) as e:
                logger.error(f"Vector search unavailable: {e}")
                return VerifyResponse(status="failed", verified=False, message="Vector database unavailable, please retry later")
            if not matches:
                logger.warning(
                    f"No matching face found during verification. "
                    f"company_id={company_id}, user_id={user_id}, device_id={device_id}"
                    f"actice_name=verify, 'device_id': {device_id}, 'ip_address': {ip_address}, 'user_agent': {user_agent}"
                    f"'liveness_score': {liveness_score}, 'type_search': {search_mode}"
                )

                return VerifyResponse(
                    status="no_match", verified=False,
                    message="No matching face found", liveness_score=liveness_score
                )
            # If 1:1 mode, filter by user_id
            if search_mode == "1:1" and user_id:
                matches = [m for m in matches if m["user_id"] == str(user_id)]
                if not matches:
                    logger.warning(
                        f"No matching face found during verification in 1:1 mode. "
                        f"company_id={company_id}, user_id={user_id}, device_id={device_id}"
                        f"actice_name=verify, 'device_id': {device_id}, 'ip_address': {ip_address}, 'user_agent': {user_agent}"
                        f"'liveness_score': {liveness_score}, 'type_search': {search_mode}"
                    )

                    return VerifyResponse(
                        status="no_match", verified=False,
                        message="Face does not match user",
                        liveness_score=liveness_score
                    )
            # Get best match and determine verification result
            best = matches[0]
            verified = best["similarity"] >= settings.VERIFY_THRESHOLD
            
            # Additional safety check: require gap between top match and second match
            # This prevents false positives when top 2 matches are too similar
            if verified and len(matches) > 1:
                second_best = matches[1]
                gap = best["similarity"] - second_best["similarity"]
                # Use configurable gap threshold instead of hardcoded value
                if gap < settings.DUPLICATE_GAP_THRESHOLD:
                    logger.warning(
                        f"Verification rejected: similarity gap too small "
                        f"(best={best['similarity']:.4f}, second={second_best['similarity']:.4f}, gap={gap:.4f}, threshold={settings.DUPLICATE_GAP_THRESHOLD})"
                    )
                    verified = False
            
            match_objs = [
                VerifyMatch(
                    user_id=UUID(m["user_id"]),
                    profile_id=UUID(m["profile_id"]),
                    similarity=m["similarity"],
                    confidence=m["similarity"],
                    is_primary=m["is_primary"]
                )
                for m in matches
            ]
            status = "match" if verified else "no_match"
            matched_user_id = UUID(best["user_id"]) if verified else None
            matched_profile_id = UUID(best["profile_id"]) if verified else None
            
            # Log all matches for diagnostics
            matches_str = " | ".join([f"user={m['user_id'][:8]}, sim={m['similarity']:.4f}" for m in matches[:3]])
            logger.info(
                f"Verification result: status={status}, "
                f"company_id={company_id}, user_id={user_id}, device_id={device_id}, "
                f"matched_user_id={matched_user_id}, matched_profile_id={matched_profile_id}, "
                f"best_similarity={best['similarity']:.4f}, threshold={settings.VERIFY_THRESHOLD}, "
                f"top_matches=[{matches_str}], liveness_score={liveness_score}"
            )
            
            # Store verification image if needed
            image_path = None
            verification_id = uuid4()
            should_store = False
            if self.minio:
                if verified and settings.IMAGE_STORE_VERIFICATIONS:
                    should_store = True
                elif not verified and settings.IMAGE_STORE_FAILED_VERIFICATIONS:
                    should_store = True
            if should_store:
                try:
                    optimized = ImageOptimizer.optimize_for_storage(image)
                    image_path = self.minio.upload_verification_image(
                        optimized, verification_id, matched_user_id
                    )
                except Exception as e:
                    logger.error(f"MinIO verify image failed: {e}")
                    
            # Add record to attendance system 
            if verified and record_attendance and self.attendance_client and matched_user_id:
                try:
                    batching = self.batching_service
                    
                    if batching:
                        logger.info(
                            "Recording attendance via gRPC client..."
                            f" company_id={company_id}, employee_id={matched_user_id}, "
                            f"verification_score={best['similarity']:.4f}"
                        )
                        record = RawAttendanceRecord(
                            company_id=company_id,
                            employee_id=matched_user_id,
                            record_time=int(datetime.utcnow().timestamp()),  # Unix timestamp
                            device_id=device_id,
                            verification_method="face",
                            verification_score=best["similarity"],
                            face_image_url=image_path,
                            location_coordinates=location_coordinates
                        )
                        ok = batching.enqueue_record(record)
                        if not ok:
                            logger.error("Failed to enqueue attendance record.")
                    else:
                        logger.error("Attendance batching service not available in app state.")
                    
                except Exception as e:
                    logger.error(f"Attendance record failed: {e}")
            
            return VerifyResponse(
                status=status,
                verified=verified,
                matches=match_objs,
                best_match=match_objs[0] if verified else None,
                message="Face verified successfully" if verified else "Face does not match",
                liveness_score=liveness_score
            )

        except Exception as e:
            logger.error(f"Verify error: {e}")
            return VerifyResponse(status="failed", verified=False, message=str(e))
        finally:
            db.close()

    # ---------------- Update Profile ----------------
    async def update_profile(
        self,
        profile_id: UUID,
        company_id: UUID,
        image_base64: Optional[str] = None,
        make_primary: Optional[bool] = None,
        metadata: Optional[Dict] = None
    ) -> Dict:
        db = self.SessionLocal()
        try:
            profile = db.query(FaceProfile).filter(
                FaceProfile.profile_id == profile_id,
                FaceProfile.company_id == company_id,
                FaceProfile.deleted_at.is_(None)
            ).first()
            if not profile:
                return {"status": "failed", "message": "Profile not found"}

            if image_base64:
                image = self._decode_image(image_base64)
                if image is None:
                    return {"status": "failed", "message": "Invalid image format"}
                embedding, quality, _ = self.embedding_gen.get_embedding_with_quality(image)
                if embedding is None:
                    return {"status": "failed", "message": "No face detected"}
                profile.embedding = embedding.tolist()
                profile.quality_score = quality
                profile.indexed = False
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                self.index_manager.add_embedding(
                    str(profile_id), str(company_id), str(profile.user_id), embedding, profile.is_primary
                )

            if make_primary:
                db.query(FaceProfile).filter(
                    and_(
                        FaceProfile.user_id == profile.user_id,
                        FaceProfile.company_id == company_id,
                        FaceProfile.profile_id != profile_id,
                        FaceProfile.deleted_at.is_(None)
                    )
                ).update({"is_primary": False})
                profile.is_primary = True

            if metadata:
                profile.meta_data.update(metadata)

            profile.updated_at = datetime.utcnow()
            profile.index_version = self.index_manager.index_version
            db.commit()

            if image_base64:
                self.index_manager.save_index()

            self._log_audit(
                company_id=company_id,
                actor_id=uuid4(),
                action_category="face_profile",
                action_name="update_profile",
                resource_type="face_profile",
                resource_id=str(profile_id),
                status="updated",
                details={
                    "make_primary": make_primary,
                    "metadata_updated": bool(metadata),
                },
                ip_address=None,
                user_agent=None
            )
            
            return {"status": "ok", "message": "Profile updated successfully"}
        except Exception as e:
            logger.error(f"Update profile error: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e)}
        finally:
            db.close()

    # ---------------- Delete Profile ----------------
    async def delete_profile(
        self,
        profile_id: UUID,
        company_id: UUID,
        metadata: Optional[Dict] = None,
        hard_delete: bool = False
    ) -> Dict:
        db = self.SessionLocal()
        try:
            profile = db.query(FaceProfile).filter(
                FaceProfile.profile_id == profile_id,
                FaceProfile.company_id == company_id
            ).first()
            if not profile:
                return {"status": "failed", "message": "Profile not found"}

            if hard_delete:
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                db.delete(profile)
                msg = "Profile permanently deleted"
            else:
                profile.deleted_at = datetime.utcnow()
                profile.is_primary = False
                self.index_manager.remove_embedding(str(profile_id), str(company_id))
                msg = "Profile soft deleted"

            db.commit()
            self.index_manager.save_index()
            actor_id = uuid4()
            # Audit log could be added here
            self._log_audit(
                company_id=company_id,
                actor_id=actor_id,
                action_category="face_profile",
                action_name="delete_profile",
                resource_type="face_profile",
                resource_id=str(profile_id),
                status="hard_deleted" if hard_delete else "soft_deleted",
                details={
                    "hard_delete": hard_delete,
                    "session_user": metadata or {}  
                },
                ip_address=None,
                user_agent=None
            )
            return {"status": "ok", "message": msg}
        except Exception as e:
            logger.error(f"Delete profile error: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e)}
        finally:
            db.close()

    # ---------------- Retrieval ----------------
    async def get_user_profiles(self, user_id: UUID, company_id: Optional[UUID]) -> List[FaceProfileResponse]:
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

    # ---------------- Reindex ----------------
    async def reindex(self, force: bool = False) -> ReindexResponse:
        start = datetime.utcnow()
        db = self.SessionLocal()
        try:
            if not force and self.index_manager.last_rebuild:
                delta = datetime.utcnow() - self.index_manager.last_rebuild
                if delta.total_seconds() < settings.VECTOR_INDEX_REBUILD_INTERVAL:
                    return ReindexResponse(
                        status="skipped",
                        message=f"Recent rebuild ({delta.total_seconds():.0f}s ago)",
                        profiles_indexed=0,
                        duration_seconds=0
                    )
            profiles = db.query(FaceProfile).filter(
                FaceProfile.deleted_at.is_(None)
            ).all()
            embeddings = [
                (str(p.profile_id), str(p.user_id), self._to_numpy_embedding(p.embedding), p.is_primary)
                for p in profiles
            ]
            self.index_manager.rebuild_index(embeddings)
            db.query(FaceProfile).filter(FaceProfile.deleted_at.is_(None)).update({
                "indexed": True,
                "index_version": self.index_manager.index_version
            })
            db.commit()
            duration = (datetime.utcnow() - start).total_seconds()
            return ReindexResponse(
                status="ok",
                message="Index rebuilt successfully",
                profiles_indexed=len(profiles),
                duration_seconds=duration
            )
        except Exception as e:
            logger.error(f"Reindex error: {e}")
            db.rollback()
            return ReindexResponse(
                status="failed",
                message=str(e),
                profiles_indexed=0,
                duration_seconds=0
            )
        finally:
            db.close()

    # ---------------- Cleanup Profiles For Company ----------------
    async def cleanup_profiles_for_company(self, company_id: UUID, metadata: Optional[Dict]) -> Dict:
        db = self.SessionLocal()
        try:
            # Get profile IDs first
            profile_ids = db.query(FaceProfile.profile_id).filter(
                FaceProfile.company_id == company_id,
                FaceProfile.deleted_at.isnot(None)
            ).all()
            
            profile_ids = [str(p[0]) for p in profile_ids]
            count = len(profile_ids)
            
            if count == 0:
                return {"status": "ok", "message": "No profiles to cleanup", "profiles_cleaned": 0}
            
            # Batch remove from index
            for pid in profile_ids:
                try:
                    self.index_manager.remove_embedding(pid, str(company_id))
                except Exception as e:
                    logger.warning(f"Index removal failed for {pid}: {e}")
            
            # Batch delete from DB
            deleted_count = db.query(FaceProfile).filter(
                FaceProfile.company_id == company_id,
                FaceProfile.deleted_at.isnot(None)
            ).delete(synchronize_session=False)
            
            db.commit()
            self.index_manager.save_index()
            self._log_audit(
                company_id=company_id,
                actor_id=uuid4(),
                action_category="face_profile",
                action_name="cleanup_profiles_for_company",
                resource_type="face_profile",
                resource_id=None,
                status="cleaned",
                details={
                    "profiles_cleaned": deleted_count,
                    "session_user": metadata or {}
                },
                ip_address=None,
                user_agent=None
            )
            logger.info(f"Cleaned up {deleted_count} profiles for company {company_id}")
            return {"status": "ok", "message": f"Cleaned {deleted_count} profiles", "profiles_cleaned": deleted_count}
            
        except Exception as e:
            logger.error(f"Cleanup error: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e), "profiles_cleaned": 0}
        finally:
            db.close()
    
    # ---------------- Cleanup ----------------
    async def cleanup_soft_deleted(self):
        db = self.SessionLocal()
        try:
            cutoff = datetime.utcnow() - timedelta(days=settings.SOFT_DELETE_RETENTION)
            to_delete = db.query(FaceProfile).filter(
                FaceProfile.deleted_at.isnot(None),
                FaceProfile.deleted_at < cutoff
            ).all()
            for p in to_delete:
                db.delete(p)
            db.commit()
            logger.info(f"Cleaned {len(to_delete)} soft-deleted profiles.")
        except Exception as e:
            logger.error(f"Cleanup error: {e}")
            db.rollback()
        finally:
            db.close()