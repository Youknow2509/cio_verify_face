"""
Example usage of the enhanced face verification service with ScyllaDB and MinIO

This script demonstrates:
1. Face enrollment with image storage in MinIO
2. Face verification with state tracking in ScyllaDB
3. Querying verification history from ScyllaDB
"""
import asyncio
import base64
import uuid
from datetime import datetime
import cv2
import numpy as np

from app.services.face_service import FaceService
from app.services.scylladb_manager import ScyllaDBManager
from app.services.minio_manager import MinIOManager


async def example_enrollment():
    """Example: Enroll a new face with MinIO storage"""
    print("\n=== Example 1: Face Enrollment with MinIO Storage ===")
    
    # Initialize service
    face_service = FaceService()
    
    # Create sample IDs
    user_id = uuid.uuid4()
    company_id = uuid.uuid4()
    print(f"User ID: {user_id}")
    print(f"Company ID: {company_id}")
    
    # Load a sample image (replace with actual image path)
    image_path = "./data/sample_face.jpg"
    try:
        image = cv2.imread(image_path)
        if image is None:
            print(f"Warning: Could not load sample image from {image_path}")
            print("Please provide a valid face image to test enrollment")
            return None
        
        # Encode image to base64
        _, buffer = cv2.imencode('.jpg', image)
        image_base64 = base64.b64encode(buffer).decode('utf-8')
        
        # Enroll face
        result = await face_service.enroll_face(
            user_id=user_id,
            company_id=company_id,
            image_base64=image_base64,
            device_id="device_001",
            make_primary=True,
            metadata={"name": "John Doe", "department": "Engineering"}
        )
        
        print(f"Enrollment status: {result.status}")
        print(f"Profile ID: {result.profile_id}")
        print(f"Quality score: {result.quality_score}")
        print(f"Message: {result.message}")
        
        # The face image is now stored in MinIO at:
        # minio://face-images/faces/{user_id}/{profile_id}.jpg
        
        # The enrollment state is recorded in ScyllaDB enrollment_states table
        
        return company_id, user_id, result.profile_id
        
    except FileNotFoundError:
        print(f"Sample image not found at {image_path}")
        print("Please create a sample image or update the path")
        return None


async def example_verification(company_id: uuid.UUID):
    """Example: Verify a face with ScyllaDB state tracking"""
    print("\n=== Example 2: Face Verification with ScyllaDB Tracking ===")
    
    # Initialize service
    face_service = FaceService()
    
    # Load a sample image for verification
    image_path = "./data/sample_face.jpg"
    try:
        image = cv2.imread(image_path)
        if image is None:
            print(f"Warning: Could not load sample image from {image_path}")
            return
        
        # Encode image to base64
        _, buffer = cv2.imencode('.jpg', image)
        image_base64 = base64.b64encode(buffer).decode('utf-8')
        
        # Verify face (1:N mode - search across all users)
        result = await face_service.verify_face(
            image_base64=image_base64,
            company_id=company_id,
            device_id="device_002",
            search_mode="1:N",
            top_k=3
        )
        
        print(f"Verification status: {result.status}")
        print(f"Verified: {result.verified}")
        print(f"Number of matches: {len(result.matches)}")
        
        if result.best_match:
            print(f"\nBest match:")
            print(f"  User ID: {result.best_match.user_id}")
            print(f"  Profile ID: {result.best_match.profile_id}")
            print(f"  Similarity: {result.best_match.similarity:.4f}")
            print(f"  Is primary: {result.best_match.is_primary}")
        
        if result.liveness_score:
            print(f"Liveness score: {result.liveness_score:.4f}")
        
        # If verified, the verification image is stored in MinIO at:
        # minio://verification-images/verifications/{date}/{user_id}/{verification_id}.jpg
        
        # The verification state is recorded in ScyllaDB authentication_states table
        # and indexed in user_verifications table for quick history lookup
        
    except FileNotFoundError:
        print(f"Sample image not found at {image_path}")


async def example_query_verification_history(company_id: uuid.UUID):
    """Example: Query verification history from ScyllaDB"""
    print("\n=== Example 3: Query Verification History from ScyllaDB ===")
    
    try:
        # Initialize ScyllaDB manager
        scylladb = ScyllaDBManager()
        
        # Create a sample user ID (in real scenario, use actual user ID)
        user_id = uuid.uuid4()
        print(f"Querying verification history for user: {user_id}")
        
        # Get recent verifications for user
        verifications = scylladb.get_user_verifications(user_id, company_id, limit=10)
        
        if verifications:
            print(f"\nFound {len(verifications)} recent verifications:")
            for v in verifications:
                print(f"\n  Verification ID: {v['verification_id']}")
                print(f"  Timestamp: {v['timestamp']}")
                print(f"  Verified: {v['verified']}")
                print(f"  Similarity: {v['similarity_score']:.4f}")
                print(f"  Device: {v['device_id']}")
        else:
            print("No verification history found for this user")
        
        # Query a specific verification state
        verification_id = uuid.uuid4()
        state = scylladb.get_verification_state(verification_id)
        
        if state:
            print(f"\nVerification state for {verification_id}:")
            print(f"  Status: {state['status']}")
            print(f"  Verified: {state['verified']}")
            print(f"  User ID: {state['user_id']}")
            print(f"  Image path: {state['image_path']}")
        else:
            print(f"No state found for verification {verification_id}")
            
    except Exception as e:
        print(f"Error querying ScyllaDB: {e}")
        print("Make sure ScyllaDB is running and accessible")


async def example_retrieve_image_from_minio():
    """Example: Retrieve stored face image from MinIO"""
    print("\n=== Example 4: Retrieve Face Image from MinIO ===")
    
    try:
        # Initialize MinIO manager
        minio = MinIOManager()
        
        # Example object path (replace with actual path from enrollment)
        user_id = uuid.uuid4()
        profile_id = uuid.uuid4()
        object_name = f"faces/{user_id}/{profile_id}.jpg"
        
        print(f"Attempting to retrieve: {object_name}")
        
        # Get image data
        image_data = minio.get_face_image(object_name)
        
        if image_data:
            print(f"Successfully retrieved image ({len(image_data)} bytes)")
            
            # Convert bytes to OpenCV image
            nparr = np.frombuffer(image_data, np.uint8)
            image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
            
            if image is not None:
                print(f"Image shape: {image.shape}")
                # You can now use the image for processing
        else:
            print("Image not found (this is expected if you haven't enrolled yet)")
        
        # Generate presigned URL for temporary access
        url = minio.get_presigned_url(
            bucket="face-images",
            object_name=object_name,
            expiry=3600  # 1 hour
        )
        
        if url:
            print(f"\nPresigned URL (valid for 1 hour):")
            print(url)
            
    except Exception as e:
        print(f"Error accessing MinIO: {e}")
        print("Make sure MinIO is running and accessible")


async def main():
    """Run all examples"""
    print("=" * 70)
    print("Face Verification Service - ScyllaDB & MinIO Integration Examples")
    print("=" * 70)
    
    # Run examples
    enrolled = await example_enrollment()
    
    if enrolled:
        company_id, user_id, profile_id = enrolled
        await example_verification(company_id)
        await example_query_verification_history(company_id)
    await example_retrieve_image_from_minio()
    
    print("\n" + "=" * 70)
    print("Examples completed!")
    print("=" * 70)
    print("\nNote: Some examples may show 'not found' messages if you haven't")
    print("enrolled faces yet. This is normal and demonstrates error handling.")


if __name__ == "__main__":
    asyncio.run(main())
