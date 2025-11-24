"""Database utilities."""

import uuid

def get_face_profile_partition_name(company_id: str | uuid.UUID) -> str:
    """Get face profile partition name in postgres based on user_id."""
    return f"face_profiles_p_{str(company_id).replace('-', '')}"
