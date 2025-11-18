"""Middleware handling user JWTs only.

This middleware attempts to validate an Authorization Bearer token as a user
JWT (only if no device auth was already attached by `DeviceAuthMiddleware`).
If a valid user token is present it attaches `request.state.user` and
`request.state.auth_type = 'user'`. If no bearer token is present it is
non-blocking and simply calls the next handler.
"""
from __future__ import annotations

from typing import Optional
import logging

from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.services.auth_client import get_client, AuthClient

_LOGGER = logging.getLogger(__name__)


class UserAuthMiddleware(BaseHTTPMiddleware):
    """Non-blocking user authentication middleware.

    Skips if `request.state.auth_type == 'device'` (device already authenticated).
    """

    def __init__(self, app, client: Optional[AuthClient] = None):
        super().__init__(app)
        self._client = client or get_client()

    async def dispatch(self, request: Request, call_next):
        # If device already authenticated, skip user auth
        if getattr(request.state, "auth_type", None) == "device":
            return await call_next(request)

        auth_hdr = request.headers.get("authorization") or request.headers.get("Authorization")
        if not auth_hdr:
            return await call_next(request)

        parts = auth_hdr.split()
        if len(parts) != 2 or parts[0].lower() != "bearer":
            return await call_next(request)

        token = parts[1]
        client = getattr(request.app.state, "auth_client", self._client)

        try:
            resp = client.parse_user_token(token)
        except Exception:
            _LOGGER.exception("AuthService parse_user_token error")
            return JSONResponse(status_code=503, content={"detail": "Auth service unavailable"})

        if resp is None or not getattr(resp, "user_id", None):
            # invalid user token — non-blocking here (RequireAuthMiddleware enforces presence)
            return await call_next(request)

        request.state.auth_type = "user"
        request.state.user = {
            "user_id": resp.user_id,
            "roles": getattr(resp, "roles", None),
            "token_id": getattr(resp, "token_id", None),
            "company_id": getattr(resp, "company_id", None),
            "expires_at": getattr(resp, "exprires_at", None),
        }

        return await call_next(request)
