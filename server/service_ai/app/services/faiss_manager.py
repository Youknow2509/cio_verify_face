"""
FAISS index manager for efficient face embedding search
"""
import logging
import numpy as np
import faiss
import os
import pickle
from typing import List, Tuple, Optional, Dict
from datetime import datetime
import threading

from app.core.config import settings

logger = logging.getLogger(__name__)


class FAISSIndexManager:
    """Manage FAISS index for face embeddings"""
    
    def __init__(self):
        """Initialize FAISS index manager"""
        self.dimension = settings.EMBEDDING_DIMENSION
        self.index = None
        self.profile_ids = []  # Map index position to profile_id
        self.user_ids = []  # Map index position to user_id
        self.is_primary = []  # Map index position to is_primary flag
        self.index_file = os.path.join(settings.INDEX_PATH, "faiss_index.bin")
        self.metadata_file = os.path.join(settings.INDEX_PATH, "index_metadata.pkl")
        self.lock = threading.Lock()
        self.last_rebuild = None
        
        self._initialize_index()
    
    def _initialize_index(self):
        """Initialize or load FAISS index"""
        try:
            if os.path.exists(self.index_file) and os.path.exists(self.metadata_file):
                logger.info("Loading existing FAISS index...")
                self._load_index()
            else:
                logger.info("Creating new FAISS index...")
                self._create_index()
            
            logger.info(f"FAISS index initialized with {self.get_size()} embeddings")
        except Exception as e:
            logger.error(f"Error initializing FAISS index: {e}")
            self._create_index()
    
    def _create_index(self):
        """Create a new FAISS index"""
        with self.lock:
            if settings.FAISS_INDEX_TYPE == "Flat":
                # Flat index for exact search (best for < 10k embeddings)
                self.index = faiss.IndexFlatIP(self.dimension)  # Inner product (cosine similarity)
            elif settings.FAISS_INDEX_TYPE == "IVF":
                # IVF index for faster search (good for > 10k embeddings)
                quantizer = faiss.IndexFlatIP(self.dimension)
                self.index = faiss.IndexIVFFlat(quantizer, self.dimension, 100)
            else:
                # Default to Flat
                self.index = faiss.IndexFlatIP(self.dimension)
            
            self.profile_ids = []
            self.user_ids = []
            self.is_primary = []
            self.last_rebuild = datetime.utcnow()
    
    def _load_index(self):
        """Load FAISS index from disk"""
        with self.lock:
            try:
                self.index = faiss.read_index(self.index_file)
                
                with open(self.metadata_file, 'rb') as f:
                    metadata = pickle.load(f)
                    self.profile_ids = metadata['profile_ids']
                    self.user_ids = metadata['user_ids']
                    self.is_primary = metadata['is_primary']
                    self.last_rebuild = metadata.get('last_rebuild', datetime.utcnow())
                
                logger.info("FAISS index loaded successfully")
            except Exception as e:
                logger.error(f"Error loading index: {e}")
                raise
    
    def save_index(self):
        """Save FAISS index to disk"""
        with self.lock:
            try:
                faiss.write_index(self.index, self.index_file)
                
                metadata = {
                    'profile_ids': self.profile_ids,
                    'user_ids': self.user_ids,
                    'is_primary': self.is_primary,
                    'last_rebuild': self.last_rebuild
                }
                
                with open(self.metadata_file, 'wb') as f:
                    pickle.dump(metadata, f)
                
                logger.info("FAISS index saved successfully")
            except Exception as e:
                logger.error(f"Error saving index: {e}")
                raise
    
    def add_embedding(self, profile_id: str, user_id: str, embedding: np.ndarray, 
                     is_primary: bool = False):
        """
        Add a single embedding to the index
        
        Args:
            profile_id: Profile UUID as string
            user_id: User UUID as string
            embedding: Face embedding vector
            is_primary: Whether this is the primary profile
        """
        with self.lock:
            try:
                # Ensure embedding is normalized
                embedding = embedding.reshape(1, -1).astype(np.float32)
                embedding = embedding / np.linalg.norm(embedding)
                
                # Add to index
                self.index.add(embedding)
                self.profile_ids.append(profile_id)
                self.user_ids.append(user_id)
                self.is_primary.append(is_primary)
                
                logger.debug(f"Added embedding for profile {profile_id}")
            except Exception as e:
                logger.error(f"Error adding embedding: {e}")
                raise
    
    def remove_embedding(self, profile_id: str):
        """
        Remove an embedding from the index
        
        Args:
            profile_id: Profile UUID as string
        """
        with self.lock:
            try:
                if profile_id not in self.profile_ids:
                    logger.warning(f"Profile {profile_id} not found in index")
                    return
                
                # Find index position
                idx = self.profile_ids.index(profile_id)
                
                # Remove from metadata
                self.profile_ids.pop(idx)
                self.user_ids.pop(idx)
                self.is_primary.pop(idx)
                
                # FAISS doesn't support direct removal, so we need to rebuild
                # For now, mark for rebuild
                logger.info(f"Removed embedding for profile {profile_id}, index rebuild recommended")
            except Exception as e:
                logger.error(f"Error removing embedding: {e}")
                raise
    
    def search(self, embedding: np.ndarray, k: int = 5) -> List[Dict]:
        """
        Search for similar embeddings
        
        Args:
            embedding: Query embedding
            k: Number of results to return
            
        Returns:
            List of matches with profile_id, user_id, similarity, is_primary
        """
        with self.lock:
            try:
                if self.get_size() == 0:
                    return []
                
                # Ensure embedding is normalized
                embedding = embedding.reshape(1, -1).astype(np.float32)
                embedding = embedding / np.linalg.norm(embedding)
                
                # Search
                k = min(k, self.get_size())
                distances, indices = self.index.search(embedding, k)
                
                # Format results
                results = []
                for dist, idx in zip(distances[0], indices[0]):
                    if idx >= 0 and idx < len(self.profile_ids):
                        results.append({
                            'profile_id': self.profile_ids[idx],
                            'user_id': self.user_ids[idx],
                            'similarity': float(dist),
                            'is_primary': self.is_primary[idx]
                        })
                
                return results
            except Exception as e:
                logger.error(f"Error searching index: {e}")
                raise
    
    def rebuild_index(self, embeddings: List[Tuple[str, str, np.ndarray, bool]]):
        """
        Rebuild the entire index
        
        Args:
            embeddings: List of (profile_id, user_id, embedding, is_primary) tuples
        """
        with self.lock:
            try:
                logger.info(f"Rebuilding index with {len(embeddings)} embeddings...")
                
                # Create new index
                self._create_index()
                
                # Add all embeddings
                for profile_id, user_id, embedding, is_primary in embeddings:
                    embedding = embedding.reshape(1, -1).astype(np.float32)
                    embedding = embedding / np.linalg.norm(embedding)
                    
                    self.index.add(embedding)
                    self.profile_ids.append(profile_id)
                    self.user_ids.append(user_id)
                    self.is_primary.append(is_primary)
                
                self.last_rebuild = datetime.utcnow()
                
                # Save to disk
                self.save_index()
                
                logger.info(f"Index rebuilt successfully with {len(embeddings)} embeddings")
            except Exception as e:
                logger.error(f"Error rebuilding index: {e}")
                raise
    
    def get_size(self) -> int:
        """Get the number of embeddings in the index"""
        return len(self.profile_ids)
    
    def clear(self):
        """Clear the index"""
        with self.lock:
            self._create_index()
            logger.info("Index cleared")
