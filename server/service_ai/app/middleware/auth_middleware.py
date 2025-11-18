"""FastAPI middleware that validates Authorization bearer tokens via AuthService.

The middleware checks the `Authorization: Bearer <token>` header for requests
that start with `/api/` and attaches a `request.state.user` dict with parsed
token fields on success. On failure it returns `401 Unauthorized`.
"""
from __future__ import annotations

from typing import Optional, Iterable
import logging

from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.services.auth_client import get_client, AuthClient

_LOGGER = logging.getLogger(__name__)


def _default_skip_paths() -> Iterable[str]:
    return ("/health", "/metrics", "/docs", "/redoc", "/openapi.json", "/")


class AuthMiddleware(BaseHTTPMiddleware):
    """Middleware that validates incoming bearer tokens via AuthService.

    Behavior:
    - Skips validation for common public paths (health, docs, metrics)
    - For paths starting with `/api/` it requires `Authorization: Bearer <token>`
    - On successful validation attaches `request.state.user` with fields:
      `user_id`, `roles`, `token_id`, `company_id`, `expires_at`.
    """

    def __init__(self, app, client: Optional[AuthClient] = None, *, skip_paths: Optional[Iterable[str]] = None):
        super().__init__(app)
        self._client = client or get_client()
        self._skip_paths = tuple(skip_paths) if skip_paths is not None else tuple(_default_skip_paths())

    async def dispatch(self, request: Request, call_next):
        path = request.url.path

        # Skip public paths
        for p in self._skip_paths:
            if path == p or path.startswith(p.rstrip("/") + "/"):
                return await call_next(request)

        # Only enforce auth on /api endpoints (keep others public)
        if not path.startswith("/api/"):
            return await call_next(request)

        # Determine whether this request is from a device or a user.
        # Detection order (prefer device):
        # 1. If `X-Device-Token` header present => device
        # 2. Else if `X-Device-Id` present together with Authorization Bearer => device
        # 3. Else if Authorization Bearer present => user
        # 4. Else => unauthorized

        headers = request.headers
        # case-insensitive header access via .get()
        device_token = headers.get("x-device-token") or headers.get("X-Device-Token")
        device_id = headers.get("x-device-id") or headers.get("X-Device-Id")

        auth_hdr = headers.get("authorization") or headers.get("Authorization")
        bearer_token: Optional[str] = None
        if auth_hdr:
            parts = auth_hdr.split()
            if len(parts) == 2 and parts[0].lower() == "bearer":
                bearer_token = parts[1]

        # Prefer device token if explicit token header present
        client = getattr(request.app.state, "auth_client", self._client)

        # Case A: explicit device token header
        if device_token:
            if not device_id:
                return JSONResponse(status_code=401, content={"detail": "Missing X-Device-Id for device token"})
            try:
                resp = client.parse_device_token(device_token, device_id)
            except Exception:
                _LOGGER.exception("AuthService parse_device_token error")
                return JSONResponse(status_code=503, content={"detail": "Auth service unavailable"})

            if resp is None or not getattr(resp, "device_id", None):
                return JSONResponse(status_code=401, content={"detail": "Invalid or expired device token"})

            request.state.auth_type = "device"
            request.state.device = {
                "device_id": resp.device_id,
            }
            return await call_next(request)

        # Case B: Authorization bearer + explicit device id header -> treat as device token
        if bearer_token and device_id:
            try:
                resp = client.parse_device_token(bearer_token, device_id)
            except Exception:
                _LOGGER.exception("AuthService parse_device_token error")
                return JSONResponse(status_code=503, content={"detail": "Auth service unavailable"})

            if resp is not None and getattr(resp, "device_id", None):
                request.state.auth_type = "device"
                request.state.device = {
                    "device_id": resp.device_id,
                }
                return await call_next(request)

        # Case C: Authorization bearer => user token
        if bearer_token:
            try:
                resp = client.parse_user_token(bearer_token)
            except Exception:
                _LOGGER.exception("AuthService parse_user_token error")
                return JSONResponse(status_code=503, content={"detail": "Auth service unavailable"})

            if resp is None or not getattr(resp, "user_id", None):
                return JSONResponse(status_code=401, content={"detail": "Invalid or expired user token"})

            # Attach user info to request.state for downstream handlers
            request.state.auth_type = "user"
            request.state.user = {
                "user_id": resp.user_id,
                "roles": getattr(resp, "roles", None),
                "token_id": getattr(resp, "token_id", None),
                "company_id": getattr(resp, "company_id", None),
                "expires_at": getattr(resp, "exprires_at", None),
            }
            return await call_next(request)

        # No valid auth presented
        return JSONResponse(status_code=401, content={"detail": "Missing authorization headers"})
