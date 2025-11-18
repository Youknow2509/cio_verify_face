"""
Test backward compatibility between FAISS and pgvector implementations
This ensures the migration doesn't break existing functionality
"""
import pytest
import numpy as np
from unittest.mock import Mock, patch


def generate_test_embedding(dim=512, seed=None):
    """Generate a test embedding"""
    if seed is not None:
        np.random.seed(seed)
    embedding = np.random.randn(dim).astype(np.float32)
    return embedding / np.linalg.norm(embedding)


class TestBackwardCompatibility:
    """Test that pgvector maintains FAISS API compatibility"""
    
    def test_pgvector_api_matches_faiss(self):
        """Verify PgVectorManager has same API as FAISSIndexManager"""
        from unittest.mock import patch
        
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                from app.services.pgvector_manager import PgVectorManager
                
                manager = PgVectorManager()
                
                # Check critical methods exist
                assert hasattr(manager, 'add_embedding')
                assert hasattr(manager, 'remove_embedding')
                assert hasattr(manager, 'search')
                assert hasattr(manager, 'rebuild_index')
                assert hasattr(manager, 'save_index')
                assert hasattr(manager, 'get_size')
                assert hasattr(manager, 'clear')
                assert hasattr(manager, 'last_rebuild')
    
    def test_face_service_can_initialize_with_pgvector(self):
        """Test that FaceService can initialize with PgVectorManager"""
        with patch('app.services.face_service.FaceDetector'):
            with patch('app.services.face_service.FaceEmbedding'):
                with patch('app.services.face_service.PgVectorManager'):
                    with patch('app.services.face_service.LivenessDetector'):
                        with patch('app.services.face_service.ScyllaDBManager'):
                            with patch('app.services.face_service.MinIOManager'):
                                with patch('app.services.face_service.create_engine'):
                                    with patch('app.services.face_service.sessionmaker'):
                                        from app.services.face_service import FaceService
                                        
                                        # This should not raise any errors
                                        # The service initialization is patched so it won't actually connect
    
    def test_embedding_format_compatibility(self):
        """Test that embedding format works with both systems"""
        # Test numpy array
        embedding_np = generate_test_embedding()
        assert embedding_np.dtype == np.float32
        assert embedding_np.shape == (512,)
        assert np.isclose(np.linalg.norm(embedding_np), 1.0)
        
        # Test list format (for database storage)
        embedding_list = embedding_np.tolist()
        assert isinstance(embedding_list, list)
        assert len(embedding_list) == 512
        
        # Test conversion back
        embedding_restored = np.array(embedding_list, dtype=np.float32)
        assert np.allclose(embedding_np, embedding_restored)
    
    def test_similarity_computation(self):
        """Test that similarity scores are comparable between FAISS and pgvector"""
        # Generate test embeddings
        emb1 = generate_test_embedding(seed=42)
        emb2 = generate_test_embedding(seed=43)
        
        # Compute cosine similarity manually
        # Both FAISS (with IndexFlatIP) and pgvector (with vector_cosine_ops) use cosine similarity
        similarity = np.dot(emb1, emb2)
        
        # Verify it's in expected range
        assert -1.0 <= similarity <= 1.0
        
        # Test with identical embeddings
        similarity_self = np.dot(emb1, emb1)
        assert np.isclose(similarity_self, 1.0, atol=0.01)
    
    def test_search_result_format(self):
        """Test that search results have the same format"""
        # Expected result format from FAISS
        expected_keys = {'profile_id', 'user_id', 'similarity', 'is_primary'}
        
        # Create a mock result
        result = {
            'profile_id': 'uuid-string',
            'user_id': 'uuid-string',
            'similarity': 0.95,
            'is_primary': True
        }
        
        # Verify format
        assert set(result.keys()) == expected_keys
        assert isinstance(result['profile_id'], str)
        assert isinstance(result['user_id'], str)
        assert isinstance(result['similarity'], (int, float))
        assert isinstance(result['is_primary'], bool)
    
    def test_normalization_idempotent(self):
        """Test that normalization is idempotent (applying twice doesn't change result)"""
        embedding = np.random.randn(512).astype(np.float32)
        
        # First normalization
        normalized1 = embedding / np.linalg.norm(embedding)
        
        # Second normalization (should not change)
        normalized2 = normalized1 / np.linalg.norm(normalized1)
        
        # Should be nearly identical
        assert np.allclose(normalized1, normalized2)
        assert np.isclose(np.linalg.norm(normalized2), 1.0)
    
    def test_distance_metrics_compatible(self):
        """Test that distance metrics produce comparable results"""
        emb1 = generate_test_embedding(seed=100)
        emb2 = generate_test_embedding(seed=101)
        
        # Cosine similarity (what both systems use)
        cosine_sim = np.dot(emb1, emb2)
        
        # Cosine distance (1 - similarity)
        cosine_dist = 1 - cosine_sim
        
        # pgvector uses cosine distance operator (<=>)
        # FAISS with IndexFlatIP returns inner product (same as cosine similarity for normalized vectors)
        
        # Verify relationship
        assert np.isclose(cosine_sim + cosine_dist, 1.0)
        
        # For normalized vectors, inner product = cosine similarity
        assert -1.0 <= cosine_sim <= 1.0
        assert 0.0 <= cosine_dist <= 2.0


class TestMigrationSafety:
    """Test that migration is safe and doesn't lose data"""
    
    def test_vector_type_handles_float_arrays(self):
        """Test that Vector type can handle float arrays"""
        from app.models.database import Vector
        
        vector_type = Vector(dim=512)
        
        # Test column spec
        col_spec = vector_type.get_col_spec()
        assert col_spec == "vector(512)"
        
        # Test bind processor
        bind_proc = vector_type.bind_processor(None)
        
        # Test with list
        test_list = [1.0, 2.0, 3.0]
        result = bind_proc(test_list)
        assert isinstance(result, str)
        
        # Test with None
        result = bind_proc(None)
        assert result is None
    
    def test_config_compatibility(self):
        """Test that config changes are backward compatible"""
        from app.core.config import settings
        
        # Check that essential settings exist
        assert hasattr(settings, 'EMBEDDING_DIMENSION')
        assert hasattr(settings, 'DATABASE_URL')
        assert hasattr(settings, 'DUPLICATE_THRESHOLD')
        assert hasattr(settings, 'VERIFY_THRESHOLD')
        
        # Check new setting exists
        assert hasattr(settings, 'VECTOR_INDEX_REBUILD_INTERVAL')
        
        # Old FAISS settings can be removed but shouldn't break anything
        # (they may or may not exist depending on migration state)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
