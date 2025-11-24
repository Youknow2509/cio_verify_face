"""
MinIO manager for storing face images
"""
import logging
import io
from typing import Optional
from uuid import UUID
from datetime import timedelta
from minio import Minio
from minio.error import S3Error

from app.core.config import settings

logger = logging.getLogger(__name__)


class MinIOManager:
    """Manager for MinIO object storage operations"""
    
    def __init__(self):
        """Initialize MinIO client"""
        self.client = None
        self._connect()
    
    def _connect(self):
        """Connect to MinIO"""
        try:
            logger.info(f"Connecting to MinIO at {settings.MINIO_ENDPOINT}")
            
            self.client = Minio(
                settings.MINIO_ENDPOINT,
                access_key=settings.MINIO_ACCESS_KEY,
                secret_key=settings.MINIO_SECRET_KEY,
                secure=settings.MINIO_SECURE
            )
            
            # Create buckets if they don't exist
            self._create_buckets()
            
            logger.info("MinIO connection established successfully")
            
        except Exception as e:
            logger.error(f"Error connecting to MinIO: {e}")
            raise
    
    def _create_buckets(self):
        """Create required buckets if they don't exist"""
        try:
            buckets = [
                settings.MINIO_BUCKET_FACES,
                settings.MINIO_BUCKET_VERIFICATIONS
            ]
            
            for bucket in buckets:
                if not self.client.bucket_exists(bucket):
                    self.client.make_bucket(bucket)
                    logger.info(f"Created bucket: {bucket}")
                else:
                    logger.debug(f"Bucket already exists: {bucket}")
                    
        except S3Error as e:
            logger.error(f"Error creating buckets: {e}")
            raise
    
    def upload_face_image(
        self,
        image_data: bytes,
        user_id: UUID,
        profile_id: UUID,
        file_extension: str = "jpg"
    ) -> Optional[str]:
        """
        Upload face enrollment image to MinIO
        
        Args:
            image_data: Image binary data
            user_id: User ID
            profile_id: Profile ID
            file_extension: File extension (default: jpg)
            
        Returns:
            Object path in MinIO or None if failed
        """
        try:
            # Create object name: faces/{user_id}/{profile_id}.{ext}
            object_name = f"faces/{user_id}/{profile_id}.{file_extension}"
            
            # Upload to MinIO
            self.client.put_object(
                settings.MINIO_BUCKET_FACES,
                object_name,
                io.BytesIO(image_data),
                length=len(image_data),
                content_type=f"image/{file_extension}"
            )
            
            logger.info(f"Uploaded face image: {object_name}")
            return object_name
            
        except S3Error as e:
            logger.error(f"Error uploading face image: {e}")
            return None
    
    def upload_verification_image(
        self,
        image_data: bytes,
        verification_id: UUID,
        user_id: Optional[UUID] = None,
        file_extension: str = "jpg"
    ) -> Optional[str]:
        """
        Upload verification image to MinIO
        
        Args:
            image_data: Image binary data
            verification_id: Verification ID
            user_id: User ID (if matched)
            file_extension: File extension (default: jpg)
            
        Returns:
            Object path in MinIO or None if failed
        """
        try:
            # Create object name: verifications/{date}/{user_id or 'unknown'}/{verification_id}.{ext}
            from datetime import datetime
            date_str = datetime.utcnow().strftime("%Y-%m-%d")
            user_part = str(user_id) if user_id else "unknown"
            object_name = f"verifications/{date_str}/{user_part}/{verification_id}.{file_extension}"
            
            # Upload to MinIO
            self.client.put_object(
                settings.MINIO_BUCKET_VERIFICATIONS,
                object_name,
                io.BytesIO(image_data),
                length=len(image_data),
                content_type=f"image/{file_extension}"
            )
            
            logger.info(f"Uploaded verification image: {object_name}")
            return object_name
            
        except S3Error as e:
            logger.error(f"Error uploading verification image: {e}")
            return None
    
    def get_face_image(self, object_name: str) -> Optional[bytes]:
        """
        Download face image from MinIO
        
        Args:
            object_name: Object path in MinIO
            
        Returns:
            Image binary data or None if failed
        """
        try:
            response = self.client.get_object(
                settings.MINIO_BUCKET_FACES,
                object_name
            )
            
            data = response.read()
            response.close()
            response.release_conn()
            
            return data
            
        except S3Error as e:
            logger.error(f"Error downloading face image: {e}")
            return None
    
    def get_verification_image(self, object_name: str) -> Optional[bytes]:
        """
        Download verification image from MinIO
        
        Args:
            object_name: Object path in MinIO
            
        Returns:
            Image binary data or None if failed
        """
        try:
            response = self.client.get_object(
                settings.MINIO_BUCKET_VERIFICATIONS,
                object_name
            )
            
            data = response.read()
            response.close()
            response.release_conn()
            
            return data
            
        except S3Error as e:
            logger.error(f"Error downloading verification image: {e}")
            return None
    
    def get_presigned_url(
        self,
        bucket: str,
        object_name: str,
        expiry: int = 3600
    ) -> Optional[str]:
        """
        Get presigned URL for temporary access to an object
        
        Args:
            bucket: Bucket name
            object_name: Object path
            expiry: URL expiry time in seconds (default: 1 hour)
            
        Returns:
            Presigned URL or None if failed
        """
        try:
            url = self.client.presigned_get_object(
                bucket,
                object_name,
                expires=timedelta(seconds=expiry)
            )
            return url
            
        except S3Error as e:
            logger.error(f"Error generating presigned URL: {e}")
            return None
    
    def delete_face_image(self, object_name: str) -> bool:
        """
        Delete face image from MinIO
        
        Args:
            object_name: Object path in MinIO
            
        Returns:
            True if successful, False otherwise
        """
        try:
            self.client.remove_object(
                settings.MINIO_BUCKET_FACES,
                object_name
            )
            logger.info(f"Deleted face image: {object_name}")
            return True
            
        except S3Error as e:
            logger.error(f"Error deleting face image: {e}")
            return False
    
    def delete_verification_image(self, object_name: str) -> bool:
        """
        Delete verification image from MinIO
        
        Args:
            object_name: Object path in MinIO
            
        Returns:
            True if successful, False otherwise
        """
        try:
            self.client.remove_object(
                settings.MINIO_BUCKET_VERIFICATIONS,
                object_name
            )
            logger.info(f"Deleted verification image: {object_name}")
            return True
            
        except S3Error as e:
            logger.error(f"Error deleting verification image: {e}")
            return False
