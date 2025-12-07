"""
Milvus manager for efficient face embedding search (migrated from pgvector).
"""
import logging
from datetime import datetime
from typing import Dict, List, Optional
from uuid import UUID

import numpy as np
from pymilvus import (
    Collection,
    CollectionSchema,
    DataType,
    FieldSchema,
    connections,
    utility,
)
from pymilvus.exceptions import MilvusException
from sqlalchemy import create_engine, text
from sqlalchemy.exc import OperationalError
from sqlalchemy.orm import Session, sessionmaker
from sqlalchemy.pool import QueuePool

from app.core.config import settings
from app.models.schemas import FaceProfileResponse
from app.utils.database import get_face_profile_partition_name


class VectorDBUnavailable(Exception):
    """Raised when the vector database is unavailable (e.g., restarting/recovery)."""


logger = logging.getLogger(__name__)


class MilvusManager:
    """Manage face embeddings using Milvus vector database (with Postgres metadata)."""

    def __init__(self):
        self.dimension = settings.EMBEDDING_DIMENSION

        # Postgres (metadata + backward compatibility)
        self.engine = create_engine(
            settings.DATABASE_URL,
            poolclass=QueuePool,
            pool_size=20,
            max_overflow=30,
            pool_pre_ping=True,
            pool_recycle=3600,
            echo=False,
        )
        self.SessionLocal = sessionmaker(bind=self.engine)

        # Milvus config
        self.milvus_alias = "default"
        self.collection_name = settings.MILVUS_COLLECTION
        self.milvus_db = settings.MILVUS_DB
        self.metric_type = settings.MILVUS_METRIC_TYPE
        self.index_type = settings.MILVUS_INDEX_TYPE
        self.nlist = settings.MILVUS_NLIST
        self.nprobe = settings.MILVUS_NPROBE
        self.index_version = settings.VECTOR_DB_INDEX_VERSION

        self._connect_milvus()
        self._ensure_collection()
        logger.info("MilvusManager initialized successfully")

    # ------------------------ Internal helpers ------------------------
    def _connect_milvus(self):
        try:
            if settings.MILVUS_URI:
                connections.connect(
                    alias=self.milvus_alias,
                    uri=settings.MILVUS_URI,
                    token=settings.MILVUS_TOKEN,
                )
            else:
                connections.connect(
                    alias=self.milvus_alias,
                    host=settings.MILVUS_HOST,
                    port=str(settings.MILVUS_PORT),
                    user=settings.MILVUS_USERNAME,
                    password=settings.MILVUS_PASSWORD,
                    secure=settings.MILVUS_SECURE,
                    db_name=self.milvus_db,
                )
        except Exception as e:
            logger.error(f"Milvus connection failed: {e}")
            raise

    def _get_collection(self) -> Collection:
        return Collection(self.collection_name, using=self.milvus_alias)

    def _ensure_collection(self):
        try:
            if not utility.has_collection(self.collection_name, using=self.milvus_alias):
                fields = [
                    FieldSchema(name="profile_id", dtype=DataType.VARCHAR, is_primary=True, max_length=64),
                    FieldSchema(name="user_id", dtype=DataType.VARCHAR, max_length=64),
                    FieldSchema(name="company_id", dtype=DataType.VARCHAR, max_length=64),
                    FieldSchema(name="is_primary", dtype=DataType.BOOL),
                    FieldSchema(name="indexed", dtype=DataType.BOOL),
                    FieldSchema(name="updated_at", dtype=DataType.FLOAT),
                    FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=self.dimension),
                ]
                schema = CollectionSchema(fields=fields, description="Face embeddings")
                Collection(name=self.collection_name, schema=schema, using=self.milvus_alias)

            coll = self._get_collection()
            self._create_index_if_needed(coll)
            try:
                coll.load()
            except MilvusException as e:
                if "index not found" in str(e).lower():
                    # Create index then retry load
                    self._create_index_if_needed(coll)
                    coll.load()
                else:
                    raise
        except Exception as e:
            logger.error(f"Error ensuring Milvus collection: {e}")
            raise

    def _create_index_if_needed(self, coll: Collection):
        has_embedding_index = any(idx.field_name == "embedding" for idx in coll.indexes)
        if not has_embedding_index:
            params = {"index_type": self.index_type, "metric_type": self.metric_type, "params": {"nlist": self.nlist}}
            try:
                coll.create_index(field_name="embedding", index_params=params)
            except Exception as e:
                logger.error(f"Create index failed: {e}")
                raise

    def _ensure_partition(self, coll: Collection, company_id: str):
        if not coll.has_partition(company_id):
            coll.create_partition(company_id)
        return coll

    def _normalize_embedding(self, embedding: np.ndarray) -> np.ndarray:
        arr = np.asarray(embedding, dtype=np.float32).flatten()
        if arr.size != self.dimension:
            raise ValueError(f"Embedding dimension mismatch: expected {self.dimension}, got {arr.size}")
        if not np.all(np.isfinite(arr)):
            raise ValueError("Embedding contains non-finite values")
        norm = np.linalg.norm(arr)
        if norm == 0:
            raise ValueError("Embedding norm is zero; cannot normalize")
        return arr / norm

    def _lookup_company_id(self, profile_id: str, db: Session) -> Optional[str]:
        sql_raw = text(
            """
            SELECT company_id FROM face_profiles
            WHERE profile_id = CAST(:profile_id AS uuid)
            LIMIT 1
            """
        )
        row = db.execute(sql_raw, {"profile_id": profile_id}).fetchone()
        return str(row[0]) if row else None

    # ------------------------ Connectivity ------------------------
    def check_connection(self) -> bool:
        ok_pg = True
        try:
            with self.engine.connect() as conn:
                conn.execute(text("SELECT 1"))
        except Exception as e:
            ok_pg = False
            logger.error(f"PostgreSQL connection failed: {e}")

        try:
            utility.get_server_version(using=self.milvus_alias)
            return ok_pg and True
        except Exception as e:
            logger.error(f"Milvus connection failed: {e}")
            return False

    # ------------------------ Postgres metadata helpers ------------------------
    def check_employee_exist_in_company(self, company_id: UUID, employee_id: UUID) -> bool:
        with self.SessionLocal() as session:
            sql_raw = text(
                """
                SELECT employee_id FROM employees
                WHERE company_id = :company_id AND employee_id = :employee_id
                """
            )
            result = session.execute(
                sql_raw,
                {"company_id": str(company_id), "employee_id": str(employee_id)},
            ).fetchone()
            return result is not None

    def get_profile_face_employee(
        self,
        company_id: UUID,
        employee_id: UUID,
        page_size: int = 20,
        page_number: int = 1,
    ) -> List[FaceProfileResponse]:
        limit = page_size
        offset = (page_number - 1) * page_size
        partition_table_name = get_face_profile_partition_name(company_id)

        with self.SessionLocal() as session:
            sql_raw = text(
                f"""
                SELECT
                    profile_id, user_id, company_id, embedding_version,
                    is_primary, created_at, updated_at, deleted_at,
                    meta_data, quality_score
                FROM {partition_table_name}
                WHERE company_id = :company_id AND user_id = :employee_id
                LIMIT :limit OFFSET :offset
                """
            )
            result = session.execute(
                sql_raw,
                {
                    "company_id": str(company_id),
                    "employee_id": str(employee_id),
                    "limit": limit,
                    "offset": offset,
                },
            ).fetchall()

            face_profiles = []
            for row in result:
                face_profiles.append(
                    FaceProfileResponse(
                        profile_id=row[0],
                        user_id=row[1],
                        company_id=row[2],
                        embedding_version=row[3],
                        is_primary=row[4],
                        created_at=row[5],
                        updated_at=row[6],
                        deleted_at=row[7],
                        metadata=row[8],
                        quality_score=row[9],
                    )
                )
            return face_profiles

    # ------------------------ Milvus operations ------------------------
    def ensure_company_partition(self, company_id: str) -> bool:
        try:
            coll = self._get_collection()
            self._ensure_partition(coll, company_id)
            return True
        except Exception as e:
            logger.error(f"Error ensuring Milvus partition for company {company_id}: {e}")
            return False

    def add_embedding(self, profile_id: str, company_id: str, user_id: str, embedding: np.ndarray, is_primary: bool = False):
        try:
            # Normalize embedding for IP (Inner Product) metric
            # IP metric expects normalized vectors (L2 norm = 1)
            arr = np.asarray(embedding, dtype=np.float32).flatten()
            if arr.size != self.dimension:
                raise ValueError(f"Embedding dimension mismatch: expected {self.dimension}, got {arr.size}")
            if not np.all(np.isfinite(arr)):
                raise ValueError("Embedding contains non-finite values")
            
            # Normalize: divide by L2 norm to get unit vector
            norm = np.linalg.norm(arr)
            if norm == 0:
                raise ValueError("Embedding norm is zero; cannot normalize")
            arr_normalized = arr / norm
            
            embedding_list = arr_normalized.tolist()

            coll = self._get_collection()
            self._ensure_partition(coll, company_id)
            
            # Delete old embedding if exists
            try:
                coll.delete(expr=f'profile_id == "{profile_id}"')
            except Exception:
                pass
            
            # Insert new embedding into Milvus (unnormalized)
            coll.insert(
                [
                    [profile_id],
                    [user_id],
                    [company_id],
                    [is_primary],
                    [True],
                    [datetime.utcnow().timestamp()],
                    [embedding_list],
                ],
                partition_name=company_id,
            )
            coll.flush()
            
            # Reload collection to ensure new data is searchable
            try:
                coll.release()
            except Exception:
                pass
            coll.load()
            
            logger.info(f"Added embedding to Milvus: profile_id={profile_id}, company_id={company_id}")
        except Exception as e:
            logger.error(f"Error adding embedding to Milvus: {e}")
            raise

    def remove_embedding(self, profile_id: str, company_id: str):
        db = self.SessionLocal()
        try:
            query = text(
                """
                UPDATE face_profiles
                SET indexed = false,
                    updated_at = CURRENT_TIMESTAMP
                WHERE profile_id = CAST(:profile_id AS uuid)
                  AND company_id = CAST(:company_id AS uuid)
                """
            )
            db.execute(query, {"profile_id": profile_id, "company_id": company_id})
            db.commit()

            coll = self._get_collection()
            expr = f'profile_id == "{profile_id}" and company_id == "{company_id}"'
            coll.delete(expr=expr)
            coll.flush()
            logger.info(f"Removed embedding for profile {profile_id} / company {company_id} from Milvus")
        except Exception as e:
            logger.error(f"Error removing embedding: {e}")
            db.rollback()
            raise
        finally:
            db.close()

    def search(self, company_id: str, embedding: np.ndarray, k: int = 5) -> List[Dict]:
        try:
            # Normalize embedding for IP (Inner Product) metric
            # IP metric expects normalized vectors (L2 norm = 1)
            arr = np.asarray(embedding, dtype=np.float32).flatten()
            if arr.size != self.dimension:
                raise ValueError(f"Embedding dimension mismatch: expected {self.dimension}, got {arr.size}")
            if not np.all(np.isfinite(arr)):
                raise ValueError("Embedding contains non-finite values")
            
            # Normalize: divide by L2 norm to get unit vector
            norm = np.linalg.norm(arr)
            if norm == 0:
                raise ValueError("Embedding norm is zero; cannot normalize")
            arr_normalized = arr / norm
            
            embedding_list = arr_normalized.tolist()

            coll = self._get_collection()
            
            # Ensure collection is loaded before searching
            try:
                coll.load()
            except Exception as e:
                logger.warning(f"Load collection warning (may already be loaded): {e}")

            search_params = {
                "metric_type": self.metric_type,
                "params": {"nprobe": self.nprobe},
            }
            expr = f'company_id == "{company_id}" and indexed == true'

            res = coll.search(
                data=[embedding_list],
                anns_field="embedding",
                param=search_params,
                limit=k,
                expr=expr,
                output_fields=["profile_id", "user_id", "is_primary"],
            )
            
            matches: List[Dict] = []
            for hit in res[0]:
                distance = float(hit.distance)
                # For IP metric: distance = dot product of normalized vectors (range [0, 1])
                # IP metric with normalized vectors: 1.0 = identical, 0.0 = orthogonal
                # Simply use distance as similarity directly
                similarity = distance
                matches.append(
                    {
                        "profile_id": hit.entity.get("profile_id"),
                        "user_id": hit.entity.get("user_id"),
                        "is_primary": bool(hit.entity.get("is_primary")),
                        "similarity": similarity,
                        "distance": distance,  # Add raw distance for diagnostics
                    }
                )
            
            if matches:
                logger.debug(
                    f"Search for company {company_id}: returned {len(matches)} matches. "
                    f"Top 3: {', '.join([f'sim={m['similarity']:.4f}' for m in matches[:3]])}"
                )
            logger.info(f"Search returned {len(matches)} matches for company {company_id}")
            return matches
        except Exception as e:
            logger.error(f"Error searching embeddings: {e}")
            raise VectorDBUnavailable("Vector database unavailable or in recovery; please retry later") from e

    def rebuild_index(self, embeddings: List[tuple]):
        db: Optional[Session] = None
        try:
            coll = self._get_collection()
            for idx in coll.indexes:
                try:
                    coll.drop_index(index_name=idx.index_name)
                except Exception:
                    pass
            coll.delete(expr="profile_id != ''")

            db = self.SessionLocal()
            inserted = 0
            for item in embeddings:
                if len(item) == 5:
                    profile_id, user_id, embedding, is_primary, company_id = item
                else:
                    profile_id, user_id, embedding, is_primary = item
                    company_id = self._lookup_company_id(profile_id, db)
                    if company_id is None:
                        continue

                # Normalize embedding for IP metric
                arr = np.asarray(embedding, dtype=np.float32).flatten()
                norm = np.linalg.norm(arr)
                if norm == 0:
                    logger.warning(f"Skipping embedding with zero norm for profile {profile_id}")
                    continue
                arr_normalized = arr / norm
                emb_list = arr_normalized.tolist()
                self._ensure_partition(coll, company_id)
                coll.insert(
                    [
                        [profile_id],
                        [user_id],
                        [company_id],
                        [is_primary],
                        [True],
                        [datetime.utcnow().timestamp()],
                        [emb_list],
                    ],
                    partition_name=company_id,
                )
                inserted += 1

            coll.flush()
            self._create_index_if_needed(coll)
            logger.info(f"Rebuilt Milvus index with {inserted} embeddings")
        except Exception as e:
            logger.error(f"Error rebuilding Milvus index: {e}")
            raise
        finally:
            try:
                if db:
                    db.close()
            except Exception:
                pass

    def save_index(self):
        try:
            coll = self._get_collection()
            coll.flush()
        except Exception as e:
            logger.error(f"Error flushing Milvus index: {e}")

    def get_size(self) -> int:
        try:
            coll = self._get_collection()
            return int(utility.get_collection_stats(coll.name, using=self.milvus_alias)["row_count"])
        except Exception as e:
            logger.error(f"Error getting Milvus size: {e}")
            return 0

    def clear(self):
        try:
            coll = self._get_collection()
            coll.delete(expr="profile_id != ''")
            coll.flush()
            logger.info("Milvus index cleared")
        except Exception as e:
            logger.error(f"Error clearing Milvus index: {e}")
            raise

    @property
    def last_rebuild(self):
        try:
            coll = self._get_collection()
            stats = utility.get_collection_stats(coll.name, using=self.milvus_alias)
            return datetime.utcnow() if not stats else datetime.utcnow()
        except Exception:
            return datetime.utcnow()


# Backward-compatible alias used across services
PgVectorManager = MilvusManager
