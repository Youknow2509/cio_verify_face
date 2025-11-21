"""
FaceService sử dụng ScyllaDBManager đã tối ưu:
- Ghi audit theo schema mới (audit_logs)
- Ghi log enrollment (face_enrollment_logs)
- Giữ các luồng enroll/verify/update/delete như bản nâng cấp trước
"""

import logging
import base64
import numpy as np
import json
import ast
import cv2
from typing import Optional, List, Dict
from uuid import UUID, uuid4
from datetime import datetime, timedelta
from sqlalchemy import create_engine, and_, text
from sqlalchemy.orm import sessionmaker, Session

from app.core.config import settings
from app.services.face_detector import FaceDetector
from app.services.face_embedding import FaceEmbedding
from app.services.pgvector_manager import PgVectorManager
from app.services.liveness_detector import LivenessDetector
from app.services.scylladb_manager import ScyllaDBManager
from app.services.minio_manager import MinIOManager
from app.utils.image_optimizer import ImageOptimizer
from app.models.database import Base, FaceProfile
from app.models.schemas import (
    EnrollResponse, VerifyResponse, VerifyMatch,
    FaceProfileResponse, ReindexResponse
)
from app.services.attendance_client import get_client as get_attendance_client

logger = logging.getLogger(__name__)


class FaceService:
    def __init__(self):
        logger.info("Initializing FaceService ...")
        self.detector = FaceDetector()
        self.embedding_gen = FaceEmbedding(self.detector)
        self.index_manager = PgVectorManager()
        self.liveness_detector = LivenessDetector()

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
            self.attendance_client = get_attendance_client(
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
    def _year_month(self, dt: Optional[datetime] = None) -> str:
        dt = dt or datetime.utcnow()
        return dt.strftime("%Y%m")

    def _ensure_company_partition(self, db: Session, company_id: UUID) -> bool:
        partition_name = f"face_profiles_p_{str(company_id).replace('-', '')}"
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
        db: Session,
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
        year_month = self._year_month(created_at)
        if error_message:
            details = details or {}
            details["error_message"] = error_message
        details = details or {}

        # Scylla preferred
        if self.scylladb and company_id:
            try:
                self.scylladb.store_audit_action(
                    company_id=company_id,
                    year_month=year_month,
                    created_at=created_at,
                    actor_id=actor_id,
                    action_category=action_category,
                    action_name=action_name,
                    resource_type=resource_type,
                    resource_id=resource_id or "",
                    details={str(k): str(v) for k, v in details.items()},
                    ip_address=ip_address or "",
                    user_agent=user_agent or "",
                    status=status
                )
                return
            except Exception as e:
                logger.warning(f"Scylla audit failed -> fallback: {e}")

    def _log_face_enrollment(
        self,
        *,
        company_id: UUID,
        employee_id: UUID,
        action_type: str,
        status: str,
        image_url: Optional[str],
        failure_reason: Optional[str],
        metadata: Optional[Dict]
    ):
        if not self.scylladb:
            return
        try:
            self.scylladb.store_face_enrollment_log(
                company_id=company_id,
                year_month=self._year_month(),
                created_at=datetime.utcnow(),
                employee_id=employee_id,
                action_type=action_type,
                status=status,
                image_url=image_url or "",
                failure_reason=failure_reason or "",
                metadata={str(k): str(v) for k, v in (metadata or {}).items()}
            )
        except Exception as e:
            logger.warning(f"Scylla enrollment log failed: {e}")

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
        db = self.SessionLocal()
        try:
            if not self._ensure_company_partition(db, company_id):
                return EnrollResponse(status="failed", message="Partition creation failed")

            if not self.index_manager.ensure_company_partition(str(company_id)):
                return EnrollResponse(status="failed", message="Vector partition failed")

            image = self._decode_image(image_base64)
            if image is None:
                self._log_audit(
                    db, company_id=company_id, actor_id=user_id,
                    action_category="face_enrollment", action_name="enroll",
                    resource_type="face_profile", resource_id=None,
                    status="failed",
                    details={"device_id": device_id},
                    ip_address=ip_address, user_agent=user_agent,
                    error_message="invalid_image"
                )
                self._log_face_enrollment(
                    company_id=company_id, employee_id=user_id,
                    action_type="enroll", status="failed",
                    image_url=None, failure_reason="invalid_image",
                    metadata=metadata
                )
                return EnrollResponse(status="failed", message="Invalid image format")

            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    self._log_audit(
                        db, company_id=company_id, actor_id=user_id,
                        action_category="face_enrollment", action_name="enroll",
                        resource_type="face_profile", resource_id=None,
                        status="failed",
                        details={"device_id": device_id, "liveness_score": liveness_score},
                        ip_address=ip_address, user_agent=user_agent,
                        error_message="liveness_failed"
                    )
                    self._log_face_enrollment(
                        company_id=company_id, employee_id=user_id,
                        action_type="enroll", status="failed",
                        image_url=None, failure_reason="liveness_failed",
                        metadata=metadata
                    )
                    return EnrollResponse(status="failed", message="Liveness check failed")

            embedding, quality, _ = self.embedding_gen.get_embedding_with_quality(image)
            if embedding is None:
                self._log_audit(
                    db, company_id=company_id, actor_id=user_id,
                    action_category="face_enrollment", action_name="enroll",
                    resource_type="face_profile", resource_id=None,
                    status="failed",
                    details={"device_id": device_id},
                    ip_address=ip_address, user_agent=user_agent,
                    error_message="no_face"
                )
                self._log_face_enrollment(
                    company_id=company_id, employee_id=user_id,
                    action_type="enroll", status="failed",
                    image_url=None, failure_reason="no_face",
                    metadata=metadata
                )
                return EnrollResponse(status="failed", message="No face detected in image")

            if quality < 0.3:
                self._log_audit(
                    db, company_id=company_id, actor_id=user_id,
                    action_category="face_enrollment", action_name="enroll",
                    resource_type="face_profile", resource_id=None,
                    status="failed",
                    details={"device_id": device_id, "quality_score": quality},
                    ip_address=ip_address, user_agent=user_agent,
                    error_message="low_quality"
                )
                self._log_face_enrollment(
                    company_id=company_id, employee_id=user_id,
                    action_type="enroll", status="failed",
                    image_url=None, failure_reason="low_quality",
                    metadata={"quality": quality}
                )
                return EnrollResponse(
                    status="failed",
                    message=f"Image quality too low ({quality:.2f})",
                    quality_score=quality
                )

            matches = self.index_manager.search(str(company_id), embedding, k=5)
            if matches:
                best = matches[0]
                if best["similarity"] > settings.DUPLICATE_THRESHOLD and best["user_id"] != str(user_id):
                    self._log_audit(
                        db, company_id=company_id, actor_id=user_id,
                        action_category="face_enrollment", action_name="enroll_duplicate",
                        resource_type="face_profile", resource_id=best["profile_id"],
                        status="duplicate",
                        details={
                            "device_id": device_id,
                            "similarity_score": best["similarity"],
                            "matched_user": best["user_id"]
                        },
                        ip_address=ip_address, user_agent=user_agent
                    )
                    self._log_face_enrollment(
                        company_id=company_id, employee_id=user_id,
                        action_type="enroll", status="duplicate",
                        image_url=None, failure_reason="duplicate",
                        metadata={"matched_user": best["user_id"], "similarity": best["similarity"]}
                    )
                    return EnrollResponse(
                        status="duplicate",
                        message="Face already enrolled for another user",
                        duplicate_profiles=[{
                            "user_id": best["user_id"],
                            "similarity": best["similarity"]
                        }]
                    )

            profile = FaceProfile(
                user_id=user_id,
                company_id=company_id,
                embedding=embedding.tolist(),
                embedding_version=settings.FACE_EMBEDDING_MODEL,
                is_primary=make_primary,
                quality_score=quality,
                meta_data=metadata or {},
                indexed=False,
                index_version=self.index_manager.index_version
            )
            db.add(profile)
            db.flush()

            image_path = None
            if self.minio and settings.IMAGE_STORE_ENROLLMENTS:
                try:
                    optimized = ImageOptimizer.optimize_for_storage(image)
                    image_path = self.minio.upload_face_image(optimized, user_id, profile.profile_id)
                    if image_path:
                        profile.enroll_image_path = image_path
                except Exception as e:
                    logger.error(f"MinIO enrollment store failed: {e}")

            if make_primary:
                db.query(FaceProfile).filter(
                    and_(
                        FaceProfile.user_id == user_id,
                        FaceProfile.company_id == company_id,
                        FaceProfile.profile_id != profile.profile_id,
                        FaceProfile.deleted_at.is_(None)
                    )
                ).update({"is_primary": False})

            try:
                self.index_manager.add_embedding(
                    str(profile.profile_id),
                    str(company_id),
                    str(user_id),
                    embedding,
                    make_primary
                )
                profile.indexed = True
                db.commit()
                self.index_manager.save_index()
            except Exception as e:
                logger.error(f"Index add failed: {e}")
                db.rollback()
                return EnrollResponse(status="failed", message="Indexing failed")

            self._log_audit(
                db, company_id=company_id, actor_id=user_id,
                action_category="face_enrollment", action_name="enroll",
                resource_type="face_profile", resource_id=str(profile.profile_id),
                status="success",
                details={"device_id": device_id, "quality_score": quality, "image_path": image_path or ""},
                ip_address=ip_address, user_agent=user_agent
            )
            self._log_face_enrollment(
                company_id=company_id, employee_id=user_id,
                action_type="enroll", status="success",
                image_url=image_path, failure_reason=None,
                metadata={"quality": quality}
            )

            if self.scylladb:
                try:
                    self.scylladb.save_enrollment_state(
                        profile_id=profile.profile_id,
                        company_id=company_id,
                        user_id=user_id,
                        device_id=device_id,
                        status="ok",
                        quality_score=quality,
                        metadata=metadata,
                        image_path=image_path
                    )
                except Exception as e:
                    logger.warning(f"save_enrollment_state failed: {e}")

            return EnrollResponse(
                status="ok",
                profile_id=profile.profile_id,
                message="Face enrolled successfully",
                quality_score=quality
            )

        except Exception as e:
            logger.error(f"Enroll error: {e}")
            db.rollback()
            self._log_audit(
                db, company_id=company_id, actor_id=user_id,
                action_category="face_enrollment", action_name="enroll",
                resource_type="face_profile", resource_id=None,
                status="failed",
                details={"device_id": device_id},
                error_message=str(e)
            )
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
        record_attendance: bool = False,
        location_coordinates: Optional[str] = None,
        ip_address: Optional[str] = None,
        user_agent: Optional[str] = None
    ) -> VerifyResponse:
        db = self.SessionLocal()
        try:
            image = self._decode_image(image_base64)
            if image is None:
                return VerifyResponse(status="failed", verified=False, message="Invalid image format")

            liveness_score = None
            if settings.LIVENESS_ENABLED:
                is_live, liveness_score = self.liveness_detector.detect_liveness(image)
                if not is_live:
                    self._log_audit(
                        db, company_id=company_id, actor_id=user_id,
                        action_category="face_verification", action_name="verify",
                        resource_type="face_image", resource_id=None,
                        status="failed",
                        details={"device_id": device_id, "liveness_score": liveness_score},
                        ip_address=ip_address, user_agent=user_agent,
                        error_message="liveness_failed"
                    )
                    return VerifyResponse(
                        status="failed", verified=False,
                        message="Liveness check failed", liveness_score=liveness_score
                    )

            embedding = self.embedding_gen.get_embedding(image)
            if embedding is None:
                return VerifyResponse(status="failed", verified=False, message="No face detected")

            matches = self.index_manager.search(str(company_id), embedding, k=top_k)
            if not matches:
                self._log_audit(
                    db, company_id=company_id, actor_id=user_id,
                    action_category="face_verification", action_name="verify",
                    resource_type="face_image", resource_id=None,
                    status="no_match",
                    details={"device_id": device_id},
                    ip_address=ip_address, user_agent=user_agent
                )
                return VerifyResponse(
                    status="no_match", verified=False,
                    message="No matching face found", liveness_score=liveness_score
                )

            if search_mode == "1:1" and user_id:
                matches = [m for m in matches if m["user_id"] == str(user_id)]
                if not matches:
                    self._log_audit(
                        db, company_id=company_id, actor_id=user_id,
                        action_category="face_verification", action_name="verify",
                        resource_type="face_image", resource_id=None,
                        status="no_match",
                        details={"device_id": device_id},
                        ip_address=ip_address, user_agent=user_agent
                    )
                    return VerifyResponse(
                        status="no_match", verified=False,
                        message="Face does not match user",
                        liveness_score=liveness_score
                    )

            best = matches[0]
            verified = best["similarity"] >= settings.VERIFY_THRESHOLD

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

            self._log_audit(
                db, company_id=company_id, actor_id=matched_user_id,
                action_category="face_verification", action_name="verify",
                resource_type="face_profile" if verified else "face_image",
                resource_id=str(matched_profile_id) if matched_profile_id else None,
                status=status,
                details={
                    "device_id": device_id,
                    "similarity_score": best["similarity"],
                    "liveness_score": liveness_score if liveness_score is not None else ""
                },
                ip_address=ip_address, user_agent=user_agent
            )

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

            logger.info(f"Verification result: verified={verified}, user_id={matched_user_id}, similarity={best['similarity']:.4f}, attendance_client={self.attendance_client is not None}")  # TODO: remove it
            if verified and record_attendance and self.attendance_client and matched_user_id:
                logger.info("Recording attendance via gRPC...")
                try:
                    from app.grpc_generated import attendance_pb2
                    session_info = attendance_pb2.ServiceSessionInfo(
                        service_name=settings.SERVICE_NAME,
                        service_id=settings.SERVICE_ID,
                        client_ip=ip_address or "",
                        client_agent=user_agent or ""
                    )
                    # TODO: handler worker batch send attendance records
                    batch_req_input = attendance_pb2.ServiceAddBatchAttendanceInput(
                        company_id=str(company_id),
                        employee_id=str(matched_user_id),
                        device_id=device_id or "",
                        record_time=int(datetime.utcnow().timestamp()),
                        verification_method="face",
                        verification_score=float(best["similarity"]),
                        face_image_url=image_path or "",
                        location_coordinates=location_coordinates or "",
                        session=session_info
                    )
                    self.attendance_client.service_add_batch_attendance(
                        [batch_req_input]
                    )
                except Exception as e:
                    logger.warning(f"Attendance rpc failed: {e}")
            else:
                logger.info("Attendance gRPC client not configured; skipping attendance record.")
            #  TODO: clean it
            # if self.scylladb:
            #     try:
            #         self.scylladb.save_verification_state(
            #             verification_id=verification_id,
            #             company_id=company_id,
            #             user_id=matched_user_id,
            #             profile_id=matched_profile_id,
            #             device_id=device_id,
            #             status=status,
            #             verified=verified,
            #             similarity_score=best["similarity"],
            #             liveness_score=liveness_score,
            #             metadata={},
            #             image_path=image_path
            #         )
            #     except Exception as e:
            #         logger.warning(f"save_verification_state failed: {e}")

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
            self._log_audit(
                db, company_id=company_id, actor_id=user_id,
                action_category="face_verification", action_name="verify",
                resource_type="face_image", resource_id=None,
                status="failed",
                details={"device_id": device_id},
                error_message=str(e)
            )
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
                db, company_id=company_id, actor_id=profile.user_id,
                action_category="face_profile", action_name="update",
                resource_type="face_profile", resource_id=str(profile_id),
                status="success",
                details={"is_primary": profile.is_primary, "quality_score": profile.quality_score}
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

            self._log_audit(
                db, company_id=company_id, actor_id=profile.user_id,
                action_category="face_profile", action_name="delete",
                resource_type="face_profile", resource_id=str(profile_id),
                status="success",
                details={"hard_delete": hard_delete}
            )
            return {"status": "ok", "message": msg}
        except Exception as e:
            logger.error(f"Delete profile error: {e}")
            db.rollback()
            return {"status": "failed", "message": str(e)}
        finally:
            db.close()

    # ---------------- Retrieval ----------------
    async def get_user_profiles(self, user_id: UUID, company_id: UUID) -> List[FaceProfileResponse]:
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