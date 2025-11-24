"""
PostgreSQL pgvector manager for efficient face embedding search in microservices
"""
import logging
import numpy as np
from typing import List, Dict, Optional
from datetime import datetime
from sqlalchemy import create_engine, text, and_
from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy.pool import QueuePool
from app.utils.database import get_face_profile_partition_name

from app.core.config import settings

logger = logging.getLogger(__name__)


class PgVectorManager:
    """Manage face embeddings using PostgreSQL pgvector extension"""
    
    def __init__(self):
        """Initialize pgvector manager"""
        self.dimension = settings.EMBEDDING_DIMENSION
        
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
        self.index_version = settings.VECTOR_DB_INDEX_VERSION
        self._initialize_pgvector()
        logger.info("PgVectorManager initialized successfully")
    
    def check_connection(self):
        """Check connection to PostgreSQL database"""
        try:
            with self.engine.connect() as conn:
                conn.execute(text("SELECT 1"))
            logger.info("PostgreSQL connection successful")
        except Exception as e:
            logger.error(f"PostgreSQL connection failed: {e}")
            raise
    
    def _initialize_pgvector(self):
        """Initialize pgvector extension if not already enabled"""
        try:
            with self.engine.connect() as conn:
                # Enable pgvector extension
                conn.execute(text("CREATE EXTENSION IF NOT EXISTS vector"))
                conn.commit()
                
                # Verify the face_profiles table has the vector column
                # This will be handled by the migration script
                logger.info("pgvector extension verified")
        except Exception as e:
            logger.error(f"Error initializing pgvector: {e}")
            raise
    
    def ensure_company_partition(self, company_id: str) -> bool:
        """
        Ensure that a partition exists for the given company_id.
        Creates the partition if it doesn't exist.
        
        Args:
            company_id: Company UUID as string
            
        Returns:
            True if partition exists or was created successfully, False otherwise
        """
        db = self.SessionLocal()
        try:
            partition_name = get_face_profile_partition_name(company_id)
            
            # Check if partition already exists
            check_query = text("""
                SELECT EXISTS (
                    SELECT 1 FROM pg_tables 
                    WHERE tablename = :partition_name
                );
            """)
            result = db.execute(check_query, {"partition_name": partition_name}).scalar()
            
            if result:
                logger.debug(f"Partition {partition_name} already exists")
                return True
            
            # Create partition - table name and UUID value must be inserted directly
            # (cannot use parameter binding for table names in DDL)
            create_partition_query = text(f"""
                CREATE TABLE IF NOT EXISTS {partition_name}
                PARTITION OF face_profiles
                FOR VALUES IN ('{company_id}');
            """)
            db.execute(create_partition_query)
            db.commit()
            logger.info(f"Created partition {partition_name} for company_id: {company_id}")
            return True
            
        except Exception as e:
            logger.error(f"Error ensuring partition for company {company_id}: {e}")
            db.rollback()
            return False
        finally:
            db.close()
    
    def add_embedding(self, profile_id: str, company_id: str, user_id: str, embedding: np.ndarray, 
                     is_primary: bool = False):
        """
        Add or update an embedding in the database
        
        Args:
            profile_id: Profile UUID as string
            company_id: Company UUID as string
            user_id: User UUID as string
            embedding: Face embedding vector
            is_primary: Whether this is the primary profile
            
        Note: This method updates the embedding vector in an existing profile.
        The actual profile creation happens in face_service.py
        """
        db = self.SessionLocal()
        try:
            # Ensure embedding is normalized for cosine similarity
            embedding = embedding.flatten().astype(np.float32)
            embedding = embedding / np.linalg.norm(embedding)
            
            # Convert numpy array to list for PostgreSQL
            embedding_list = embedding.tolist()
            
            # Update the embedding vector in the face_profiles table
            query = text("""
                UPDATE face_profiles
                SET embedding = CAST(:embedding AS vector),
                    indexed = true,
                    updated_at = CURRENT_TIMESTAMP
                WHERE profile_id = CAST(:profile_id AS uuid) 
                    AND company_id = CAST(:company_id AS uuid)
                    AND deleted_at IS NULL
            """)
            
            result = db.execute(
                query,
                {
                    "embedding": str(embedding_list),
                    "profile_id": profile_id,
                    "company_id": company_id
                }
            )
            db.commit()
            
            if result.rowcount == 0:
                logger.warning(f"Profile {profile_id} not found or already deleted for company {company_id}")
            else:
                logger.debug(f"Updated embedding for profile {profile_id} and company {company_id}")
                
        except Exception as e:
            logger.error(f"Error adding/updating embedding: {e}")
            db.rollback()
            raise
        finally:
            db.close()
    
    def remove_embedding(self, profile_id: str, company_id: str):
        """
        Mark embedding as not indexed (soft removal from vector search)
        
        Args:
            profile_id: Profile UUID as string
            company_id: Company UUID as string
            
        Note: We don't actually delete the row, just mark it as not indexed.
        Actual deletion is handled by soft delete in face_service.py
        """
        db = self.SessionLocal()
        try:
            query = text("""
                UPDATE face_profiles
                SET indexed = false,
                    updated_at = CURRENT_TIMESTAMP
                WHERE profile_id = CAST(:profile_id AS uuid)
                  AND company_id = CAST(:company_id AS uuid)
            """)
            
            db.execute(query, {"profile_id": profile_id, "company_id": company_id})
            db.commit()
            
            logger.info(f"Marked embedding as not indexed for profile {profile_id} and company {company_id}")
        except Exception as e:
            logger.error(f"Error removing embedding: {e}")
            db.rollback()
            raise
        finally:
            db.close()
    
    def search(self, company_id: str, embedding: np.ndarray, k: int = 5) -> List[Dict]:
        """
        Search for similar embeddings using cosine distance
        
        Args:
            company_id: Company UUID as string
            embedding: Query embedding vector
            k: Number of results to return
            
        Returns:
            List of matches with profile_id, user_id, similarity, is_primary
        """
        db = self.SessionLocal()
        try:
            # Ensure embedding is normalized
            embedding = embedding.flatten().astype(np.float32)
            embedding = embedding / np.linalg.norm(embedding)
            
            # Convert to list for PostgreSQL
            embedding_list = embedding.tolist()
            
            # Use cosine distance operator (<=>)
            # pgvector's <=> operator computes cosine distance (0 = identical, 2 = opposite)
            # We convert to similarity score (0-1 range) where 1 = identical
            query = text("""
                SELECT 
                    profile_id::text,
                    user_id::text,
                    is_primary,
                    1 - (embedding <=> CAST(:embedding AS vector)) AS similarity
                FROM face_profiles
                WHERE deleted_at IS NULL
                  AND indexed = true
                  AND company_id = CAST(:company_id AS uuid)
                ORDER BY embedding <=> CAST(:embedding AS vector)
                LIMIT :k
            """)
            
            result = db.execute(
                query,
                {
                    "embedding": str(embedding_list),
                    "company_id": company_id,
                    "k": k
                }
            )
            
            # Format results
            matches = []
            for row in result:
                matches.append({
                    'profile_id': row[0],
                    'user_id': row[1],
                    'is_primary': row[2],
                    'similarity': float(row[3])
                })
            
            return matches
            
        except Exception as e:
            logger.error(f"Error searching embeddings: {e}")
            raise
        finally:
            db.close()
    
    def rebuild_index(self, embeddings: List[tuple]):
        """
        Rebuild the vector index (for compatibility with FAISS API)
        
        Args:
            embeddings: List of (profile_id, user_id, embedding, is_primary) tuples
            
        Note: With pgvector, we don't need to rebuild an index structure.
        This method updates all embeddings in the database and marks them as indexed.
        The database index on the embedding column is automatically maintained.
        """
        db = self.SessionLocal()
        try:
            logger.info(f"Rebuilding index with {len(embeddings)} embeddings...")
            
            # First, mark all as not indexed
            db.execute(text("UPDATE face_profiles SET indexed = false WHERE deleted_at IS NULL"))
            
            # Update each embedding
            updated_count = 0
            for profile_id, user_id, embedding, is_primary in embeddings:
                # Normalize embedding
                embedding = embedding.flatten().astype(np.float32)
                embedding = embedding / np.linalg.norm(embedding)
                embedding_list = embedding.tolist()
                
                # Update the embedding
                query = text("""
                    UPDATE face_profiles
                    SET embedding = CAST(:embedding AS vector),
                        indexed = true,
                        updated_at = CURRENT_TIMESTAMP
                    WHERE profile_id = CAST(:profile_id AS uuid)
                      AND deleted_at IS NULL
                """)
                
                result = db.execute(
                    query,
                    {
                        "embedding": str(embedding_list),
                        "profile_id": profile_id
                    }
                )
                
                if result.rowcount > 0:
                    updated_count += 1
            
            db.commit()
            logger.info(f"Index rebuilt successfully with {updated_count} embeddings")
            
        except Exception as e:
            logger.error(f"Error rebuilding index: {e}")
            db.rollback()
            raise
        finally:
            db.close()
    
    def save_index(self):
        """
        Save index to disk (for compatibility with FAISS API)
        
        Note: With pgvector, the index is automatically persisted in PostgreSQL.
        This is a no-op method for API compatibility.
        """
        # No-op: PostgreSQL automatically persists data
        logger.debug("Index save called (no-op for pgvector)")
        pass
    
    def get_size(self) -> int:
        """Get the number of indexed embeddings"""
        db = self.SessionLocal()
        try:
            query = text("""
                SELECT COUNT(*)
                FROM face_profiles
                WHERE deleted_at IS NULL
                  AND indexed = true
            """)
            
            result = db.execute(query)
            count = result.scalar()
            
            return count if count else 0
            
        except Exception as e:
            logger.error(f"Error getting index size: {e}")
            return 0
        finally:
            db.close()
    
    def clear(self):
        """Clear all indexed embeddings"""
        db = self.SessionLocal()
        try:
            query = text("""
                UPDATE face_profiles
                SET indexed = false,
                    updated_at = CURRENT_TIMESTAMP
                WHERE deleted_at IS NULL
            """)
            
            db.execute(query)
            db.commit()
            logger.info("Index cleared (all embeddings marked as not indexed)")
            
        except Exception as e:
            logger.error(f"Error clearing index: {e}")
            db.rollback()
            raise
        finally:
            db.close()
    
    @property
    def last_rebuild(self):
        """
        Get last rebuild time (for compatibility with FAISS API)
        
        Note: With pgvector, we return the most recent update time
        """
        db = self.SessionLocal()
        try:
            query = text("""
                SELECT MAX(updated_at)
                FROM face_profiles
                WHERE indexed = true
            """)
            
            result = db.execute(query)
            last_update = result.scalar()
            
            return last_update if last_update else datetime.utcnow()
            
        except Exception as e:
            logger.error(f"Error getting last rebuild time: {e}")
            return datetime.utcnow()
        finally:
            db.close()
