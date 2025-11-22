"""
Upgraded gRPC client wrapper for the Auth service.
Provides retry, backoff, metadata, exception handling, and robust channel management.
"""

from __future__ import annotations

import logging
import time
from typing import Optional, Sequence, Any, Callable

import grpc
from google.protobuf import empty_pb2

from app.grpc_generated import auth_pb2, auth_pb2_grpc
from app.core.config import settings

_LOGGER = logging.getLogger(__name__)


# ============================================================
#                     CUSTOM EXCEPTION
# ============================================================
class AuthRPCError(RuntimeError):
    """Custom exception wrapper for Auth RPC errors."""
    def __init__(self, message: str, status_code: grpc.StatusCode, original: grpc.RpcError):
        super().__init__(f"{message} (grpc_status={status_code.name})")
        self.status_code = status_code
        self.original = original


# ============================================================
#                     AUTH CLIENT
# ============================================================
class AuthClient:
    DEFAULT_TARGET = "localhost:50051"

    NON_RETRY_CODES = {
        grpc.StatusCode.INVALID_ARGUMENT,
        grpc.StatusCode.PERMISSION_DENIED,
        grpc.StatusCode.UNAUTHENTICATED,
        grpc.StatusCode.NOT_FOUND,
        grpc.StatusCode.ALREADY_EXISTS,
        grpc.StatusCode.FAILED_PRECONDITION,
        grpc.StatusCode.OUT_OF_RANGE,
        grpc.StatusCode.UNIMPLEMENTED,
    }

    def __init__(
        self,
        target: Optional[str] = None,
        *,
        timeout: float = 5.0,
        max_retries: int = 0,
        retry_backoff: float = 0.5,
        exponential_backoff: bool = True,
        metadata: Optional[Sequence[tuple[str, str]]] = None,
        channel_options: Optional[Sequence[tuple[str, Any]]] = None,
        raise_exception: bool = False,
        insecure: bool = True,
    ):
        self._target = target or getattr(settings, 'GRPC_AUTH_URL', self.DEFAULT_TARGET)
        self._timeout = timeout
        self._max_retries = max_retries
        self._retry_backoff = retry_backoff
        self._exponential_backoff = exponential_backoff
        self._default_metadata = metadata or []
        self._raise_exception = raise_exception
        self._insecure = insecure

        self._channel_options = channel_options or (
            ("grpc.keepalive_time_ms", settings.GRPC_CLIENT_KEEPALIVE_TIME_MS),
            ("grpc.keepalive_timeout_ms", settings.GRPC_CLIENT_KEEPALIVE_TIMEOUT_MS),
            ("grpc.keepalive_permit_without_calls", settings.GRPC_CLIENT_KEEPALIVE_PERMIT_WITHOUT_CALLS),
            ("grpc.http2.max_pings_without_data", settings.GRPC_CLIENT_HTTP2_MAX_PINGS_WITHOUT_DATA),
            ("grpc.http2.min_time_between_pings_ms", settings.GRPC_CLIENT_HTTP2_MIN_TIME_BETWEEN_PINGS_MS),
            ("grpc.http2.min_ping_interval_without_data_ms", settings.GRPC_CLIENT_HTTP2_MIN_PING_INTERVAL_WITHOUT_DATA_MS),
        )

        self._channel: Optional[grpc.Channel] = None
        self._stub: Optional[auth_pb2_grpc.AuthServiceStub] = None
        self._ensure_channel()

    # --------------------------------------------------------
    #               CHANNEL MANAGEMENT
    # --------------------------------------------------------
    def _ensure_channel(self) -> None:
        if self._channel is None:
            if self._insecure:
                self._channel = grpc.insecure_channel(self._target, options=self._channel_options)
            else:
                creds = grpc.ssl_channel_credentials()
                self._channel = grpc.secure_channel(self._target, creds, options=self._channel_options)
            self._stub = auth_pb2_grpc.AuthServiceStub(self._channel)

    def recreate_channel(self) -> None:
        self.close()
        self._channel = None
        self._stub = None
        self._ensure_channel()

    def __enter__(self) -> "AuthClient":
        self._ensure_channel()
        return self

    def __exit__(self, exc_type, exc, tb):
        self.close()

    def close(self) -> None:
        try:
            if self._channel:
                self._channel.close()
        except Exception:
            _LOGGER.exception("Error closing AuthClient channel")

    # --------------------------------------------------------
    #               INTERNAL RETRY WRAPPER
    # --------------------------------------------------------
    def _call_with_retry(
        self,
        func: Callable[..., Any],
        request: Any,
        *,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None,
        description: str = "",
    ):
        attempt = 0
        backoff = self._retry_backoff
        deadline = timeout or self._timeout

        while True:
            try:
                return func(request, timeout=deadline, metadata=metadata or self._default_metadata)

            except grpc.RpcError as exc:
                attempt += 1
                code = exc.code()
                msg = description or func.__name__

                can_retry = attempt <= self._max_retries and code not in self.NON_RETRY_CODES

                if not can_retry:
                    _LOGGER.error("RPC '%s' failed (attempt %d) code=%s details=%s",
                                  msg, attempt, code.name, exc.details())
                    if self._raise_exception:
                        raise AuthRPCError(f"RPC '{msg}' failed", code, exc)
                    return None

                _LOGGER.warning("Retrying '%s' (attempt %d/%d) code=%s",
                                msg, attempt, self._max_retries, code.name)

                if backoff > 0:
                    time.sleep(backoff)
                    if self._exponential_backoff:
                        backoff *= 2

    # --------------------------------------------------------
    #                 HEALTH CHECK
    # --------------------------------------------------------
    def check_connection(self, *, timeout: Optional[float] = None, metadata=None) -> bool:
        req = empty_pb2.Empty()
        try:
            self._stub.HealthCheck(req, timeout=timeout or self._timeout, metadata=metadata or self._default_metadata)
            return True
        except grpc.RpcError as exc:
            if self._raise_exception:
                raise AuthRPCError("HealthCheck failed", exc.code(), exc)
            _LOGGER.error("HealthCheck RPC failed: %s", exc)
            return False

    # --------------------------------------------------------
    #               AUTH SERVICE METHODS
    # --------------------------------------------------------
    def create_user_token(self, user_id: str, roles: int, *, metadata=None):
        req = auth_pb2.CreateUserTokenRequest(user_id=user_id, roles=roles)
        return self._call_with_retry(
            self._stub.CreateUserToken,
            req,
            metadata=metadata,
            description="CreateUserToken",
        )

    def create_service_token(self, service_id: str, *, metadata=None):
        req = auth_pb2.CreateServiceTokenRequest(service_id=service_id)
        return self._call_with_retry(
            self._stub.CreateServiceToken,
            req,
            metadata=metadata,
            description="CreateServiceToken",
        )

    def create_device_token(self, device_id: str, company_id: str, *, metadata=None):
        req = auth_pb2.CreateDeviceTokenRequest(device_id=device_id, company_id=company_id)
        return self._call_with_retry(
            self._stub.CreateDeviceToken,
            req,
            metadata=metadata,
            description="CreateDeviceToken",
        )

    def parse_user_token(self, token: str, *, metadata=None):
        req = auth_pb2.ParseUserTokenRequest(token=token)
        return self._call_with_retry(
            self._stub.ParseUserToken,
            req,
            metadata=metadata,
            description="ParseUserToken",
        )

    def parse_service_token(self, service_id: str, *, metadata=None):
        req = auth_pb2.ParseServiceTokenRequest(service_id=service_id)
        return self._call_with_retry(
            self._stub.ParseServiceToken,
            req,
            metadata=metadata,
            description="ParseServiceToken",
        )

    def parse_device_token(self, token: str, device_id: str, *, metadata=None):
        req = auth_pb2.ParseDeviceTokenRequest(token=token, device_id=device_id)
        return self._call_with_retry(
            self._stub.ParseDeviceToken,
            req,
            metadata=metadata,
            description="ParseDeviceToken",
        )


# ============================================================
#                       FACTORY
# ============================================================
def get_client(
    target: Optional[str] = None,
    timeout: float = 5.0,
    *,
    max_retries: int = 0,
    retry_backoff: float = 0.5,
    exponential_backoff: bool = True,
    raise_exception: bool = False,
    insecure: bool = True,
    metadata: Optional[Sequence[tuple[str, str]]] = None,
) -> AuthClient:
    return AuthClient(
        target=target,
        timeout=timeout,
        max_retries=max_retries,
        retry_backoff=retry_backoff,
        exponential_backoff=exponential_backoff,
        raise_exception=raise_exception,
        insecure=insecure,
        metadata=metadata,
    )
