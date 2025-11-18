"""Middleware that enforces that a request has been authenticated.

This should be placed after device/user auth middlewares. It will enforce
that requests under `/api/` have either `request.state.device` or
`request.state.user` set. It skips common public paths.
"""
from __future__ import annotations

from typing import Iterable
import logging

from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

_LOGGER = logging.getLogger(__name__)


def _default_skip_paths() -> Iterable[str]:
    return ("/health", "/metrics", "/docs", "/redoc", "/openapi.json", "/")


class RequireAuthMiddleware(BaseHTTPMiddleware):
    """Enforce authentication for `/api/` endpoints.

    Place after DeviceAuthMiddleware and UserAuthMiddleware.
    """

    def __init__(self, app, *, skip_paths: Iterable[str] | None = None):
        super().__init__(app)
        self._skip_paths = tuple(skip_paths) if skip_paths is not None else tuple(_default_skip_paths())

    async def dispatch(self, request: Request, call_next):
        path = request.url.path

        # Skip public paths
        for p in self._skip_paths:
            if path == p or path.startswith(p.rstrip('/') + '/'):
                return await call_next(request)

        # Only enforce for /api/
        if not path.startswith('/api/'):
            return await call_next(request)

        # If either user or device auth was attached, proceed
        if getattr(request.state, 'auth_type', None) in ('user', 'device'):
            return await call_next(request)

        # Not authenticated
        return JSONResponse(status_code=401, content={"detail": "Authentication required"})
