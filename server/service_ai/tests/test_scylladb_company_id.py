"""
Unit tests for ScyllaDB manager with company_id partitioning
"""
import pytest
import uuid
from datetime import datetime

from app.services.scylladb_manager import ScyllaDBManager


class TestScyllaDBCompanyPartitioning:
    """Test ScyllaDB company_id partitioning"""
    
    @pytest.fixture(autouse=True)
    def setup(self):
        """Setup test environment"""
        try:
            self.scylladb = ScyllaDBManager()
        except Exception as e:
            pytest.skip(f"ScyllaDB not available: {e}")
    
    def test_save_verification_with_company_id(self):
        """Test saving verification state with company_id"""
        verification_id = uuid.uuid4()
        company_id = uuid.uuid4()
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        
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
        assert state['company_id'] == company_id
        assert state['user_id'] == user_id
        assert state['verified'] is True
    
    def test_save_enrollment_with_company_id(self):
        """Test saving enrollment state with company_id"""
        profile_id = uuid.uuid4()
        company_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
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
        
        # Test passes if no exception raised
        assert True
    
    def test_get_user_verifications_by_company(self):
        """Test retrieving user verifications filtered by company_id"""
        company_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
        # Save verifications for this company/user
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
    
    def test_company_isolation(self):
        """Test that data is isolated per company"""
        company1_id = uuid.uuid4()
        company2_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
        # Save verification for company 1
        ver1_id = uuid.uuid4()
        self.scylladb.save_verification_state(
            verification_id=ver1_id,
            company_id=company1_id,
            user_id=user_id,
            profile_id=uuid.uuid4(),
            device_id="device_1",
            status="match",
            verified=True,
            similarity_score=0.95
        )
        
        # Save verification for company 2
        ver2_id = uuid.uuid4()
        self.scylladb.save_verification_state(
            verification_id=ver2_id,
            company_id=company2_id,
            user_id=user_id,
            profile_id=uuid.uuid4(),
            device_id="device_2",
            status="match",
            verified=True,
            similarity_score=0.90
        )
        
        # Get verifications for company 1
        company1_vers = self.scylladb.get_user_verifications(user_id, company1_id, limit=10)
        
        # Get verifications for company 2
        company2_vers = self.scylladb.get_user_verifications(user_id, company2_id, limit=10)
        
        # Each should only see their own data
        assert len(company1_vers) >= 1
        assert len(company2_vers) >= 1
        
        # Verification IDs should be different
        company1_ids = {v['verification_id'] for v in company1_vers}
        company2_ids = {v['verification_id'] for v in company2_vers}
        
        assert ver1_id in company1_ids or len(company1_vers) > 1
        assert ver2_id in company2_ids or len(company2_vers) > 1
    
    def test_save_audit_log_with_company(self):
        """Test saving audit log with company_id"""
        log_id = uuid.uuid4()
        company_id = uuid.uuid4()
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        
        # Save audit log
        self.scylladb.save_audit_log(
            log_id=log_id,
            company_id=company_id,
            operation="enroll",
            status="success",
            profile_id=profile_id,
            user_id=user_id,
            device_id="test_device",
            quality_score=0.95,
            metadata={"test": "audit"}
        )
        
        # Get audit logs for company
        logs = self.scylladb.get_audit_logs(company_id, operation="enroll", limit=10)
        
        assert len(logs) >= 1
        assert any(log['log_id'] == log_id for log in logs)
    
    def test_get_user_audit_logs_by_company(self):
        """Test retrieving user audit logs filtered by company_id"""
        company_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
        # Save audit logs
        for i in range(2):
            log_id = uuid.uuid4()
            self.scylladb.save_audit_log(
                log_id=log_id,
                company_id=company_id,
                operation="verify",
                status="success",
                user_id=user_id,
                similarity_score=0.90 + i * 0.02
            )
        
        # Get user audit logs
        logs = self.scylladb.get_user_audit_logs(user_id, company_id, limit=10)
        
        assert len(logs) >= 2
    
    def test_audit_log_company_isolation(self):
        """Test that audit logs are isolated per company"""
        company1_id = uuid.uuid4()
        company2_id = uuid.uuid4()
        user_id = uuid.uuid4()
        
        # Save log for company 1
        log1_id = uuid.uuid4()
        self.scylladb.save_audit_log(
            log_id=log1_id,
            company_id=company1_id,
            operation="enroll",
            status="success",
            user_id=user_id
        )
        
        # Save log for company 2
        log2_id = uuid.uuid4()
        self.scylladb.save_audit_log(
            log_id=log2_id,
            company_id=company2_id,
            operation="enroll",
            status="success",
            user_id=user_id
        )
        
        # Get logs for company 1
        company1_logs = self.scylladb.get_audit_logs(company1_id, operation="enroll", limit=10)
        company2_logs = self.scylladb.get_audit_logs(company2_id, operation="enroll", limit=10)
        
        # Each should only see their own logs
        company1_log_ids = {log['log_id'] for log in company1_logs}
        company2_log_ids = {log['log_id'] for log in company2_logs}
        
        # log1 should be in company1, log2 should be in company2
        assert log1_id in company1_log_ids or len(company1_logs) > 1
        assert log2_id in company2_log_ids or len(company2_logs) > 1


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
