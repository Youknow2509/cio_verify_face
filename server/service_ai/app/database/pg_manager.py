"""
    PostgreSQL Database Manager Module
"""

import logging
from typing import List
from uuid import UUID
from sqlalchemy import QueuePool, create_engine
from app.core.config import settings
from sqlalchemy.orm import sessionmaker, Session

from app.models.schemas import FaceProfileResponse
from app.utils.cache import get_employee_company_cache_key, get_ttl_cache_short
from sqlalchemy import text

from app.utils.database import get_face_profile_partition_name
logger = logging.getLogger(__name__)

class PGManager:
    def __init__(self):
        """Initialize pG manager"""        
        # Create database engine with connection pooling for better performance
        self.engine = create_engine(
            settings.DATABASE_URL,
            poolclass=QueuePool,
            pool_size=10,
            max_overflow=20,
            pool_pre_ping=True,  # Verify connections before using
            echo=False
        )
        self.SessionLocal = sessionmaker(bind=self.engine)
        if not self.check_connection():
            logger.error("Database connection failed during PGManager initialization")
            raise ConnectionError("Failed to connect to the database")
        
        logger.info("PGManager initialized successfully")
    
        
    def check_connection(self) -> bool:
        """Check database connection"""
        try:
            with self.engine.connect() as connection:
                connection.execute(text("SELECT 1"))
            return True
        except Exception as e:
            logger.error(f"Database connection check failed: {e}")
            return False
        
    def check_employee_exist_in_company(self, company_id: UUID, employee_id: UUID) -> bool:
        """
            Check if an employee exists in a company with tables:
            CREATE TABLE IF NOT EXISTS employees (
                employee_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
                company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
                employee_code VARCHAR(50) NOT NULL,
                department VARCHAR(100),
                position VARCHAR(100),
                hire_date DATE,
                salary DECIMAL(12,2),
                status int2 DEFAULT 0 NOT NULL CHECK (status IN (0, 1, 2)), -- 0: active, 1: inactive, 2: on leave
                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
            );
        """
        with self.SessionLocal() as session:
            sql_raw = text("""
                SELECT employee_id FROM employees WHERE company_id = :company_id AND employee_id = :employee_id
            """)
            result = session.execute(
                sql_raw,
                {"company_id": str(company_id), "employee_id": str(employee_id)}
            ).fetchone()
            if result is not None:
                return True
            return False
    
    def get_profile_face_employee(self, 
        company_id: UUID, 
        employee_id: UUID,
        page_size: int = 20,
        page_number: int = 1
    ) -> List[FaceProfileResponse]:
        """
            Get face profiles of an employee in a company with tables:
            CREATE TABLE face_profiles (
                profile_id UUID NOT NULL,
                user_id UUID NOT NULL,
                company_id UUID,
                embedding vector(512) NOT NULL,
                embedding_version VARCHAR(50) NOT NULL,
                enroll_image_path TEXT,
                is_primary BOOLEAN DEFAULT false NOT NULL,
                quality_score FLOAT,
                meta_data JSONB DEFAULT '{}' NOT NULL,
                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
                updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
                deleted_at TIMESTAMP WITH TIME ZONE,
                indexed BOOLEAN DEFAULT false NOT NULL,
                index_version INTEGER DEFAULT 0 NOT NULL,
                PRIMARY KEY (profile_id, company_id)
            ) PARTITION BY LIST (company_id);
        """
        limit = page_size
        offset = (page_number - 1) * page_size
        # Check table partition exists
        partition_table_name = get_face_profile_partition_name(company_id)
        with self.SessionLocal() as session:
            sql_raw = text(f"""
                SELECT
                    profile_id, user_id, company_id, embedding_version,
                    is_primary, created_at, updated_at, deleted_at,
                    meta_data, quality_score
                FROM {partition_table_name}
                WHERE company_id = :company_id AND user_id = :employee_id
                LIMIT :limit OFFSET :offset
            """)
            result = session.execute(
                sql_raw,
                {
                    "company_id": str(company_id),
                    "employee_id": str(employee_id),
                    "limit": limit,
                    "offset": offset
                }
            ).fetchall()
            face_profiles = []
            for row in result:
                face_profile = FaceProfileResponse(
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
                face_profiles.append(face_profile)
            return face_profiles
        