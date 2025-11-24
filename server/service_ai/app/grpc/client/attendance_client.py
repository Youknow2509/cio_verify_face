from __future__ import annotations

import logging
import time
from typing import (
    Optional,
    Iterator,
    Iterable,
    Callable,
    Generator,
    Sequence,
    Dict,
    Any,
    Tuple,
    Union,
)

import grpc
from google.protobuf import empty_pb2

from app.grpc_generated import attendance_pb2, attendance_pb2_grpc
from app.core.config import settings

_LOGGER = logging.getLogger(__name__)


class AttendanceRPCError(RuntimeError):
    """Ngoại lệ chung cho lỗi RPC AttendanceService."""
    def __init__(self, message: str, status_code: grpc.StatusCode, original: grpc.RpcError):
        super().__init__(f"{message} (grpc_status={status_code.name})")
        self.status_code = status_code
        self.original = original


class AttendanceClient:
    """
    Client gRPC cho AttendanceService.
    """

    DEFAULT_TARGET = "localhost:50052"

    def __init__(
        self,
        target: Optional[str] = None,
        *,
        timeout: float = 5.0,
        max_retries: int = 0,
        retry_backoff: float = 0.5,
        exponential_backoff: bool = True,
        channel_options: Optional[Sequence[tuple[str, Any]]] = None,
        raise_exception: bool = False,
        insecure: bool = True,
        metadata: Optional[Sequence[tuple[str, str]]] = None,
    ):
        """
        Args:
            target: host:port của AttendanceService.
            timeout: timeout mỗi RPC (deadline).
            max_retries: số lần retry khi lỗi (Status != OK và không thuộc nhóm không retry).
            retry_backoff: thời gian nghỉ giữa các lần retry (giây) cho lần đầu.
            exponential_backoff: nếu True => backoff *= 2 sau mỗi lần retry.
            channel_options: tuỳ chọn channel grpc (list các tuple).
            raise_exception: nếu True => ném AttendanceRPCError thay vì trả về None khi lỗi.
            insecure: dùng insecure_channel (False => yêu cầu ssl credentials bên ngoài).
            metadata: metadata chung đính kèm mọi RPC (có thể override khi gọi từng hàm).
        """
        self._target = target or getattr(settings, "GRPC_ATTENDANCE_URL", self.DEFAULT_TARGET)
        self._timeout = timeout
        self._max_retries = max_retries
        self._retry_backoff = retry_backoff
        self._exponential_backoff = exponential_backoff
        self._raise_exception = raise_exception
        self._default_metadata = metadata or []
        self._channel_options = channel_options or (
            ("grpc.keepalive_time_ms", settings.GRPC_CLIENT_KEEPALIVE_TIME_MS),                 # Ping mỗi 120s
            ("grpc.keepalive_timeout_ms", settings.GRPC_CLIENT_KEEPALIVE_TIMEOUT_MS),               # 20s chờ trước khi reset
            ("grpc.keepalive_permit_without_calls", settings.GRPC_CLIENT_KEEPALIVE_PERMIT_WITHOUT_CALLS),          # Không ping khi idle
            ("grpc.http2.max_pings_without_data", settings.GRPC_CLIENT_HTTP2_MAX_PINGS_WITHOUT_DATA),            # Chỉ 1 ping idle
            ("grpc.http2.min_time_between_pings_ms", settings.GRPC_CLIENT_HTTP2_MIN_TIME_BETWEEN_PINGS_MS),    # Tối thiểu 60s giữa 2 ping
            ("grpc.http2.min_ping_interval_without_data_ms", settings.GRPC_CLIENT_HTTP2_MIN_PING_INTERVAL_WITHOUT_DATA_MS),
        )

        self._insecure = insecure
        self._channel: Optional[grpc.Channel] = None
        self._stub: Optional[attendance_pb2_grpc.AttendanceServiceStub] = None

        self._ensure_channel()

    # ------------- Channel Management -------------
    def _ensure_channel(self) -> None:
        if self._channel is None:
            if self._insecure:
                self._channel = grpc.insecure_channel(self._target, options=self._channel_options)
            else:
                creds = grpc.ssl_channel_credentials()
                self._channel = grpc.secure_channel(self._target, creds, options=self._channel_options)
            self._stub = attendance_pb2_grpc.AttendanceServiceStub(self._channel)

    def recreate_channel(self) -> None:
        """Đóng và tạo lại channel (khi cần)."""
        self.close()
        self._channel = None
        self._stub = None
        self._ensure_channel()

    # ------------- Context Manager -------------
    def __enter__(self) -> "AttendanceClient":
        self._ensure_channel()
        return self

    def __exit__(self, exc_type, exc, tb):
        self.close()

    # ------------- Health -------------
    def check_connection(self, *, timeout: Optional[float] = None, metadata: Optional[Sequence[tuple[str, str]]] = None) -> bool:
        """
        Kiểm tra kết nối đến AttendanceService.
        Returns:
            True nếu HealthCheck OK, False nếu lỗi (và không raise).
        """
        req = empty_pb2.Empty()
        try:
            self._stub.HealthCheck(req, timeout=timeout or self._timeout, metadata=metadata or self._default_metadata)
            return True
        except grpc.RpcError as exc:
            if self._raise_exception:
                raise AttendanceRPCError("HealthCheck failed", exc.code(), exc)
            _LOGGER.error("HealthCheck rpc failed: %s", exc)
            return False

    # ------------- Internal Helpers -------------
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

    def _call_with_retry(
        self,
        func: Callable[..., Any],
        request_or_iterator: Any,
        *,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None,
        is_stream: bool = False,
        description: str = "",
    ) -> Any:
        """
        Wrapper thực thi RPC với retry.
        Args:
            func: stub method.
            request_or_iterator: request message hoặc iterator streaming.
            timeout: override deadline.
            metadata: metadata cho call.
            is_stream: True nếu là client-streaming RPC.
            description: mô tả ngắn (log).
        """
        attempt = 0
        backoff = self._retry_backoff
        deadline = timeout or self._timeout

        while True:
            try:
                if is_stream:
                    return func(request_or_iterator, timeout=deadline, metadata=metadata or self._default_metadata)
                return func(request_or_iterator, timeout=deadline, metadata=metadata or self._default_metadata)
            except grpc.RpcError as exc:
                code = exc.code()
                attempt += 1
                can_retry = attempt <= self._max_retries and code not in self.NON_RETRY_CODES
                msg = description or func.__name__
                if not can_retry:
                    _LOGGER.error("RPC '%s' failed (attempt %d) code=%s details=%s",
                                  msg, attempt, code.name, exc.details())
                    if self._raise_exception:
                        raise AttendanceRPCError(f"RPC '{msg}' failed", code, exc)
                    return None
                _LOGGER.warning("Retrying RPC '%s' (attempt %d/%d) code=%s",
                                msg, attempt, self._max_retries, code.name)
                if backoff > 0:
                    time.sleep(backoff)
                    if self._exponential_backoff:
                        backoff *= 2

    # ------------- Close -------------
    def close(self) -> None:
        if self._channel:
            try:
                self._channel.close()
            except Exception:
                _LOGGER.exception("Error closing attendance client channel")

    # ------------- Single Attendance -------------
    def add_attendance(
        self,
        *,
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        session: attendance_pb2.SessionInfo,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.AddAttendanceOutput]:
        """Ghi một bản ghi điểm danh đơn lẻ."""
        if verification_score < 0:
            _LOGGER.warning("verification_score < 0, auto clamp to 0")
            verification_score = 0.0
        req = attendance_pb2.AddAttendanceInput(
            company_id=company_id,
            employee_id=employee_id,
            device_id=device_id,
            record_time=record_time,
            verification_method=verification_method,
            verification_score=verification_score,
            face_image_url=face_image_url,
            location_coordinates=location_coordinates,
            session=session,
        )
        return self._call_with_retry(
            self._stub.AddAttendance,
            req,
            timeout=timeout,
            metadata=metadata,
            description="AddAttendance",
        )

    # ------------- Records Pagination (Raw) -------------
    def get_attendance_records(
        self,
        *,
        company_id: str,
        employee_id: str,
        year_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.GetAttendanceRecordsOutput]:
        """Lấy danh sách điểm danh theo tháng (phân trang)."""
        if page_size <= 0:
            page_size = 50
        req = attendance_pb2.GetAttendanceRecordsInput(
            company_id=company_id,
            employee_id=employee_id,
            year_month=year_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session,
        )
        return self._call_with_retry(
            self._stub.GetAttendanceRecords,
            req,
            timeout=timeout,
            metadata=metadata,
            description="GetAttendanceRecords",
        )

    def get_attendance_records_employee(
        self,
        *,
        company_id: str,
        employee_id: str,
        year_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.GetAttendanceRecordsOutput]:
        """Lấy danh sách điểm danh của một nhân viên (phân trang)."""
        if page_size <= 0:
            page_size = 50
        req = attendance_pb2.GetAttendanceRecordsEmployeeInput(
            company_id=company_id,
            employee_id=employee_id,
            year_month=year_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session,
        )
        return self._call_with_retry(
            self._stub.GetAttendanceRecordsEmployee,
            req,
            timeout=timeout,
            metadata=metadata,
            description="GetAttendanceRecordsEmployee",
        )

    # ------------- Daily Summary (Raw) -------------
    def get_daily_attendance_summary(
        self,
        *,
        company_id: str,
        employee_id: str,
        summary_month: str,
        work_date: int,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.GetDailyAttendanceSummaryOutput]:
        """Lấy tổng hợp điểm danh ngày (phân trang)."""
        if page_size <= 0:
            page_size = 50
        req = attendance_pb2.GetDailyAttendanceSummaryInput(
            company_id=company_id,
            employee_id=employee_id,
            summary_month=summary_month,
            work_date=work_date,
            page_size=page_size,
            page_stage=page_stage,
            session=session,
        )
        return self._call_with_retry(
            self._stub.GetDailyAttendanceSummary,
            req,
            timeout=timeout,
            metadata=metadata,
            description="GetDailyAttendanceSummary",
        )

    def get_daily_attendance_summary_employee(
        self,
        *,
        company_id: str,
        employee_id: str,
        summary_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.GetDailyAttendanceSummaryOutput]:
        """Lấy tổng hợp điểm danh theo tháng (của 1 nhân viên, phân trang)."""
        if page_size <= 0:
            page_size = 50
        req = attendance_pb2.GetDailyAttendanceSummaryEmployeeInput(
            company_id=company_id,
            employee_id=employee_id,
            summary_month=summary_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session,
        )
        return self._call_with_retry(
            self._stub.GetDailyAttendanceSummaryEmployee,
            req,
            timeout=timeout,
            metadata=metadata,
            description="GetDailyAttendanceSummaryEmployee",
        )

    # ------------- Batch Attendance (Streaming) -------------
    def add_batch_attendance(
        self,
        attendance_records: Iterator[attendance_pb2.AddAttendanceInput],
        *,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Optional[attendance_pb2.AddAttendanceOutput]:
        """
        Ghi nhiều bản ghi điểm danh bằng client streaming.
        Timeout mặc định: tăng gấp 10 nếu không truyền cụ thể.
        """
        return self._call_with_retry(
            self._stub.AddBatchAttendance,
            attendance_records,
            timeout=timeout or (self._timeout * 10),
            metadata=metadata,
            is_stream=True,
            description="AddBatchAttendance",
        )

    def service_add_batch_attendance(
        self,
        request_iterator,
        timeout: Optional[float] = None,
        metadata: Optional[Sequence[Tuple[str, str]]] = None,
    ) -> Optional[attendance_pb2.ServiceAddBatchAttendanceOutput]:
        """Stream multiple attendance records"""
        try:
            _LOGGER.debug("Starting service_add_batch_attendance stream")
            
            if not self._channel or not self._stub:
                _LOGGER.error("Channel or stub not initialized")
                return None
                
            response = self._stub.ServiceAddBatchAttendance(
                request_iterator,
                timeout=timeout,
                metadata=metadata or []
            )
            
            _LOGGER.debug("service_add_batch_attendance completed: status=%s", 
                        getattr(response, 'status_code', 'unknown'))
            return response
            
        except grpc.RpcError as e:
            _LOGGER.error(
                "gRPC error in service_add_batch_attendance: code=%s details=%s",
                e.code(), e.details()
            )
            return None
        except Exception as e:
            _LOGGER.exception("Unexpected error in service_add_batch_attendance: %s", str(e))
            return None

    # ------------- Pagination Helpers (High-Level Iterators) -------------
    @staticmethod
    def _decode_page_stage_next(page_stage_next: bytes) -> str:
        """Giải mã page_stage_next (bytes) -> str (utf-8) hoặc hex nếu không phải utf-8."""
        if not page_stage_next:
            return ""
        try:
            return page_stage_next.decode("utf-8")
        except UnicodeDecodeError:
            return page_stage_next.hex()

    def paginate_attendance_records(
        self,
        *,
        company_id: str,
        employee_id: str,
        year_month: str,
        page_size: int,
        session: attendance_pb2.SessionInfo,
        use_employee_api: bool = False,
        start_page_stage: str = "",
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Generator[attendance_pb2.AttendanceRecordInfo, None, None]:
        """
        Generator trả về tất cả AttendanceRecordInfo qua các trang.
        """
        page_stage = start_page_stage
        fetch_fn = self.get_attendance_records_employee if use_employee_api else self.get_attendance_records

        while True:
            resp: attendance_pb2.GetAttendanceRecordsOutput = fetch_fn(
                company_id=company_id,
                employee_id=employee_id,
                year_month=year_month,
                page_size=page_size,
                page_stage=page_stage,
                session=session,
                metadata=metadata,
            )
            if not resp:
                break
            for rec in resp.records:
                yield rec
            next_stage = self._decode_page_stage_next(resp.page_stage_next)
            if not next_stage or next_stage == page_stage:
                break
            page_stage = next_stage

    def paginate_daily_summary(
        self,
        *,
        company_id: str,
        employee_id: str,
        summary_month: str,
        page_size: int,
        session: attendance_pb2.SessionInfo,
        work_date: Optional[int] = None,
        use_employee_api: bool = False,
        start_page_stage: str = "",
        metadata: Optional[Sequence[tuple[str, str]]] = None
    ) -> Generator[attendance_pb2.DailyAttendanceSummaryInfo, None, None]:
        """
        Generator trả về toàn bộ DailyAttendanceSummaryInfo.
        Nếu work_date được cung cấp -> gọi get_daily_attendance_summary (1 ngày cụ thể).
        Nếu không -> dùng *Employee API* để lấy theo tháng.
        """
        page_stage = start_page_stage

        while True:
            if use_employee_api or work_date is None:
                resp: attendance_pb2.GetDailyAttendanceSummaryOutput = self.get_daily_attendance_summary_employee(
                    company_id=company_id,
                    employee_id=employee_id,
                    summary_month=summary_month,
                    page_size=page_size,
                    page_stage=page_stage,
                    session=session,
                    metadata=metadata,
                )
            else:
                resp = self.get_daily_attendance_summary(
                    company_id=company_id,
                    employee_id=employee_id,
                    summary_month=summary_month,
                    work_date=work_date,
                    page_size=page_size,
                    page_stage=page_stage,
                    session=session,
                    metadata=metadata,
                )

            if not resp:
                break

            for rec in resp.records:
                yield rec

            next_stage = self._decode_page_stage_next(resp.page_stage_next)
            if not next_stage or next_stage == page_stage:
                break
            page_stage = next_stage

    # ------------- Builder Helpers -------------
    @staticmethod
    def build_add_attendance_input(
        *,
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        session: attendance_pb2.SessionInfo,
    ) -> attendance_pb2.AddAttendanceInput:
        return attendance_pb2.AddAttendanceInput(
            company_id=company_id,
            employee_id=employee_id,
            device_id=device_id,
            record_time=record_time,
            verification_method=verification_method,
            verification_score=verification_score,
            face_image_url=face_image_url,
            location_coordinates=location_coordinates,
            session=session,
        )

    @staticmethod
    def build_service_add_attendance_input(
        *,
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        session: attendance_pb2.ServiceSessionInfo,
    ) -> attendance_pb2.ServiceAddBatchAttendanceInput:
        return attendance_pb2.ServiceAddBatchAttendanceInput(
            company_id=company_id,
            employee_id=employee_id,
            device_id=device_id,
            record_time=record_time,
            verification_method=verification_method,
            verification_score=verification_score,
            face_image_url=face_image_url,
            location_coordinates=location_coordinates,
            session=session,
        )

    @staticmethod
    def build_service_session_info(
        *,
        service_name: str,
        service_id: str,
        client_ip: str,
        client_agent: str,
    ) -> attendance_pb2.ServiceSessionInfo:
        return attendance_pb2.ServiceSessionInfo(
            service_name=service_name,
            service_id=service_id,
            client_ip=client_ip,
            client_agent=client_agent,
        )

    # ------------- Iter Generators For Streaming -------------
    @staticmethod
    def iter_add_attendance_records(
        records: Iterable[Dict[str, Any]],
        session: attendance_pb2.SessionInfo,
    ) -> Generator[attendance_pb2.AddAttendanceInput, None, None]:
        """
        records: Iterable[dict] gồm:
            company_id, employee_id, device_id?, record_time,
            verification_method?, verification_score?,
            face_image_url?, location_coordinates?
        """
        for r in records:
            yield attendance_pb2.AddAttendanceInput(
                company_id=r["company_id"],
                employee_id=r["employee_id"],
                device_id=r.get("device_id", ""),
                record_time=r["record_time"],
                verification_method=r.get("verification_method", "face"),
                verification_score=float(r.get("verification_score", 0.0)),
                face_image_url=r.get("face_image_url", ""),
                location_coordinates=r.get("location_coordinates", ""),
                session=session,
            )

    @staticmethod
    def iter_service_add_attendance_records(
        records: Iterable[Dict[str, Any]],
        session: attendance_pb2.ServiceSessionInfo,
    ) -> Generator[attendance_pb2.ServiceAddBatchAttendanceInput, None, None]:
        """
        records: Iterable[dict] như trên (không có service_token theo proto).
        """
        for r in records:
            yield attendance_pb2.ServiceAddBatchAttendanceInput(
                company_id=r["company_id"],
                employee_id=r["employee_id"],
                device_id=r.get("device_id", ""),
                record_time=r["record_time"],
                verification_method=r.get("verification_method", "face"),
                verification_score=float(r.get("verification_score", 0.0)),
                face_image_url=r.get("face_image_url", ""),
                location_coordinates=r.get("location_coordinates", ""),
                session=session,
            )


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
) -> AttendanceClient:
    """Factory tạo AttendanceClient với cấu hình thường dùng."""
    return AttendanceClient(
        target=target,
        timeout=timeout,
        max_retries=max_retries,
        retry_backoff=retry_backoff,
        exponential_backoff=exponential_backoff,
        raise_exception=raise_exception,
        insecure=insecure,
        metadata=metadata,
    )