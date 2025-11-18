"""gRPC client wrapper for the Auth service.

Provides a small, testable wrapper around the generated gRPC stubs
in `app.grpc_generated` so other parts of the application can call
auth methods conveniently.
"""
from __future__ import annotations

import os
import logging
from typing import Optional

import grpc

from ..grpc_generated import auth_pb2, auth_pb2_grpc
from app.core.config import settings

_LOGGER = logging.getLogger(__name__)


class AuthClient:
    """Simple gRPC client for the `auth.AuthService`.

    The client uses `GRPC_CLIENT_URL` environment variable if no
    `target` is provided. The value should be in the form
    `host:port` (e.g. `localhost:50051`).
    """

    def __init__(self, target: Optional[str] = None, *, timeout: float = 5.0):
        self._target = target or settings.GRPC_CLIENT_URL
        self._timeout = timeout
        self._channel = grpc.insecure_channel(self._target)
        self._stub = auth_pb2_grpc.AuthServiceStub(self._channel)

    def close(self) -> None:
        try:
            self._channel.close()
        except Exception:
            _LOGGER.exception("Error closing auth client channel")

    def create_user_token(self, user_id: str, roles: int) -> Optional[auth_pb2.CreateUserTokenResponse]:
        req = auth_pb2.CreateUserTokenRequest(user_id=user_id, roles=roles)
        try:
            return self._stub.CreateUserToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("CreateUserToken rpc failed: %s", exc)
            return None

    def create_service_token(self, service_id: str) -> Optional[auth_pb2.CreateServiceTokenResponse]:
        req = auth_pb2.CreateServiceTokenRequest(service_id=service_id)
        try:
            return self._stub.CreateServiceToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("CreateServiceToken rpc failed: %s", exc)
            return None

    def create_device_token(self, device_id: str, company_id: str) -> Optional[auth_pb2.CreateDeviceTokenResponse]:
        req = auth_pb2.CreateDeviceTokenRequest(device_id=device_id, company_id=company_id)
        try:
            return self._stub.CreateDeviceToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("CreateDeviceToken rpc failed: %s", exc)
            return None

    def parse_user_token(self, token: str) -> Optional[auth_pb2.ParseUserTokenResponse]:
        req = auth_pb2.ParseUserTokenRequest(token=token)
        try:
            return self._stub.ParseUserToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("ParseUserToken rpc failed: %s", exc)
            return None

    def parse_service_token(self, service_id: str) -> Optional[auth_pb2.ParseServiceTokenResponse]:
        req = auth_pb2.ParseServiceTokenRequest(service_id=service_id)
        try:
            return self._stub.ParseServiceToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("ParseServiceToken rpc failed: %s", exc)
            return None

    def parse_device_token(self, token: str, device_id: str) -> Optional[auth_pb2.ParseDeviceTokenResponse]:
        req = auth_pb2.ParseDeviceTokenRequest(token=token, device_id=device_id)
        try:
            return self._stub.ParseDeviceToken(req, timeout=self._timeout)
        except grpc.RpcError as exc:
            _LOGGER.error("ParseDeviceToken rpc failed: %s", exc)
            return None


def get_client(target: Optional[str] = None, timeout: float = 5.0) -> AuthClient:
    return AuthClient(target, timeout=timeout)
