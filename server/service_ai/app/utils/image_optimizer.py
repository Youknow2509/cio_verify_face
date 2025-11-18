"""
Image optimization utilities for bandwidth and storage optimization
"""
import cv2
import numpy as np
from typing import Tuple, Optional
import logging

from app.core.config import settings

logger = logging.getLogger(__name__)


class ImageOptimizer:
    """Optimizes images for transmission and storage"""
    
    @staticmethod
    def resize_image(image: np.ndarray, max_size: int = None) -> np.ndarray:
        """
        Resize image to maximum dimensions while maintaining aspect ratio
        
        Args:
            image: Input image as numpy array
            max_size: Maximum width or height (uses settings if None)
            
        Returns:
            Resized image
        """
        if max_size is None:
            max_size = settings.IMAGE_MAX_SIZE
        
        height, width = image.shape[:2]
        
        # Skip if already smaller
        if height <= max_size and width <= max_size:
            return image
        
        # Calculate new dimensions
        if height > width:
            new_height = max_size
            new_width = int(width * (max_size / height))
        else:
            new_width = max_size
            new_height = int(height * (max_size / width))
        
        resized = cv2.resize(image, (new_width, new_height), interpolation=cv2.INTER_AREA)
        logger.debug(f"Resized image from {width}x{height} to {new_width}x{new_height}")
        
        return resized
    
    @staticmethod
    def compress_image(
        image: np.ndarray,
        quality: int = None,
        format: str = 'jpg'
    ) -> Tuple[bytes, int]:
        """
        Compress image with specified quality
        
        Args:
            image: Input image as numpy array
            quality: JPEG quality 1-100 (uses settings if None)
            format: Image format ('jpg' or 'png')
            
        Returns:
            Tuple of (compressed bytes, original size reduction percentage)
        """
        if quality is None:
            quality = settings.IMAGE_QUALITY
        
        # Calculate original size
        original_size = image.nbytes
        
        # Encode with compression
        if format == 'jpg':
            encode_params = [cv2.IMWRITE_JPEG_QUALITY, quality]
            success, buffer = cv2.imencode('.jpg', image, encode_params)
        else:  # png
            encode_params = [cv2.IMWRITE_PNG_COMPRESSION, 9]
            success, buffer = cv2.imencode('.png', image, encode_params)
        
        if not success:
            raise ValueError("Failed to encode image")
        
        compressed_bytes = buffer.tobytes()
        compressed_size = len(compressed_bytes)
        reduction = int((1 - compressed_size / original_size) * 100)
        
        logger.debug(f"Compressed image: {original_size} -> {compressed_size} bytes ({reduction}% reduction)")
        
        return compressed_bytes, reduction
    
    @staticmethod
    def optimize_for_storage(image: np.ndarray) -> bytes:
        """
        Optimize image for storage (resize + compress)
        
        Args:
            image: Input image as numpy array
            
        Returns:
            Optimized image as bytes
        """
        # Resize if needed
        resized = ImageOptimizer.resize_image(image)
        
        # Compress
        compressed, _ = ImageOptimizer.compress_image(resized)
        
        return compressed
    
    @staticmethod
    def estimate_base64_overhead(image_bytes: bytes) -> dict:
        """
        Estimate bandwidth overhead of base64 encoding
        
        Args:
            image_bytes: Original image bytes
            
        Returns:
            Dictionary with size information
        """
        import base64
        
        original_size = len(image_bytes)
        base64_encoded = base64.b64encode(image_bytes)
        base64_size = len(base64_encoded)
        overhead = base64_size - original_size
        overhead_percent = (overhead / original_size) * 100
        
        return {
            'original_size': original_size,
            'base64_size': base64_size,
            'overhead_bytes': overhead,
            'overhead_percent': round(overhead_percent, 2)
        }
