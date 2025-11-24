"""
    Utilities for caching.
"""

# =============================================================
# Utils get ttl cache
# =============================================================
def get_ttl_cache_default() -> int:
    """Get default TTL for cache in seconds"""
    return 3600  # 1 hour

def get_ttl_cache_short() -> int:
    """Get short TTL for cache in seconds"""
    return 300  # 5 minutes

def get_ttl_cache_long() -> int:
    """Get long TTL for cache in seconds"""
    return 86400  # 24 hours

# =============================================================
# Utils get key cache
# =============================================================

def get_user_profile_face_cache_key(user_id_hash: str, company_id_hash: str, page_number: int, page_size: int) -> str:
    """Generate cache key for user profile"""
    return f"employee:profile_face:{company_id_hash}:{user_id_hash}:{page_number}:{page_size}"

def get_employee_company_cache_key(company_id_hash: str, employee_id_hash: str) -> str:
    """Generate cache key for employee in company"""
    return f"employee_in_company:{company_id_hash}:{employee_id_hash}"