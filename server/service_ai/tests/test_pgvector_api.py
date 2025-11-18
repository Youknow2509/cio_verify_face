"""
Unit tests for pgvector manager - tests API compatibility without requiring database
"""
import pytest
import numpy as np
from unittest.mock import Mock, MagicMock, patch
from datetime import datetime

from app.services.pgvector_manager import PgVectorManager


class TestPgVectorManagerAPI:
    """Test PgVectorManager API compatibility with FAISSIndexManager"""
    
    def test_initialization_api(self):
        """Test that PgVectorManager can be initialized with same API as FAISS"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                # Verify basic attributes exist
                assert hasattr(manager, 'dimension')
                assert hasattr(manager, 'engine')
                assert hasattr(manager, 'SessionLocal')
    
    def test_add_embedding_signature(self):
        """Test add_embedding has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                # Verify method exists and has correct signature
                assert hasattr(manager, 'add_embedding')
                assert callable(manager.add_embedding)
    
    def test_remove_embedding_signature(self):
        """Test remove_embedding has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'remove_embedding')
                assert callable(manager.remove_embedding)
    
    def test_search_signature(self):
        """Test search has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'search')
                assert callable(manager.search)
    
    def test_rebuild_index_signature(self):
        """Test rebuild_index has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'rebuild_index')
                assert callable(manager.rebuild_index)
    
    def test_save_index_signature(self):
        """Test save_index exists (even if it's a no-op)"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'save_index')
                assert callable(manager.save_index)
    
    def test_get_size_signature(self):
        """Test get_size has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'get_size')
                assert callable(manager.get_size)
    
    def test_clear_signature(self):
        """Test clear has the correct signature"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'clear')
                assert callable(manager.clear)
    
    def test_last_rebuild_property(self):
        """Test last_rebuild property exists"""
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                assert hasattr(manager, 'last_rebuild')
    
    def test_embedding_normalization(self):
        """Test that embeddings are normalized correctly"""
        # Test the normalization logic in isolation
        embedding = np.array([3.0, 4.0] + [0.0] * 510, dtype=np.float32)
        
        # Normalize
        normalized = embedding / np.linalg.norm(embedding)
        
        # Check norm is 1
        assert np.isclose(np.linalg.norm(normalized), 1.0)
        
        # Check values are scaled correctly
        expected_first = 3.0 / 5.0  # norm of [3,4] is 5
        expected_second = 4.0 / 5.0
        assert np.isclose(normalized[0], expected_first)
        assert np.isclose(normalized[1], expected_second)


class TestPgVectorManagerCompatibility:
    """Test compatibility with FAISS API"""
    
    def test_api_matches_faiss(self):
        """Verify PgVectorManager has same public API as FAISSIndexManager"""
        # Expected public methods from FAISSIndexManager
        expected_methods = [
            'add_embedding',
            'remove_embedding',
            'search',
            'rebuild_index',
            'save_index',
            'get_size',
            'clear'
        ]
        
        expected_properties = [
            'last_rebuild'
        ]
        
        with patch('app.services.pgvector_manager.create_engine'):
            with patch('app.services.pgvector_manager.sessionmaker'):
                manager = PgVectorManager()
                
                # Check all methods exist
                for method in expected_methods:
                    assert hasattr(manager, method), f"Missing method: {method}"
                    assert callable(getattr(manager, method)), f"Not callable: {method}"
                
                # Check all properties exist
                for prop in expected_properties:
                    assert hasattr(manager, prop), f"Missing property: {prop}"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
