"""
    User Service Module
"""

import json
import logging
from typing import List
from uuid import UUID
from app.database.pg_manager import PGManager
from app.database.redis_manager import RedisManager
from app.models.schemas import FaceProfileResponse
from app.utils.cache import get_employee_company_cache_key, get_ttl_cache_short, get_user_profile_face_cache_key
from app.utils.crypto import getHash

logger = logging.getLogger(__name__)

class UserService:
    """
        Service class for user-related operations.
    """

    def __init__(self, redis_client: RedisManager, postgres_client: PGManager):
        self.redis_client = redis_client
        self.postgres_client = postgres_client

    def _distribute_cache(self) -> RedisManager:
        return self.redis_client
    
    def _db(self) -> PGManager:
        return self.postgres_client
    
    async def get_profile_face_user(
        self,
        user_id: UUID,
        company_id: UUID,
        page_size: int = 20,
        page_number: int = 1
    ) -> List[FaceProfileResponse]:
        """
            Retrieve face profiles associated with a user in a specific company.

            :param user_id: UUID of the user
            :param company_id: UUID of the company
            :return: List of FaceProfileResponse objects
        """
        # Check in cache
        key = get_user_profile_face_cache_key(
            user_id_hash=getHash(str(user_id)),
            company_id_hash=getHash(str(company_id)),
            page_number=page_number,
            page_size=page_size
        )
        cached_profiles = self._distribute_cache().get(key)
        if cached_profiles is not None:
            try:
                profiles_data = json.loads(cached_profiles)
                return [FaceProfileResponse.parse_obj(p) for p in profiles_data]
            except Exception as e:
                logger.error(
                    f"Error parsing cached profiles for key {key}"
                    f" Error: {e}"
                )
                pass
        
        # Get from db
        raw_profiles = self._db().get_profile_face_employee(
            employee_id=user_id,
            company_id=company_id,
            page_size=page_size,
            page_number=page_number
        )
        profiles = [FaceProfileResponse.parse_obj(p) for p in raw_profiles]
        # Store in cache as a single JSON blob
        try:
            self._distribute_cache().set(
                key,
                json.dumps([p.dict() for p in profiles]),
                ex=get_ttl_cache_short()
            )
        except Exception as e:
            logger.error(
                f"Error caching profiles for key {key}"
                f" Error: {e}"
            )
            pass

        return profiles
    
    async def check_user_exist_in_company(
        self,
        user_id: UUID,
        company_id: UUID,
    ) -> bool:
        """
            Check if a user exists in a specific company.

            :param user_id: UUID of the user
            :param company_id: UUID of the company
            :return: True if the user exists in the company, False otherwise
        """
        # Cache in cache
        key = get_employee_company_cache_key(
            company_id_hash=getHash(str(company_id)),
            employee_id_hash=getHash(str(user_id))
        )
        resp = self._distribute_cache().get(key)
        if resp is not None and resp == "0":
            return True
        if resp is not None and resp == "1":
            return False
        
        # Check in db
        ok = self._db().check_employee_exist_in_company(
            company_id=company_id,
            employee_id=user_id,
        )
        if ok:
            self._distribute_cache().set(key, "0", ex=get_ttl_cache_short())
        else:
            self._distribute_cache().set(key, "1", ex=get_ttl_cache_short())
        return ok