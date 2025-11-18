"""
Example usage of Face Verification Service

This script demonstrates how to use the face verification service API.
"""

import requests
import base64
import json
from pathlib import Path
from uuid import uuid4


class FaceVerificationClient:
    """Client for Face Verification Service"""
    
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.api_base = f"{base_url}/api/v1/face"
    
    def health_check(self):
        """Check service health"""
        response = requests.get(f"{self.base_url}/health")
        return response.json()
    
    def enroll_face(self, user_id, image_path, make_primary=True, device_id=None):
        """
        Enroll a face for a user
        
        Args:
            user_id: UUID of user
            image_path: Path to image file
            make_primary: Set as primary profile
            device_id: Optional device ID
        """
        # Read and encode image
        with open(image_path, "rb") as f:
            image_base64 = base64.b64encode(f.read()).decode()
        
        # Prepare request
        payload = {
            "user_id": str(user_id),
            "image_base64": image_base64,
            "make_primary": make_primary,
            "metadata": {
                "enrollment_source": "example_script",
                "image_filename": Path(image_path).name
            }
        }
        
        if device_id:
            payload["device_id"] = device_id
        
        # Send request
        response = requests.post(
            f"{self.api_base}/enroll",
            json=payload
        )
        
        return response.json()
    
    def verify_face(self, image_path, user_id=None, search_mode="1:N", top_k=5):
        """
        Verify a face
        
        Args:
            image_path: Path to image file
            user_id: Optional user ID for 1:1 verification
            search_mode: "1:1" or "1:N"
            top_k: Number of top matches
        """
        # Read and encode image
        with open(image_path, "rb") as f:
            image_base64 = base64.b64encode(f.read()).decode()
        
        # Prepare request
        payload = {
            "image_base64": image_base64,
            "search_mode": search_mode,
            "top_k": top_k
        }
        
        if user_id:
            payload["user_id"] = str(user_id)
        
        # Send request
        response = requests.post(
            f"{self.api_base}/verify",
            json=payload
        )
        
        return response.json()
    
    def get_user_profiles(self, user_id):
        """Get all face profiles for a user"""
        response = requests.get(
            f"{self.api_base}/profiles/{user_id}"
        )
        return response.json()
    
    def update_profile(self, profile_id, image_path=None, make_primary=None):
        """Update a face profile"""
        payload = {}
        
        if image_path:
            with open(image_path, "rb") as f:
                image_base64 = base64.b64encode(f.read()).decode()
            payload["image_base64"] = image_base64
        
        if make_primary is not None:
            payload["make_primary"] = make_primary
        
        response = requests.put(
            f"{self.api_base}/profile/{profile_id}",
            json=payload
        )
        
        return response.json()
    
    def delete_profile(self, profile_id, hard_delete=False):
        """Delete a face profile"""
        response = requests.delete(
            f"{self.api_base}/profile/{profile_id}?hard_delete={hard_delete}"
        )
        return response.json()
    
    def reindex(self, force=False):
        """Rebuild FAISS index (admin operation)"""
        response = requests.post(
            f"{self.api_base}/reindex",
            json={"force": force}
        )
        return response.json()


def main():
    """Run all examples"""
    print("=" * 60)
    print("Face Verification Service - Example Usage")
    print("=" * 60)
    print()
    
    print("Example: Health Check")
    print("-" * 60)
    client = FaceVerificationClient()
    try:
        health = client.health_check()
        print(f"Service is {health['status']}")
        print(f"Version: {health['version']}")
        print(f"Environment: {health['environment']}")
        print()
    except Exception as e:
        print(f"Error: {e}")
        print("Make sure the service is running on http://localhost:8080")
        return
    
    print("\nFor enrollment and verification examples,")
    print("see the individual example functions in this file.")
    print("Replace placeholder image paths with actual paths.")


if __name__ == "__main__":
    main()
