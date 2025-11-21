import logging
import time
from typing import Optional, Dict, Any
from datetime import datetime
from uuid import UUID
from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider
from cassandra.query import SimpleStatement, ConsistencyLevel, PreparedStatement
from cassandra.policies import RetryPolicy

from app.core.config import settings

logger = logging.getLogger(__name__)


class ScyllaDBManager:
    """Manager for ScyllaDB operations with prepared statements & new schemas."""

    def __init__(self):
        self.cluster = None
        self.session = None
        self._prepared: Dict[str, PreparedStatement] = {}
        self._connect_with_retry()

    # ------------------------------------------------------------------
    # Connection & Setup
    # ------------------------------------------------------------------
    def _connect_with_retry(self, retries: int = 3, delay: float = 2.0):
        """Attempt to connect with limited retries."""
        attempt = 0
        while attempt < retries:
            try:
                self._connect()
                return
            except Exception as e:
                attempt += 1
                logger.error(f"ScyllaDB connection attempt {attempt} failed: {e}")
                if attempt < retries:
                    time.sleep(delay)
        raise RuntimeError("Failed to connect to ScyllaDB after retries")

    def _connect(self):
        hosts = settings.SCYLLADB_HOSTS.split(',')
        logger.info(f"Connecting to ScyllaDB hosts={hosts} port={settings.SCYLLADB_PORT}")

        self.cluster = Cluster(
            hosts,
            port=settings.SCYLLADB_PORT,
            auth_provider=PlainTextAuthProvider(
                username=settings.SCYLLADB_USERNAME,
                password=settings.SCYLLADB_PASSWORD
            ),
            protocol_version=4
        )
        self.session = self.cluster.connect()
        self._create_keyspace()
        self.session.set_keyspace(settings.SCYLLADB_KEYSPACE)
        self._create_tables()
        self._prepare_all()
        logger.info("ScyllaDB connection established")

    def _create_keyspace(self):
        query = f"""
            CREATE KEYSPACE IF NOT EXISTS {settings.SCYLLADB_KEYSPACE}
            WITH replication = {{
                'class': 'SimpleStrategy',
                'replication_factor': {getattr(settings, 'SCYLLADB_REPLICATION_FACTOR', 1)}
            }}
        """
        self.session.execute(query)
        logger.info(f"Keyspace {settings.SCYLLADB_KEYSPACE} verified.")

    def _create_tables(self):
        """
        Tables per provided schema + existing verification/enrollment states.
        """
        # Audit logs (time-series, partitioned by (company_id, year_month))
        audit_logs_query = """
            CREATE TABLE IF NOT EXISTS audit_logs (
                company_id UUID,
                year_month TEXT,
                created_at TIMESTAMP,
                actor_id UUID,
                action_category TEXT,
                action_name TEXT,
                resource_type TEXT,
                resource_id TEXT,
                details MAP<TEXT, TEXT>,
                ip_address TEXT,
                user_agent TEXT,
                status TEXT,
                PRIMARY KEY ((company_id, year_month), created_at, actor_id)
            ) WITH CLUSTERING ORDER BY (created_at DESC);
        """

        # Face enrollment logs
        face_enrollment_logs_query = """
            CREATE TABLE IF NOT EXISTS face_enrollment_logs (
                company_id UUID,
                year_month TEXT,
                created_at TIMESTAMP,
                employee_id UUID,
                action_type TEXT,
                status TEXT,
                image_url TEXT,
                failure_reason TEXT,
                metadata MAP<TEXT, TEXT>,
                PRIMARY KEY ((company_id, year_month), created_at, employee_id)
            ) WITH CLUSTERING ORDER BY (created_at DESC);
        """

        # Verification states (authentication_states)
        auth_states_query = """
            CREATE TABLE IF NOT EXISTS authentication_states (
                verification_id UUID PRIMARY KEY,
                company_id UUID,
                user_id UUID,
                profile_id UUID,
                device_id TEXT,
                status TEXT,
                verified BOOLEAN,
                similarity_score DOUBLE,
                liveness_score DOUBLE,
                timestamp TIMESTAMP,
                metadata MAP<TEXT, TEXT>,
                image_path TEXT
            );
        """

        # Enrollment states (enrollment_states)
        enrollment_states_query = """
            CREATE TABLE IF NOT EXISTS enrollment_states (
                profile_id UUID PRIMARY KEY,
                company_id UUID,
                user_id UUID,
                device_id TEXT,
                status TEXT,
                quality_score DOUBLE,
                timestamp TIMESTAMP,
                metadata MAP<TEXT, TEXT>,
                image_path TEXT
            );
        """

        for q, name in [
            (audit_logs_query, "audit_logs"),
            (face_enrollment_logs_query, "face_enrollment_logs"),
            (auth_states_query, "authentication_states"),
            (enrollment_states_query, "enrollment_states"),
        ]:
            self.session.execute(q)
            logger.info(f"Table {name} verified.")

    # ------------------------------------------------------------------
    # Prepared Statements
    # ------------------------------------------------------------------
    def _prepare_all(self):
        """Prepare statements used frequently."""

        prep_defs = {
            "insert_audit_log": """
                INSERT INTO audit_logs (
                    company_id, year_month, created_at, actor_id,
                    action_category, action_name, resource_type, resource_id,
                    details, ip_address, user_agent, status
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            "insert_face_enrollment_log": """
                INSERT INTO face_enrollment_logs (
                    company_id, year_month, created_at, employee_id,
                    action_type, status, image_url, failure_reason, metadata
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            "insert_auth_state": """
                INSERT INTO authentication_states (
                    verification_id, company_id, user_id, profile_id, device_id,
                    status, verified, similarity_score, liveness_score,
                    timestamp, metadata, image_path
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            "insert_enrollment_state": """
                INSERT INTO enrollment_states (
                    profile_id, company_id, user_id, device_id, status,
                    quality_score, timestamp, metadata, image_path
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            "select_auth_state": """
                SELECT * FROM authentication_states WHERE verification_id = ?
            """,
            "select_user_verifications": """
                SELECT verification_id, user_id, profile_id, device_id,
                       status, verified, similarity_score, liveness_score,
                       timestamp
                FROM authentication_states
                WHERE company_id = ? AND user_id = ?
                LIMIT ?
            """,
            "select_audit_logs_company_month": """
                SELECT * FROM audit_logs
                WHERE company_id = ? AND year_month = ?
                LIMIT ?
            """,
            "select_audit_logs_company_month_category": """
                SELECT * FROM audit_logs
                WHERE company_id = ? AND year_month = ? AND action_category = ?
                LIMIT ?
            """,
            "select_enrollment_logs_company_month": """
                SELECT * FROM face_enrollment_logs
                WHERE company_id = ? AND year_month = ?
                LIMIT ?
            """,
        }

        for key, cql in prep_defs.items():
            try:
                self._prepared[key] = self.session.prepare(cql)
                logger.debug(f"Prepared statement cached: {key}")
            except Exception as e:
                logger.error(f"Failed to prepare {key}: {e}")

    def _execute(self, key: str, params: tuple, cl: ConsistencyLevel = ConsistencyLevel.QUORUM):
        """Execute prepared statement safely."""
        ps = self._prepared.get(key)
        if not ps:
            logger.error(f"Prepared statement missing: {key}")
            return None
        bound = ps.bind(params)
        bound.consistency_level = cl
        return self.session.execute(bound)

    def _coerce_map(self, data: Optional[Dict]) -> Dict[str, str]:
        if not data:
            return {}
        return {str(k): str(v) for k, v in data.items()}

    # ------------------------------------------------------------------
    # Audit & Enrollment Logging
    # ------------------------------------------------------------------
    def store_audit_action(
        self,
        *,
        company_id: Optional[UUID],
        year_month: str,
        created_at: datetime,
        actor_id: Optional[UUID],
        action_category: str,
        action_name: str,
        resource_type: str,
        resource_id: str,
        details: Dict[str, str],
        ip_address: str,
        user_agent: str,
        status: str
    ):
        """
        Store an audit action aligned with provided audit_logs schema.
        """
        if company_id is None:
            logger.warning("Skipping audit log: company_id is None")
            return
        try:
            self._execute(
                "insert_audit_log",
                (
                    company_id,
                    year_month,
                    created_at,
                    actor_id,
                    action_category,
                    action_name,
                    resource_type,
                    resource_id,
                    self._coerce_map(details),
                    ip_address,
                    user_agent,
                    status
                )
            )
            logger.debug(f"Audit stored: {action_category}/{action_name}/{status}")
        except Exception as e:
            logger.error(f"Error storing audit log: {e}")

    def store_face_enrollment_log(
        self,
        *,
        company_id: UUID,
        year_month: str,
        created_at: datetime,
        employee_id: UUID,
        action_type: str,
        status: str,
        image_url: str,
        failure_reason: str,
        metadata: Dict[str, str]
    ):
        """
        Store face enrollment log per face_enrollment_logs schema.
        """
        try:
            self._execute(
                "insert_face_enrollment_log",
                (
                    company_id,
                    year_month,
                    created_at,
                    employee_id,
                    action_type,
                    status,
                    image_url,
                    failure_reason,
                    self._coerce_map(metadata),
                )
            )
            logger.debug(f"Face enrollment log stored: {status}/{failure_reason or 'ok'}")
        except Exception as e:
            logger.error(f"Error storing face enrollment log: {e}")

    # ------------------------------------------------------------------
    # Verification & Enrollment States
    # ------------------------------------------------------------------
    def save_verification_state(
        self,
        *,
        verification_id: UUID,
        company_id: UUID,
        user_id: Optional[UUID],
        profile_id: Optional[UUID],
        device_id: Optional[str],
        status: str,
        verified: bool,
        similarity_score: Optional[float],
        liveness_score: Optional[float],
        metadata: Optional[Dict[str, str]],
        image_path: Optional[str]
    ):
        try:
            self._execute(
                "insert_auth_state",
                (
                    verification_id,
                    company_id,
                    user_id,
                    profile_id,
                    device_id,
                    status,
                    verified,
                    similarity_score,
                    liveness_score,
                    datetime.utcnow(),
                    self._coerce_map(metadata),
                    image_path
                )
            )
            logger.debug(f"Verification state saved: {verification_id} verified={verified}")
        except Exception as e:
            logger.error(f"Error saving verification state: {e}")

    def save_enrollment_state(
        self,
        *,
        profile_id: UUID,
        company_id: UUID,
        user_id: UUID,
        device_id: Optional[str],
        status: str,
        quality_score: Optional[float],
        metadata: Optional[Dict[str, str]],
        image_path: Optional[str]
    ):
        try:
            self._execute(
                "insert_enrollment_state",
                (
                    profile_id,
                    company_id,
                    user_id,
                    device_id,
                    status,
                    quality_score,
                    datetime.utcnow(),
                    self._coerce_map(metadata),
                    image_path
                )
            )
            logger.debug(f"Enrollment state saved: {profile_id} status={status}")
        except Exception as e:
            logger.error(f"Error saving enrollment state: {e}")

    # ------------------------------------------------------------------
    # Queries
    # ------------------------------------------------------------------
    def get_verification_state(self, verification_id: UUID) -> Optional[Dict[str, Any]]:
        try:
            rows = self._execute("select_auth_state", (verification_id,), ConsistencyLevel.ONE)
            row = rows.one() if rows else None
            if not row:
                return None
            return {
                "verification_id": row.verification_id,
                "company_id": row.company_id,
                "user_id": row.user_id,
                "profile_id": row.profile_id,
                "device_id": row.device_id,
                "status": row.status,
                "verified": row.verified,
                "similarity_score": row.similarity_score,
                "liveness_score": row.liveness_score,
                "timestamp": row.timestamp,
                "metadata": row.metadata,
                "image_path": row.image_path
            }
        except Exception as e:
            logger.error(f"Error fetching verification state: {e}")
            return None

    def get_user_verifications(self, company_id: UUID, user_id: UUID, limit: int = 100) -> list:
        try:
            rows = self._execute(
                "select_user_verifications",
                (company_id, user_id, limit),
                ConsistencyLevel.ONE
            )
            if not rows:
                return []
            return [
                {
                    "verification_id": r.verification_id,
                    "user_id": r.user_id,
                    "profile_id": r.profile_id,
                    "device_id": r.device_id,
                    "status": r.status,
                    "verified": r.verified,
                    "similarity_score": r.similarity_score,
                    "liveness_score": r.liveness_score,
                    "timestamp": r.timestamp
                }
                for r in rows
            ]
        except Exception as e:
            logger.error(f"Error fetching user verifications: {e}")
            return []

    def get_audit_logs(
        self,
        company_id: UUID,
        year_month: str,
        action_category: Optional[str] = None,
        limit: int = 100
    ) -> list:
        """
        Fetch audit logs scoped by company + month, optionally filtered by category.
        """
        try:
            if action_category:
                rows = self._execute(
                    "select_audit_logs_company_month_category",
                    (company_id, year_month, action_category, limit),
                    ConsistencyLevel.ONE
                )
            else:
                rows = self._execute(
                    "select_audit_logs_company_month",
                    (company_id, year_month, limit),
                    ConsistencyLevel.ONE
                )
            if not rows:
                return []
            return [
                {
                    "company_id": r.company_id,
                    "year_month": r.year_month,
                    "created_at": r.created_at,
                    "actor_id": r.actor_id,
                    "action_category": r.action_category,
                    "action_name": r.action_name,
                    "resource_type": r.resource_type,
                    "resource_id": r.resource_id,
                    "details": r.details,
                    "ip_address": r.ip_address,
                    "user_agent": r.user_agent,
                    "status": r.status
                }
                for r in rows
            ]
        except Exception as e:
            logger.error(f"Error fetching audit logs: {e}")
            return []

    def get_face_enrollment_logs(
        self,
        company_id: UUID,
        year_month: str,
        limit: int = 100
    ) -> list:
        try:
            rows = self._execute(
                "select_enrollment_logs_company_month",
                (company_id, year_month, limit),
                ConsistencyLevel.ONE
            )
            if not rows:
                return []
            return [
                {
                    "company_id": r.company_id,
                    "year_month": r.year_month,
                    "created_at": r.created_at,
                    "employee_id": r.employee_id,
                    "action_type": r.action_type,
                    "status": r.status,
                    "image_url": r.image_url,
                    "failure_reason": r.failure_reason,
                    "metadata": r.metadata
                }
                for r in rows
            ]
        except Exception as e:
            logger.error(f"Error fetching face enrollment logs: {e}")
            return []

    # ------------------------------------------------------------------
    # Close
    # ------------------------------------------------------------------
    def close(self):
        try:
            if self.cluster:
                self.cluster.shutdown()
                logger.info("ScyllaDB connection closed")
        except Exception as e:
            logger.error(f"Error closing ScyllaDB connection: {e}")