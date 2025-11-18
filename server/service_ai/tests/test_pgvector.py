"""
Tests for PgVectorManager - PostgreSQL pgvector-based vector search
Run with: pytest tests/test_pgvector.py -v
"""
import pytest
import uuid
import numpy as np
from datetime import datetime

from app.services.pgvector_manager import PgVectorManager
from app.core.config import settings
from app.models.database import Base, FaceProfile
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker


class TestPgVectorManager:
    """Test PgVectorManager functionality"""
    
    @pytest.fixture(autouse=True)
    def setup(self):
        """Setup test environment"""
        try:
            # Initialize pgvector manager
            self.manager = PgVectorManager()
            
            # Create test database session
            self.engine = create_engine(settings.DATABASE_URL)
            self.SessionLocal = sessionmaker(bind=self.engine)
            
            # Verify pgvector extension is available
            with self.engine.connect() as conn:
                result = conn.execute(text("SELECT extname FROM pg_extension WHERE extname = 'vector'"))
                if result.fetchone() is None:
                    pytest.skip("pgvector extension not available")
            
            # Clean up any test data
            self._cleanup_test_data()
            
            yield
            
            # Clean up after tests
            self._cleanup_test_data()
            
        except Exception as e:
            pytest.skip(f"Database or pgvector not available: {e}")
    
    def _cleanup_test_data(self):
        """Clean up test profiles"""
        db = self.SessionLocal()
        try:
            # Delete test profiles (use a special marker in metadata)
            db.execute(text("""
                DELETE FROM face_profiles 
                WHERE meta_data->>'test_marker' = 'pgvector_test'
            """))
            db.commit()
        except Exception as e:
            print(f"Cleanup error: {e}")
            db.rollback()
        finally:
            db.close()
    
    def _create_test_profile(self, embedding=None, user_id=None, is_primary=False):
        """Helper to create a test profile"""
        if embedding is None:
            # Generate random normalized embedding
            embedding = np.random.randn(512).astype(np.float32)
            embedding = embedding / np.linalg.norm(embedding)
        
        if user_id is None:
            user_id = uuid.uuid4()
        
        profile_id = uuid.uuid4()
        
        db = self.SessionLocal()
        try:
            # Create profile
            profile = FaceProfile(
                profile_id=profile_id,
                user_id=user_id,
                embedding=embedding.tolist(),
                embedding_version="arcface_r100_v1",
                is_primary=is_primary,
                quality_score=0.95,
                meta_data={"test_marker": "pgvector_test"},
                indexed=False
            )
            db.add(profile)
            db.commit()
            
            return profile_id, user_id, embedding
        finally:
            db.close()
    
    def test_initialization(self):
        """Test PgVectorManager initialization"""
        assert self.manager is not None
        assert self.manager.dimension == settings.EMBEDDING_DIMENSION
        assert self.manager.engine is not None
    
    def test_add_embedding(self):
        """Test adding an embedding"""
        profile_id, user_id, embedding = self._create_test_profile()
        
        # Add embedding to index
        self.manager.add_embedding(
            str(profile_id),
            str(user_id),
            embedding,
            is_primary=False
        )
        
        # Verify it was indexed
        db = self.SessionLocal()
        try:
            result = db.execute(
                text("SELECT indexed FROM face_profiles WHERE profile_id = :pid"),
                {"pid": str(profile_id)}
            )
            row = result.fetchone()
            assert row is not None
            assert row[0] is True  # indexed flag should be True
        finally:
            db.close()
    
    def test_search_single_result(self):
        """Test searching for a single embedding"""
        # Create a test profile
        profile_id, user_id, embedding = self._create_test_profile()
        
        # Add to index
        self.manager.add_embedding(
            str(profile_id),
            str(user_id),
            embedding,
            is_primary=True
        )
        
        # Search with the same embedding (should find itself with high similarity)
        results = self.manager.search(embedding, k=5)
        
        assert len(results) >= 1
        best_match = results[0]
        assert best_match['profile_id'] == str(profile_id)
        assert best_match['user_id'] == str(user_id)
        assert best_match['is_primary'] is True
        assert best_match['similarity'] > 0.99  # Should be very similar to itself
    
    def test_search_multiple_results(self):
        """Test searching with multiple embeddings in the database"""
        # Create multiple test profiles
        profiles = []
        for i in range(5):
            profile_id, user_id, embedding = self._create_test_profile()
            self.manager.add_embedding(
                str(profile_id),
                str(user_id),
                embedding,
                is_primary=(i == 0)
            )
            profiles.append((profile_id, user_id, embedding))
        
        # Search with the first embedding
        query_embedding = profiles[0][2]
        results = self.manager.search(query_embedding, k=3)
        
        assert len(results) <= 3
        assert len(results) >= 1
        
        # Best match should be the first profile
        best_match = results[0]
        assert best_match['profile_id'] == str(profiles[0][0])
        assert best_match['similarity'] > 0.99
    
    def test_search_with_similar_embeddings(self):
        """Test searching with similar but not identical embeddings"""
        # Create a base embedding
        base_embedding = np.random.randn(512).astype(np.float32)
        base_embedding = base_embedding / np.linalg.norm(base_embedding)
        
        # Create profile with base embedding
        profile_id, user_id, _ = self._create_test_profile(base_embedding)
        self.manager.add_embedding(
            str(profile_id),
            str(user_id),
            base_embedding
        )
        
        # Create a slightly modified embedding (add small noise)
        noise = np.random.randn(512).astype(np.float32) * 0.1
        query_embedding = base_embedding + noise
        query_embedding = query_embedding / np.linalg.norm(query_embedding)
        
        # Search with the modified embedding
        results = self.manager.search(query_embedding, k=5)
        
        assert len(results) >= 1
        best_match = results[0]
        assert best_match['profile_id'] == str(profile_id)
        # Should still be quite similar (> 0.8)
        assert best_match['similarity'] > 0.8
    
    def test_remove_embedding(self):
        """Test removing an embedding from the index"""
        # Create and index a profile
        profile_id, user_id, embedding = self._create_test_profile()
        self.manager.add_embedding(str(profile_id), str(user_id), embedding)
        
        # Verify it's indexed
        results = self.manager.search(embedding, k=5)
        assert len(results) >= 1
        assert any(r['profile_id'] == str(profile_id) for r in results)
        
        # Remove from index
        self.manager.remove_embedding(str(profile_id))
        
        # Verify it's no longer in search results
        results = self.manager.search(embedding, k=5)
        assert not any(r['profile_id'] == str(profile_id) for r in results)
    
    def test_rebuild_index(self):
        """Test rebuilding the entire index"""
        # Create multiple profiles
        embeddings = []
        for i in range(3):
            profile_id, user_id, embedding = self._create_test_profile()
            embeddings.append((
                str(profile_id),
                str(user_id),
                embedding,
                i == 0  # first is primary
            ))
        
        # Rebuild index
        self.manager.rebuild_index(embeddings)
        
        # Verify all are indexed
        db = self.SessionLocal()
        try:
            result = db.execute(text("""
                SELECT COUNT(*) 
                FROM face_profiles 
                WHERE meta_data->>'test_marker' = 'pgvector_test'
                  AND indexed = true
            """))
            count = result.scalar()
            assert count == 3
        finally:
            db.close()
        
        # Verify search works
        query_embedding = embeddings[0][2]
        results = self.manager.search(query_embedding, k=5)
        assert len(results) >= 1
    
    def test_get_size(self):
        """Test getting the size of the index"""
        initial_size = self.manager.get_size()
        
        # Add some profiles
        for i in range(3):
            profile_id, user_id, embedding = self._create_test_profile()
            self.manager.add_embedding(str(profile_id), str(user_id), embedding)
        
        # Check size increased
        new_size = self.manager.get_size()
        assert new_size == initial_size + 3
    
    def test_clear_index(self):
        """Test clearing the index"""
        # Add some profiles
        for i in range(3):
            profile_id, user_id, embedding = self._create_test_profile()
            self.manager.add_embedding(str(profile_id), str(user_id), embedding)
        
        # Verify they're indexed
        size_before = self.manager.get_size()
        assert size_before >= 3
        
        # Clear index
        self.manager.clear()
        
        # Verify size is 0 (only for test profiles)
        db = self.SessionLocal()
        try:
            result = db.execute(text("""
                SELECT COUNT(*) 
                FROM face_profiles 
                WHERE meta_data->>'test_marker' = 'pgvector_test'
                  AND indexed = true
            """))
            count = result.scalar()
            assert count == 0
        finally:
            db.close()
    
    def test_search_empty_index(self):
        """Test searching when index is empty"""
        # Clear all test data
        self.manager.clear()
        
        # Create a random embedding
        embedding = np.random.randn(512).astype(np.float32)
        embedding = embedding / np.linalg.norm(embedding)
        
        # Search should return empty list
        results = self.manager.search(embedding, k=5)
        assert len(results) == 0
    
    def test_save_index_noop(self):
        """Test that save_index is a no-op (pgvector auto-persists)"""
        # This should not raise any errors
        self.manager.save_index()
        # No assertion needed - just verify it doesn't crash
    
    def test_last_rebuild_property(self):
        """Test the last_rebuild property"""
        last_rebuild = self.manager.last_rebuild
        assert isinstance(last_rebuild, datetime)
    
    def test_normalized_embeddings(self):
        """Test that embeddings are properly normalized"""
        # Create an unnormalized embedding
        embedding = np.array([1.0, 2.0, 3.0] + [0.0] * 509, dtype=np.float32)
        assert not np.isclose(np.linalg.norm(embedding), 1.0)
        
        profile_id, user_id, _ = self._create_test_profile(embedding)
        
        # Add to index (should normalize internally)
        self.manager.add_embedding(str(profile_id), str(user_id), embedding)
        
        # Search should work correctly
        results = self.manager.search(embedding, k=1)
        assert len(results) == 1
        assert results[0]['profile_id'] == str(profile_id)
    
    def test_concurrent_access(self):
        """Test that multiple operations work correctly (simulated concurrency)"""
        # Create multiple profiles
        profiles = []
        for i in range(5):
            profile_id, user_id, embedding = self._create_test_profile()
            self.manager.add_embedding(str(profile_id), str(user_id), embedding)
            profiles.append((profile_id, embedding))
        
        # Perform multiple searches
        for profile_id, embedding in profiles:
            results = self.manager.search(embedding, k=3)
            assert len(results) >= 1
            assert results[0]['profile_id'] == str(profile_id)
        
        # Remove some profiles
        for i in range(2):
            self.manager.remove_embedding(str(profiles[i][0]))
        
        # Verify removed profiles are not in search
        for i in range(2):
            results = self.manager.search(profiles[i][1], k=5)
            assert not any(r['profile_id'] == str(profiles[i][0]) for r in results)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
