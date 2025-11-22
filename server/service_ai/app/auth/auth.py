
# Check manager or admin role
from uuid import UUID
from app.models.enum import ADMIN_ROLE, MANAGER_ROLE
from app.models.schemas import SessionUser



def is_manager_or_admin(company_id: UUID, token_payload: SessionUser) -> bool:
    if token_payload.role == ADMIN_ROLE:
        return True
    if token_payload.role == MANAGER_ROLE and token_payload.company_id == company_id:
        return True
    return False