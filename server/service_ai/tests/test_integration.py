"""
Basic integration tests for ScyllaDB and MinIO
Run with: pytest tests/test_integration.py -v
"""
import pytest
import uuid
from datetime import datetime

from app.services.scylladb_manager import ScyllaDBManager
from app.services.minio_manager import MinIOManager


class TestScyllaDBIntegration:
    """Test ScyllaDB manager functionality"""
    
    @pytest.fixture(autouse=True)
    def setup(self):
        """Setup test environment"""
        try:
            self.scylladb = ScyllaDBManager()
        except Exception as e:
            pytest.skip(f"ScyllaDB not available: {e}")
    
    def test_save_verification_state(self):
        """Test saving verification state"""
        verification_id = uuid.uuid4()
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        company_id = uuid.uuid4()
        
        # Save verification state
        self.scylladb.save_verification_state(
            verification_id=verification_id,
            company_id=company_id,
            user_id=user_id,
            profile_id=profile_id,
            device_id="test_device",
            status="match",
            verified=True,
            similarity_score=0.95,
            liveness_score=0.88,
            metadata={"test": "data"},
            image_path="test/path.jpg"
        )
        
        # Retrieve verification state
        state = self.scylladb.get_verification_state(verification_id)
        
        assert state is not None
        assert state['verification_id'] == verification_id
        assert state['user_id'] == user_id
        assert state['profile_id'] == profile_id
        assert state['status'] == 'match'
        assert state['verified'] is True
        assert abs(state['similarity_score'] - 0.95) < 0.01
    
    def test_save_enrollment_state(self):
        """Test saving enrollment state"""
        profile_id = uuid.uuid4()
        user_id = uuid.uuid4()
        company_id = uuid.uuid4()
        
        # Save enrollment state
        self.scylladb.save_enrollment_state(
            profile_id=profile_id,
            company_id=company_id,
            user_id=user_id,
            device_id="test_device",
            status="ok",
            quality_score=0.92,
            metadata={"name": "Test User"},
            image_path="faces/test.jpg"
        )
        
        # Note: No direct getter for enrollment state in the current implementation
        # This test just ensures no errors are raised
    
    def test_get_user_verifications(self):
        """Test retrieving user verification history"""
        user_id = uuid.uuid4()
        company_id = uuid.uuid4()
        
        # Save a few verifications for the user
        for i in range(3):
            verification_id = uuid.uuid4()
            self.scylladb.save_verification_state(
                verification_id=verification_id,
                company_id=company_id,
                user_id=user_id,
                profile_id=uuid.uuid4(),
                device_id=f"device_{i}",
                status="match",
                verified=True,
                similarity_score=0.90 + i * 0.01
            )
        
        # Get verification history
        verifications = self.scylladb.get_user_verifications(user_id, company_id, limit=10)
        
        assert len(verifications) >= 3
        # Verify they are ordered by timestamp (descending)
        timestamps = [v['timestamp'] for v in verifications]
        assert timestamps == sorted(timestamps, reverse=True)


