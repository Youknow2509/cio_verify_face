import logging
import time
import base64
from typing import Optional, Dict, Any, List, Tuple
from datetime import datetime
from uuid import UUID

from cassandra.cluster import Cluster, ResultSet
from cassandra.auth import PlainTextAuthProvider
from cassandra.query import ConsistencyLevel, PreparedStatement
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

        for q, name in [
            (audit_logs_query, "audit_logs"),
            (face_enrollment_logs_query, "face_enrollment_logs"),
        ]:
            self.session.execute(q)
            logger.info(f"Table {name} verified.")

    # ------------------------------------------------------------------
    # Prepared Statements
    # ------------------------------------------------------------------
    def _prepare_all(self):
        """Prepare statements used frequently."""

        prep_defs = {
            # Insert
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
            # Select (face enrollment)
            "get_month_face_enrollment_logs_company": """
                SELECT * FROM face_enrollment_logs
                WHERE company_id = ? AND year_month = ?
            """,
            "get_month_face_enrollment_logs_with_around_time": """
                SELECT *
                FROM face_enrollment_logs
                WHERE company_id = ?
                  AND year_month = ?
                  AND created_at >= ?
                  AND created_at <= ?
            """,
            # Select (audit logs)
            "get_month_audit_logs_company": """
                SELECT * FROM audit_logs
                WHERE company_id = ? AND year_month = ?
            """,
            "get_month_audit_logs_with_around_time": """
                SELECT *
                FROM audit_logs
                WHERE company_id = ?
                  AND year_month = ?
                  AND created_at >= ?
                  AND created_at <= ?
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

    @staticmethod
    def _year_month(dt: datetime) -> str:
        """Trả về 'YYYY-MM' từ một datetime."""
        return dt.strftime("%Y-%m")

    @staticmethod
    def _encode_paging_state(paging_state: Optional[bytes]) -> Optional[str]:
        if not paging_state:
            return None
        return base64.b64encode(paging_state).decode("utf-8")

    @staticmethod
    def _decode_paging_state(paging_state_str: Optional[str]) -> Optional[bytes]:
        if not paging_state_str:
            return None
        try:
            return base64.b64decode(paging_state_str.encode("utf-8"))
        except Exception:
            logger.warning("Invalid paging_state provided; ignoring.")
            return None

    @staticmethod
    def _row_to_dict(row) -> Dict[str, Any]:
        if row is None:
            return {}
        return dict(row._asdict())

    # ------------------------------------------------------------------
    # Insert Operations
    # ------------------------------------------------------------------
    def add_audit_log(
        self,
        company_id: UUID,
        actor_id: UUID,
        action_category: str,
        action_name: str,
        resource_type: str,
        resource_id: str,
        details: Optional[Dict[str, Any]] = None,
        ip_address: Optional[str] = None,
        user_agent: Optional[str] = None,
        status: Optional[str] = None,
        created_at: Optional[datetime] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.QUORUM
    ) -> bool:
        """Thêm một bản ghi audit log."""
        created_at = created_at or datetime.utcnow()
        year_month = self._year_month(created_at)
        details_map = self._coerce_map(details)
        params = (
            company_id, year_month, created_at, actor_id,
            action_category, action_name, resource_type, resource_id,
            details_map, ip_address, user_agent, status
        )
        try:
            self._execute("insert_audit_log", params, cl=consistency)
            return True
        except Exception as e:
            logger.error(f"add_audit_log failed: {e}")
            return False

    def add_face_enrollment_log(
        self,
        company_id: UUID,
        employee_id: UUID,
        action_type: str,
        status: str,
        image_url: Optional[str] = None,
        failure_reason: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
        created_at: Optional[datetime] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.QUORUM
    ) -> bool:
        """Thêm một bản ghi face enrollment log."""
        created_at = created_at or datetime.utcnow()
        year_month = self._year_month(created_at)
        metadata_map = self._coerce_map(metadata)
        params = (
            company_id, year_month, created_at, employee_id,
            action_type, status, image_url, failure_reason, metadata_map
        )
        try:
            self._execute("insert_face_enrollment_log", params, cl=consistency)
            return True
        except Exception as e:
            logger.error(f"add_face_enrollment_log failed: {e}")
            return False

    # ------------------------------------------------------------------
    # Generic Pagination Helper
    # ------------------------------------------------------------------
    def _paged_query(
        self,
        key: str,
        params: tuple,
        page_size: int,
        paging_state: Optional[str],
        consistency: ConsistencyLevel = ConsistencyLevel.QUORUM
    ) -> Tuple[List[Dict[str, Any]], Optional[str]]:
        """
        Thực hiện truy vấn phân trang cho prepared statement đã định nghĩa.
        Trả về (list_rows, next_paging_state_str).
        """
        ps = self._prepared.get(key)
        if not ps:
            raise ValueError(f"Prepared statement not found: {key}")

        bound = ps.bind(params)
        bound.consistency_level = consistency
        bound.fetch_size = page_size

        decoded_state = self._decode_paging_state(paging_state)
        try:
            rs: ResultSet = self.session.execute(bound, paging_state=decoded_state)
        except Exception as e:
            logger.error(f"Paged query execution failed ({key}): {e}")
            return [], None

        rows = [self._row_to_dict(r) for r in rs.current_rows]
        next_state = self._encode_paging_state(rs.paging_state)
        return rows, next_state

    # ------------------------------------------------------------------
    # Select Face Enrollment Logs
    # ------------------------------------------------------------------
    def list_face_enrollment_logs_month(
        self,
        company_id: UUID,
        year_month: str,
        page_size: int = 50,
        paging_state: Optional[str] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.ONE
    ) -> Dict[str, Any]:
        """
        Lấy danh sách face enrollment logs theo tháng (partition), có phân trang.
        """
        rows, next_state = self._paged_query(
            "get_month_face_enrollment_logs_company",
            (company_id, year_month),
            page_size=page_size,
            paging_state=paging_state,
            consistency=consistency
        )
        return {
            "items": rows,
            "next_paging_state": next_state,
            "page_size": page_size
        }

    def list_face_enrollment_logs_between(
        self,
        company_id: UUID,
        year_month: str,
        start_time: datetime,
        end_time: datetime,
        page_size: int = 50,
        paging_state: Optional[str] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.ONE
    ) -> Dict[str, Any]:
        """
        Lấy danh sách face enrollment logs trong khoảng thời gian thuộc cùng partition year_month.
        Lưu ý: start_time và end_time phải nằm trong cùng year_month partition để có hiệu quả.
        """
        rows, next_state = self._paged_query(
            "get_month_face_enrollment_logs_with_around_time",
            (company_id, year_month, start_time, end_time),
            page_size=page_size,
            paging_state=paging_state,
            consistency=consistency
        )
        return {
            "items": rows,
            "next_paging_state": next_state,
            "page_size": page_size,
            "filter": {
                "start_time": start_time.isoformat(),
                "end_time": end_time.isoformat()
            }
        }

    # ------------------------------------------------------------------
    # Select Audit Logs
    # ------------------------------------------------------------------
    def list_audit_logs_month(
        self,
        company_id: UUID,
        year_month: str,
        page_size: int = 50,
        paging_state: Optional[str] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.ONE
    ) -> Dict[str, Any]:
        """
        Lấy audit logs theo tháng với phân trang.
        """
        rows, next_state = self._paged_query(
            "get_month_audit_logs_company",
            (company_id, year_month),
            page_size=page_size,
            paging_state=paging_state,
            consistency=consistency
        )
        return {
            "items": rows,
            "next_paging_state": next_state,
            "page_size": page_size
        }

    def list_audit_logs_between(
        self,
        company_id: UUID,
        year_month: str,
        start_time: datetime,
        end_time: datetime,
        page_size: int = 50,
        paging_state: Optional[str] = None,
        consistency: ConsistencyLevel = ConsistencyLevel.ONE
    ) -> Dict[str, Any]:
        """
        Lấy audit logs trong khoảng thời gian (phải cùng partition year_month).
        """
        rows, next_state = self._paged_query(
            "get_month_audit_logs_with_around_time",
            (company_id, year_month, start_time, end_time),
            page_size=page_size,
            paging_state=paging_state,
            consistency=consistency
        )
        return {
            "items": rows,
            "next_paging_state": next_state,
            "page_size": page_size,
            "filter": {
                "start_time": start_time.isoformat(),
                "end_time": end_time.isoformat()
            }
        }

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