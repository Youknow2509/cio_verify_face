"""
Face embedding generation using ArcFace
"""
import logging
import numpy as np
from typing import Optional

from app.core.config import settings
from app.core.face_detector import FaceDetector

logger = logging.getLogger(__name__)


class FaceEmbedding:
    """Generate face embeddings using ArcFace from InsightFace"""
    
    def __init__(self, detector: FaceDetector):
        """
        Initialize embedding generator
        
        Args:
            detector: FaceDetector instance
        """
        self.detector = detector
        self.embedding_version = settings.FACE_EMBEDDING_MODEL
        logger.info(f"Embedding generator initialized with version: {self.embedding_version}")
    
    def get_embedding(self, image: np.ndarray) -> Optional[np.ndarray]:
        """
        Get face embedding from image
        
        Args:
            image: Input image as numpy array (BGR format)
            
        Returns:
            Embedding vector (512-dim) or None if no face detected
        """
        try:
            # Detect faces
            faces = self.detector.detect_faces(image)
            
            if not faces:
                logger.warning("No faces detected in image")
                return None
            
            # Select best face
            best_face = self.detector.select_best_face(faces)
            
            if best_face is None or best_face['embedding'] is None:
                logger.warning("Could not extract embedding")
                return None
            
            # Normalize embedding
            embedding = best_face['embedding']
            embedding = embedding / np.linalg.norm(embedding)
            
            return embedding.astype(np.float32)
            
        except Exception as e:
            logger.error(f"Error generating embedding: {e}")
            return None
    
    def get_embedding_with_quality(self, image: np.ndarray) -> tuple:
        """
        Get face embedding with quality score
        
        Args:
            image: Input image as numpy array
            
        Returns:
            Tuple of (embedding, quality_score, face_dict)
        """
        try:
            # Detect faces
            faces = self.detector.detect_faces(image)
            
            if not faces:
                return None, 0.0, None
            
            # Select best face
            best_face = self.detector.select_best_face(faces)
            
            if best_face is None or best_face['embedding'] is None:
                return None, 0.0, None
            
            # Get quality score
            quality = self.detector.assess_quality(image, best_face)
            
            # Normalize embedding
            embedding = best_face['embedding']
            embedding = embedding / np.linalg.norm(embedding)
            
            return embedding.astype(np.float32), quality, best_face
            
        except Exception as e:
            logger.error(f"Error generating embedding with quality: {e}")
            return None, 0.0, None
    
    def compare_embeddings(self, emb1: np.ndarray, emb2: np.ndarray) -> float:
        """
        Compare two embeddings using cosine similarity
        
        Args:
            emb1: First embedding
            emb2: Second embedding
            
        Returns:
            Cosine similarity score (0-1)
        """
        try:
            # Ensure embeddings are normalized
            emb1_norm = emb1 / np.linalg.norm(emb1)
            emb2_norm = emb2 / np.linalg.norm(emb2)
            
            # Calculate cosine similarity
            similarity = np.dot(emb1_norm, emb2_norm)
            
            # Clip to [0, 1] range
            similarity = np.clip(similarity, 0.0, 1.0)
            
            return float(similarity)
            
        except Exception as e:
            logger.error(f"Error comparing embeddings: {e}")
            return 0.0