class TestMinIOIntegration:
    """Test MinIO manager functionality"""
    
    @pytest.fixture(autouse=True)
    def setup(self):
        """Setup test environment"""
        try:
            self.minio = MinIOManager()
        except Exception as e:
            pytest.skip(f"MinIO not available: {e}")
    
    def test_upload_and_retrieve_face_image(self):
        """Test uploading and retrieving face image"""
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        
        # Create test image data (simple binary data)
        test_image = b'\x89PNG\r\n\x1a\n' + b'\x00' * 100
        
        # Upload image
        object_name = self.minio.upload_face_image(
            image_data=test_image,
            user_id=user_id,
            profile_id=profile_id,
            file_extension="png"
        )
        
        assert object_name is not None
        assert f"{user_id}" in object_name
        assert f"{profile_id}" in object_name
        
        # Retrieve image
        retrieved_data = self.minio.get_face_image(object_name)
        
        assert retrieved_data is not None
        assert retrieved_data == test_image
        
        # Cleanup
        self.minio.delete_face_image(object_name)
    
    def test_upload_and_retrieve_verification_image(self):
        """Test uploading and retrieving verification image"""
        verification_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
        # Create test image data
        test_image = b'\xFF\xD8\xFF\xE0' + b'\x00' * 100  # JPEG header
        
        # Upload image
        object_name = self.minio.upload_verification_image(
            image_data=test_image,
            verification_id=verification_id,
            user_id=user_id,
            file_extension="jpg"
        )
        
        assert object_name is not None
        assert f"{verification_id}" in object_name
        
        # Retrieve image
        retrieved_data = self.minio.get_verification_image(object_name)
        
        assert retrieved_data is not None
        assert retrieved_data == test_image
        
        # Cleanup
        self.minio.delete_verification_image(object_name)
    
    def test_presigned_url_generation(self):
        """Test presigned URL generation"""
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        object_name = f"faces/{user_id}/{profile_id}.jpg"
        
        # Upload test image first
        test_image = b'\xFF\xD8\xFF\xE0' + b'\x00' * 50
        uploaded_path = self.minio.upload_face_image(
            image_data=test_image,
            user_id=user_id,
            profile_id=profile_id
        )
        
        # Generate presigned URL
        url = self.minio.get_presigned_url(
            bucket="face-images",
            object_name=uploaded_path,
            expiry=60
        )
        
        assert url is not None
        assert "face-images" in url
        
        # Cleanup
        self.minio.delete_face_image(uploaded_path)
    
    def test_delete_nonexistent_image(self):
        """Test deleting a non-existent image (should not raise error)"""
        # This should handle gracefully
        result = self.minio.delete_face_image("nonexistent/path.jpg")
        # MinIO may return False or True depending on implementation
        assert isinstance(result, bool)


class TestEndToEndIntegration:
    """Test end-to-end integration"""
    
    @pytest.fixture(autouse=True)
    def setup(self):
        """Setup test environment"""
        try:
            self.scylladb = ScyllaDBManager()
            self.minio = MinIOManager()
        except Exception as e:
            pytest.skip(f"ScyllaDB or MinIO not available: {e}")
    
    def test_enrollment_workflow(self):
        """Test complete enrollment workflow"""
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        company_id = uuid.uuid4()
        
        # 1. Upload face image to MinIO
        test_image = b'\xFF\xD8\xFF\xE0' + b'\x00' * 200
        image_path = self.minio.upload_face_image(
            image_data=test_image,
            user_id=user_id,
            profile_id=profile_id
        )
        
        assert image_path is not None
        
        # 2. Save enrollment state to ScyllaDB
        self.scylladb.save_enrollment_state(
            profile_id=profile_id,
            company_id=company_id,
            user_id=user_id,
            device_id="test_device",
            status="ok",
            quality_score=0.95,
            metadata={"name": "Test User"},
            image_path=image_path
        )
        
        # 3. Verify image can be retrieved
        retrieved_image = self.minio.get_face_image(image_path)
        assert retrieved_image == test_image
        
        # Cleanup
        self.minio.delete_face_image(image_path)
    
    def test_verification_workflow(self):
        """Test complete verification workflow"""
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        verification_id = uuid.uuid4()
        company_id = uuid.uuid4()
        
        # 1. Upload verification image to MinIO
        test_image = b'\xFF\xD8\xFF\xE0' + b'\x00' * 150
        image_path = self.minio.upload_verification_image(
            image_data=test_image,
            verification_id=verification_id,
            user_id=user_id
        )
        
        assert image_path is not None
        
        # 2. Save verification state to ScyllaDB
        self.scylladb.save_verification_state(
            verification_id=verification_id,
            company_id=company_id,
            user_id=user_id,
            profile_id=profile_id,
            device_id="test_device",
            status="match",
            verified=True,
            similarity_score=0.92,
            liveness_score=0.85,
            metadata={},
            image_path=image_path
        )
        
        # 3. Retrieve verification state
        state = self.scylladb.get_verification_state(verification_id)
        assert state is not None
        assert state['verified'] is True
        assert state['image_path'] == image_path
        
        # 4. Retrieve verification image
        retrieved_image = self.minio.get_verification_image(image_path)
        assert retrieved_image == test_image
        
        # 5. Check user verification history
        history = self.scylladb.get_user_verifications(user_id, company_id, limit=10)
        assert len(history) > 0
        assert any(v['verification_id'] == verification_id for v in history)
        
        # Cleanup
        self.minio.delete_verification_image(image_path)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
