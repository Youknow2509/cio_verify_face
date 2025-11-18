"""Auth-related FastAPI dependencies for per-endpoint enforcement.

These read the authentication information attached by the Device/User
middlewares (`request.state.auth_type`, `request.state.user`, `request.state.device`)
and raise `HTTPException(401)` when the required auth is not present.
"""
from __future__ import annotations

from typing import Dict, Any

from fastapi import Depends, HTTPException
from starlette.requests import Request


def get_current_user(request: Request) -> Dict[str, Any]:
    """Require that the request is authenticated as a user and return user info."""
    if getattr(request.state, "auth_type", None) == "user" and getattr(request.state, "user", None):
        return request.state.user
    raise HTTPException(status_code=401, detail="User authentication required")


def get_current_device(request: Request) -> Dict[str, Any]:
    """Require that the request is authenticated as a device and return device info."""
    if getattr(request.state, "auth_type", None) == "device" and getattr(request.state, "device", None):
        return request.state.device
    raise HTTPException(status_code=401, detail="Device authentication required")


def get_current_auth(request: Request) -> Dict[str, Any]:
    """Require that the request is authenticated as either user or device.

    Returns a dict with keys: `auth_type` ("user"|"device") and `info`.
    """
    auth_type = getattr(request.state, "auth_type", None)
    if auth_type == "user" and getattr(request.state, "user", None):
        return {"auth_type": "user", "info": request.state.user}
    if auth_type == "device" and getattr(request.state, "device", None):
        return {"auth_type": "device", "info": request.state.device}
    raise HTTPException(status_code=401, detail="Authentication required")
