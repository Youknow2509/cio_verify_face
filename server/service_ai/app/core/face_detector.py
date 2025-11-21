"""
Face detection and alignment using InsightFace
"""
import logging
import numpy as np
import cv2
from typing import Optional, List, Tuple
import insightface
from insightface.app import FaceAnalysis

from app.core.config import settings

logger = logging.getLogger(__name__)


class FaceDetector:
    """Face detection and alignment using RetinaFace from InsightFace"""
    
    def __init__(self):
        """Initialize face detector"""
        self.app = None
        self._initialize_detector()
    
    def _initialize_detector(self):
        """Initialize InsightFace detector"""
        try:
            logger.info("Initializing face detector...")
            
            # Determine execution providers based on compute mode
            if settings.COMPUTE_MODE.lower() == 'gpu':
                providers = ['CUDAExecutionProvider', 'CPUExecutionProvider']
                logger.info("Using GPU mode (CUDA)")
            else:
                providers = ['CPUExecutionProvider']
                logger.info("Using CPU mode")
            
            self.app = FaceAnalysis(
                name='buffalo_l',  # Using buffalo_l model which includes RetinaFace
                root=settings.MODEL_PATH,
                providers=providers
            )
            self.app.prepare(ctx_id=0, det_size=(640, 640))
            logger.info("Face detector initialized successfully")
        except Exception as e:
            logger.error(f"Failed to initialize face detector: {e}")
            raise
    
    def detect_faces(self, image: np.ndarray) -> List[dict]:
        """
        Detect faces in image
        
        Args:
            image: Input image as numpy array (BGR format)
            
        Returns:
            List of detected faces with bounding boxes and landmarks
        """
        try:
            faces = self.app.get(image)
            
            results = []
            for face in faces:
                bbox = face.bbox.astype(int)
                face_width = bbox[2] - bbox[0]
                face_height = bbox[3] - bbox[1]
                
                # Filter small faces
                if face_width < settings.MIN_FACE_SIZE or face_height < settings.MIN_FACE_SIZE:
                    continue
                
                result = {
                    'bbox': bbox,
                    'landmarks': face.kps,
                    'det_score': float(face.det_score),
                    'embedding': face.embedding if hasattr(face, 'embedding') else None,
                    'face_width': face_width,
                    'face_height': face_height,
                }
                results.append(result)
            
            return results
        
        except Exception as e:
            logger.error(f"Error detecting faces: {e}")
            raise
    
    def select_best_face(self, faces: List[dict]) -> Optional[dict]:
        """
        Select the best face from detected faces
        
        Args:
            faces: List of detected faces
            
        Returns:
            Best face or None if no faces
        """
        if not faces:
            return None
        
        # Select face with highest detection score
        best_face = max(faces, key=lambda f: f['det_score'])
        return best_face
    
    def align_face(self, image: np.ndarray, landmarks: np.ndarray) -> np.ndarray:
        """
        Align face using facial landmarks
        
        Args:
            image: Input image
            landmarks: Facial landmarks (5 points)
            
        Returns:
            Aligned face image
        """
        try:
            # Reference points for alignment (112x112 standard)
            reference = np.array([
                [38.2946, 51.6963],
                [73.5318, 51.5014],
                [56.0252, 71.7366],
                [41.5493, 92.3655],
                [70.7299, 92.2041]
            ], dtype=np.float32)
            
            # Calculate transformation matrix
            tform = cv2.estimateAffinePartial2D(landmarks, reference)[0]
            
            # Apply transformation
            aligned = cv2.warpAffine(
                image, 
                tform, 
                (112, 112), 
                borderValue=0.0
            )
            
            return aligned
            
        except Exception as e:
            logger.error(f"Error aligning face: {e}")
            raise
    
    def assess_quality(self, image: np.ndarray, face: dict) -> float:
        """
        Assess face image quality
        
        Args:
            image: Input image
            face: Detected face dict
            
        Returns:
            Quality score (0-1)
        """
        try:
            bbox = face['bbox']
            face_img = image[bbox[1]:bbox[3], bbox[0]:bbox[2]]
            
            # Calculate blur score using Laplacian variance
            gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY)
            blur_score = cv2.Laplacian(gray, cv2.CV_64F).var()
            
            # Normalize blur score (higher is better)
            blur_quality = min(blur_score / 500.0, 1.0)
            
            # Calculate brightness
            brightness = np.mean(gray) / 255.0
            brightness_quality = 1.0 - abs(brightness - 0.5) * 2
            
            # Calculate size quality
            size_quality = min(face['face_width'] / 200.0, 1.0)
            
            # Combined quality score
            quality = (blur_quality * 0.4 + 
                      brightness_quality * 0.3 + 
                      size_quality * 0.3)
            
            return float(quality)
            
        except Exception as e:
            logger.error(f"Error assessing quality: {e}")
            return 0.5
