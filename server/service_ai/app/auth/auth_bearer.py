import logging
from typing import Optional, Callable, Union

from fastapi import Request, HTTPException
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer

from app.models.schemas import SessionUser
from app.grpc.client.auth_client import AuthClient

logger = logging.getLogger(__name__)

class JWTBearer(HTTPBearer):
    def __init__(self, auto_error: bool = True):
        super().__init__(auto_error=auto_error)
        self._auth_grpc_client = None

    def _client(self) -> AuthClient:
        # Lazy initialization to avoid circular import
        if self._auth_grpc_client is None:
            from app.main import app
            self._auth_grpc_client = app.state.auth_client
            if self._auth_grpc_client is None:
                raise ValueError("Auth gRPC client is not initialized in app state")
        
        # Resolve factory vs instance
        return self._auth_grpc_client() if callable(self._auth_grpc_client) else self._auth_grpc_client

    async def __call__(self, request: Request) -> SessionUser:
        credentials: HTTPAuthorizationCredentials = await super().__call__(request)
        if credentials is None:
            raise HTTPException(status_code=403, detail="No authorization token provided")
        if credentials.scheme != "Bearer":
            raise HTTPException(status_code=403, detail="Invalid authentication scheme")

        token = credentials.credentials
        return await self.verify_jwt_via_grpc(token)

    async def verify_jwt_via_grpc(self, token: str) -> SessionUser:
        try:
            client = self._client()
            resp = client.parse_user_token(token=token)
            return SessionUser(
                user_id=resp.user_id,
                role=resp.roles,
                company_id=resp.company_id,
                session_id=resp.token_id,
                exprires_at=resp.exprires_at,
            )
        except HTTPException:
            raise
        except Exception as e:
            logger.error("gRPC error verifying token: %s", e)
            raise HTTPException(status_code=500, detail="Authentication service error") from e