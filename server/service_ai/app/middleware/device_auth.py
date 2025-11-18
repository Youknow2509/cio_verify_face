"""Middleware handling device JWTs only.

Checks for device-specific headers (`X-Device-Token` or `Authorization: Bearer <token>`
with `X-Device-Id`) and validates via `AuthClient.parse_device_token`. If a device
token is present and valid it sets `request.state.auth_type = 'device'` and
`request.state.device`. If no device token is present the middleware is non-blocking
and simply calls the next handler.
"""
from __future__ import annotations

from typing import Optional, Iterable
import logging

from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.services.auth_client import get_client, AuthClient

_LOGGER = logging.getLogger(__name__)


class DeviceAuthMiddleware(BaseHTTPMiddleware):
    """Non-blocking device authentication middleware.

    Only acts when device token/headers are present. On success attaches
    `request.state.device` and `request.state.auth_type = 'device'`.
    On failure it returns a `401` or `503` (if AuthService unavailable).
    """

    def __init__(self, app, client: Optional[AuthClient] = None, *, skip_paths: Optional[Iterable[str]] = None):
        super().__init__(app)
        self._client = client or get_client()

    async def dispatch(self, request: Request, call_next):
        headers = request.headers
        device_token = headers.get("x-device-token") or headers.get("X-Device-Token")
        device_id = headers.get("x-device-id") or headers.get("X-Device-Id")
        auth_hdr = headers.get("authorization") or headers.get("Authorization")

        bearer_token: Optional[str] = None
        if auth_hdr:
            parts = auth_hdr.split()
            if len(parts) == 2 and parts[0].lower() == "bearer":
                bearer_token = parts[1]

        # Prefer explicit device token header first
        client = getattr(request.app.state, "auth_client", self._client)

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
            request.state.device = {"device_id": resp.device_id}
            return await call_next(request)

        # If bearer token present together with device_id header, try treat as device token
        if bearer_token and device_id:
            try:
                resp = client.parse_device_token(bearer_token, device_id)
            except Exception:
                _LOGGER.exception("AuthService parse_device_token error")
                return JSONResponse(status_code=503, content={"detail": "Auth service unavailable"})

            if resp is not None and getattr(resp, "device_id", None):
                request.state.auth_type = "device"
                request.state.device = {"device_id": resp.device_id}
                return await call_next(request)

        # No device auth present — non-blocking, proceed
        return await call_next(request)
