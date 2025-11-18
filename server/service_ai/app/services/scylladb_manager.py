"""
ScyllaDB manager for storing authentication states
"""
import logging
from typing import Optional, Dict, Any
from datetime import datetime
from uuid import UUID
from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider
from cassandra.query import SimpleStatement, ConsistencyLevel

from app.core.config import settings

logger = logging.getLogger(__name__)


class ScyllaDBManager:
    """Manager for ScyllaDB operations"""
    
    def __init__(self):
        """Initialize ScyllaDB connection"""
        self.cluster = None
        self.session = None
        self._connect()
    
    def _connect(self):
        """Connect to ScyllaDB"""
        try:
            hosts = settings.SCYLLADB_HOSTS.split(',')
            logger.info(f"Connecting to ScyllaDB at {hosts}")
            
            # Create cluster connection
            self.cluster = Cluster(
                hosts,
                port=settings.SCYLLADB_PORT,
                auth_provider=PlainTextAuthProvider(
                    username=settings.SCYLLADB_USERNAME,
                    password=settings.SCYLLADB_PASSWORD
                ),
                protocol_version=4
            )
            
            # Create session
            self.session = self.cluster.connect()
            
            # Create keyspace if not exists
            self._create_keyspace()
            
            # Use keyspace
            self.session.set_keyspace(settings.SCYLLADB_KEYSPACE)
            
            # Create tables
            self._create_tables()
            
            logger.info("ScyllaDB connection established successfully")
            
        except Exception as e:
            logger.error(f"Error connecting to ScyllaDB: {e}")
            raise
    
    def _create_keyspace(self):
        """Create keyspace if not exists"""
        try:
            query = f"""
                CREATE KEYSPACE IF NOT EXISTS {settings.SCYLLADB_KEYSPACE}
                WITH replication = {{
                    'class': 'SimpleStrategy',
                    'replication_factor': 1
                }}
            """
            self.session.execute(query)
            logger.info(f"Keyspace {settings.SCYLLADB_KEYSPACE} created/verified")
        except Exception as e:
            logger.error(f"Error creating keyspace: {e}")
            raise
    
    def _create_tables(self):
        """Create required tables"""
        try:
            # Authentication verification states table
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
                )
            """
            self.session.execute(auth_states_query)
            
            # Verification history by user
            user_verifications_query = """
                CREATE TABLE IF NOT EXISTS user_verifications (
                    company_id UUID,
                    user_id UUID,
                    timestamp TIMESTAMP,
                    verification_id UUID,
                    verified BOOLEAN,
                    similarity_score DOUBLE,
                    device_id TEXT,
                    PRIMARY KEY ((company_id, user_id), timestamp)
                ) WITH CLUSTERING ORDER BY (timestamp DESC)
            """
            self.session.execute(user_verifications_query)
            
            # Enrollment states table
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
                )
            """
            self.session.execute(enrollment_states_query)
            
            # Audit logs table - for all operations
            # This replaces PostgreSQL's face_audit_logs table
            audit_logs_query = """
                CREATE TABLE IF NOT EXISTS audit_logs (
                    log_id UUID,
                    timestamp TIMESTAMP,
                    company_id UUID,
                    operation TEXT,
                    status TEXT,
                    profile_id UUID,
                    user_id UUID,
                    device_id TEXT,
                    similarity_score DOUBLE,
                    liveness_score DOUBLE,
                    quality_score DOUBLE,
                    metadata MAP<TEXT, TEXT>,
                    error_message TEXT,
                    PRIMARY KEY ((company_id, operation), timestamp, log_id)
                ) WITH CLUSTERING ORDER BY (timestamp DESC, log_id ASC)
            """
            self.session.execute(audit_logs_query)
            
            # Audit logs by user - for fast user-specific queries
            user_audit_logs_query = """
                CREATE TABLE IF NOT EXISTS user_audit_logs (
                    company_id UUID,
                    user_id UUID,
                    timestamp TIMESTAMP,
                    log_id UUID,
                    operation TEXT,
                    status TEXT,
                    profile_id UUID,
                    device_id TEXT,
                    similarity_score DOUBLE,
                    liveness_score DOUBLE,
                    quality_score DOUBLE,
                    error_message TEXT,
                    PRIMARY KEY ((company_id, user_id), timestamp, log_id)
                ) WITH CLUSTERING ORDER BY (timestamp DESC, log_id ASC)
            """
            self.session.execute(user_audit_logs_query)
            
            logger.info("ScyllaDB tables created/verified")
            
        except Exception as e:
            logger.error(f"Error creating tables: {e}")
            raise
    
    def save_verification_state(
        self,
        verification_id: UUID,
        company_id: UUID,
        user_id: Optional[UUID],
        profile_id: Optional[UUID],
        device_id: Optional[str],
        status: str,
        verified: bool,
        similarity_score: Optional[float] = None,
        liveness_score: Optional[float] = None,
        metadata: Optional[Dict[str, str]] = None,
        image_path: Optional[str] = None
    ):
        """
        Save verification state to ScyllaDB
        
        Args:
            verification_id: Unique verification ID
            user_id: User ID (if matched)
            profile_id: Profile ID (if matched)
            device_id: Device identifier
            status: Verification status (match, no_match, failed)
            verified: Whether verification passed
            similarity_score: Similarity score
            liveness_score: Liveness score
            metadata: Additional metadata
            image_path: Path to stored image in MinIO
        """
        try:
            query = """
                INSERT INTO authentication_states (
                    verification_id, company_id, user_id, profile_id, device_id,
                    status, verified, similarity_score, liveness_score,
                    timestamp, metadata, image_path
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            
            statement = SimpleStatement(query, consistency_level=ConsistencyLevel.QUORUM)
            
            self.session.execute(statement, (
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
                metadata or {},
                image_path
            ))
            
            # Also save to user verification history if user matched
            if user_id and verified:
                user_query = """
                    INSERT INTO user_verifications (
                        company_id, user_id, timestamp, verification_id, verified,
                        similarity_score, device_id
                    ) VALUES (%s, %s, %s, %s, %s, %s, %s)
                """
                user_statement = SimpleStatement(user_query, consistency_level=ConsistencyLevel.QUORUM)
                self.session.execute(user_statement, (
                    company_id,
                    user_id,
                    datetime.utcnow(),
                    verification_id,
                    verified,
                    similarity_score,
                    device_id
                ))
            
            logger.debug(f"Saved verification state: {verification_id}")
            
        except Exception as e:
            logger.error(f"Error saving verification state: {e}")
            # Don't raise - this is logging, not critical
    
    def save_enrollment_state(
        self,
        profile_id: UUID,
        company_id: UUID,
        user_id: UUID,
        device_id: Optional[str],
        status: str,
        quality_score: Optional[float] = None,
        metadata: Optional[Dict[str, str]] = None,
        image_path: Optional[str] = None
    ):
        """
        Save enrollment state to ScyllaDB
        
        Args:
            profile_id: Profile ID
            user_id: User ID
            device_id: Device identifier
            status: Enrollment status (ok, duplicate, failed)
            quality_score: Image quality score
            metadata: Additional metadata
            image_path: Path to stored image in MinIO
        """
        try:
            query = """
                INSERT INTO enrollment_states (
                    profile_id, company_id, user_id, device_id, status,
                    quality_score, timestamp, metadata, image_path
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            
            statement = SimpleStatement(query, consistency_level=ConsistencyLevel.QUORUM)
            
            self.session.execute(statement, (
                profile_id,
                company_id,
                user_id,
                device_id,
                status,
                quality_score,
                datetime.utcnow(),
                metadata or {},
                image_path
            ))
            
            logger.debug(f"Saved enrollment state: {profile_id}")
            
        except Exception as e:
            logger.error(f"Error saving enrollment state: {e}")
            # Don't raise - this is logging, not critical
    
    def get_verification_state(self, verification_id: UUID) -> Optional[Dict[str, Any]]:
        """Get verification state by ID"""
        try:
            query = "SELECT * FROM authentication_states WHERE verification_id = %s"
            rows = self.session.execute(query, (verification_id,))
            row = rows.one()
            
            if row:
                return {
                    'verification_id': row.verification_id,
                    'user_id': row.user_id,
                    'profile_id': row.profile_id,
                    'device_id': row.device_id,
                    'status': row.status,
                    'verified': row.verified,
                    'similarity_score': row.similarity_score,
                    'liveness_score': row.liveness_score,
                    'timestamp': row.timestamp,
                    'metadata': row.metadata,
                    'image_path': row.image_path
                }
            return None
            
        except Exception as e:
            logger.error(f"Error getting verification state: {e}")
            return None
    
    def get_user_verifications(
        self,
        user_id: UUID,
        company_id: UUID,
        limit: int = 100
    ) -> list:
        """Get recent verifications for a user"""
        try:
            query = """
                SELECT * FROM user_verifications 
                WHERE company_id = %s AND user_id = %s 
                LIMIT %s
            """
            rows = self.session.execute(query, (company_id, user_id, limit))
            
            return [
                {
                    'user_id': row.user_id,
                    'timestamp': row.timestamp,
                    'verification_id': row.verification_id,
                    'verified': row.verified,
                    'similarity_score': row.similarity_score,
                    'device_id': row.device_id
                }
                for row in rows
            ]
            
        except Exception as e:
            logger.error(f"Error getting user verifications: {e}")
            return []
    
    def save_audit_log(
        self,
        log_id: UUID,
        company_id: UUID,
        operation: str,
        status: str,
        profile_id: Optional[UUID] = None,
        user_id: Optional[UUID] = None,
        device_id: Optional[str] = None,
        similarity_score: Optional[float] = None,
        liveness_score: Optional[float] = None,
        quality_score: Optional[float] = None,
        metadata: Optional[Dict[str, str]] = None,
        error_message: Optional[str] = None
    ):
        """
        Save audit log to ScyllaDB
        
        This replaces PostgreSQL's face_audit_logs table for better performance
        
        Args:
            log_id: Unique log ID
            operation: Operation type (enroll, verify, update, delete)
            status: Operation status (success, failed, duplicate, etc.)
            profile_id: Profile ID (if applicable)
            user_id: User ID (if applicable)
            device_id: Device identifier
            similarity_score: Similarity score
            liveness_score: Liveness score
            quality_score: Quality score
            metadata: Additional metadata
            error_message: Error message (if failed)
        """
        try:
            timestamp = datetime.utcnow()
            
            # Save to main audit_logs table (indexed by operation)
            query = """
                INSERT INTO audit_logs (
                    log_id, timestamp, company_id, operation, status, profile_id, user_id,
                    device_id, similarity_score, liveness_score, quality_score,
                    metadata, error_message
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            statement = SimpleStatement(query, consistency_level=ConsistencyLevel.QUORUM)
            self.session.execute(statement, (
                log_id,
                timestamp,
                company_id,
                operation,
                status,
                profile_id,
                user_id,
                device_id,
                similarity_score,
                liveness_score,
                quality_score,
                metadata or {},
                error_message
            ))
            
            # Also save to user_audit_logs if user_id is present (for fast user queries)
            if user_id:
                user_query = """
                    INSERT INTO user_audit_logs (
                        company_id, user_id, timestamp, log_id, operation, status, profile_id,
                        device_id, similarity_score, liveness_score, quality_score,
                        error_message
                    ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                """
                user_statement = SimpleStatement(user_query, consistency_level=ConsistencyLevel.QUORUM)
                self.session.execute(user_statement, (
                    company_id,
                    user_id,
                    timestamp,
                    log_id,
                    operation,
                    status,
                    profile_id,
                    device_id,
                    similarity_score,
                    liveness_score,
                    quality_score,
                    error_message
                ))
            
            logger.debug(f"Saved audit log: {log_id} - {operation}/{status}")
            
        except Exception as e:
            logger.error(f"Error saving audit log: {e}")
            # Don't raise - this is logging, not critical
    
    def get_audit_logs(
        self,
        company_id: UUID,
        operation: Optional[str] = None,
        limit: int = 100
    ) -> list:
        """
        Get audit logs by operation
        
        Args:
            operation: Operation type (enroll, verify, update, delete)
            limit: Maximum number of logs to return
            
        Returns:
            List of audit log dictionaries
        """
        try:
            if operation:
                query = """
                    SELECT * FROM audit_logs 
                    WHERE company_id = %s AND operation = %s 
                    LIMIT %s
                """
                rows = self.session.execute(query, (company_id, operation, limit))
            else:
                query = """
                    SELECT * FROM audit_logs 
                    WHERE company_id = %s 
                    LIMIT %s
                """
                rows = self.session.execute(query, (company_id, limit))
            
            return [
                {
                    'log_id': row.log_id,
                    'timestamp': row.timestamp,
                    'operation': row.operation,
                    'status': row.status,
                    'profile_id': row.profile_id,
                    'user_id': row.user_id,
                    'device_id': row.device_id,
                    'similarity_score': row.similarity_score,
                    'liveness_score': row.liveness_score,
                    'quality_score': row.quality_score,
                    'metadata': row.metadata,
                    'error_message': row.error_message
                }
                for row in rows
            ]
            
        except Exception as e:
            logger.error(f"Error getting audit logs: {e}")
            return []
    
    def get_user_audit_logs(
        self,
        user_id: UUID,
        company_id: UUID,
        limit: int = 100
    ) -> list:
        """
        Get audit logs for a specific user
        
        Args:
            user_id: User ID
            limit: Maximum number of logs to return
            
        Returns:
            List of audit log dictionaries
        """
        try:
            query = """
                SELECT * FROM user_audit_logs 
                WHERE company_id = %s AND user_id = %s 
                LIMIT %s
            """
            rows = self.session.execute(query, (company_id, user_id, limit))
            
            return [
                {
                    'user_id': row.user_id,
                    'timestamp': row.timestamp,
                    'log_id': row.log_id,
                    'operation': row.operation,
                    'status': row.status,
                    'profile_id': row.profile_id,
                    'device_id': row.device_id,
                    'similarity_score': row.similarity_score,
                    'liveness_score': row.liveness_score,
                    'quality_score': row.quality_score,
                    'error_message': row.error_message
                }
                for row in rows
            ]
            
        except Exception as e:
            logger.error(f"Error getting user audit logs: {e}")
            return []
    
    def close(self):
        """Close ScyllaDB connection"""
        try:
            if self.cluster:
                self.cluster.shutdown()
                logger.info("ScyllaDB connection closed")
        except Exception as e:
            logger.error(f"Error closing ScyllaDB connection: {e}")
