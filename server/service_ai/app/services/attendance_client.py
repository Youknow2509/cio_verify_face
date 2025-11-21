
from __future__ import annotations

import logging
from typing import Optional, Iterator, Iterable, Callable, List, Generator

import grpc

from app.grpc_generated import attendance_pb2, attendance_pb2_grpc
from app.core.config import settings

_LOGGER = logging.getLogger(__name__)


class AttendanceClient:
    """Client đơn giản cho AttendanceService gRPC."""

    def __init__(
        self,
        target: Optional[str] = None,
        *,
        timeout: float = 5.0,
        max_retries: int = 0,
        retry_backoff: float = 0.5
    ):
        """
        Args:
            target: host:port (VD: 'localhost:50052')
            timeout: timeout cho mỗi RPC
            max_retries: số lần retry khi RpcError (Non-OK)
            retry_backoff: thời gian ngủ giữa các lần retry (giây)
        """
        self._target = target or getattr(settings, 'GRPC_ATTENDANCE_URL', 'localhost:50052')
        self._timeout = timeout
        self._max_retries = max_retries
        self._retry_backoff = retry_backoff
        self._channel = grpc.insecure_channel(self._target)
        self._stub = attendance_pb2_grpc.AttendanceServiceStub(self._channel)

    # ------------- Context Manager -------------
    def __enter__(self) -> "AttendanceClient":
        return self

    def __exit__(self, exc_type, exc, tb):
        self.close()

    # ------------- Internal Helpers -------------
    def _call_with_retry(self, func: Callable, *args, **kwargs):
        """Thực thi RPC với retry nếu cấu hình (max_retries > 0)."""
        attempt = 0
        while True:
            try:
                return func(*args, **kwargs)
            except grpc.RpcError as exc:
                code = exc.code()
                if attempt >= self._max_retries or code in (
                    grpc.StatusCode.INVALID_ARGUMENT,
                    grpc.StatusCode.PERMISSION_DENIED,
                    grpc.StatusCode.UNAUTHENTICATED,
                    grpc.StatusCode.NOT_FOUND,
                ):
                    _LOGGER.error("RPC failed (attempt %d): %s (%s)", attempt + 1, exc, code)
                    return None
                attempt += 1
                _LOGGER.warning("Retrying RPC (attempt %d/%d) code=%s", attempt, self._max_retries, code)
                if self._retry_backoff > 0:
                    import time
                    time.sleep(self._retry_backoff)

    # ------------- Close -------------
    def close(self) -> None:
        try:
            self._channel.close()
        except Exception:
            _LOGGER.exception("Error closing attendance client channel")

    # ------------- Single Attendance -------------
    def add_attendance(
        self,
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        session: attendance_pb2.SessionInfo
    ) -> Optional[attendance_pb2.AddAttendanceOutput]:
        """Ghi một bản ghi điểm danh đơn lẻ."""
        req = attendance_pb2.AddAttendanceInput(
            company_id=company_id,
            employee_id=employee_id,
            device_id=device_id,
            record_time=record_time,
            verification_method=verification_method,
            verification_score=verification_score,
            face_image_url=face_image_url,
            location_coordinates=location_coordinates,
            session=session
        )
        return self._call_with_retry(self._stub.AddAttendance, req, timeout=self._timeout)

    # ------------- Records Pagination -------------
    def get_attendance_records(
        self,
        company_id: str,
        employee_id: str,
        year_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo
    ) -> Optional[attendance_pb2.GetAttendanceRecordsOutput]:
        """Lấy danh sách điểm danh theo tháng (phân trang)."""
        req = attendance_pb2.GetAttendanceRecordsInput(
            company_id=company_id,
            employee_id=employee_id,
            year_month=year_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session
        )
        return self._call_with_retry(self._stub.GetAttendanceRecords, req, timeout=self._timeout)

    def get_attendance_records_employee(
        self,
        company_id: str,
        employee_id: str,
        year_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo
    ) -> Optional[attendance_pb2.GetAttendanceRecordsOutput]:
        """Lấy danh sách điểm danh của một nhân viên (phân trang)."""
        req = attendance_pb2.GetAttendanceRecordsEmployeeInput(
            company_id=company_id,
            employee_id=employee_id,
            year_month=year_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session
        )
        return self._call_with_retry(self._stub.GetAttendanceRecordsEmployee, req, timeout=self._timeout)

    # ------------- Daily Summary -------------
    def get_daily_attendance_summary(
        self,
        company_id: str,
        employee_id: str,
        summary_month: str,
        work_date: int,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo
    ) -> Optional[attendance_pb2.GetDailyAttendanceSummaryOutput]:
        """Lấy tổng hợp điểm danh ngày (phân trang nếu có nhiều block)."""
        req = attendance_pb2.GetDailyAttendanceSummaryInput(
            company_id=company_id,
            employee_id=employee_id,
            summary_month=summary_month,
            work_date=work_date,
            page_size=page_size,
            page_stage=page_stage,
            session=session
        )
        return self._call_with_retry(self._stub.GetDailyAttendanceSummary, req, timeout=self._timeout)

    def get_daily_attendance_summary_employee(
        self,
        company_id: str,
        employee_id: str,
        summary_month: str,
        page_size: int,
        page_stage: str,
        session: attendance_pb2.SessionInfo
    ) -> Optional[attendance_pb2.GetDailyAttendanceSummaryOutput]:
        """Lấy tổng hợp điểm danh theo tháng (của 1 nhân viên)."""
        req = attendance_pb2.GetDailyAttendanceSummaryEmployeeInput(
            company_id=company_id,
            employee_id=employee_id,
            summary_month=summary_month,
            page_size=page_size,
            page_stage=page_stage,
            session=session
        )
        return self._call_with_retry(self._stub.GetDailyAttendanceSummaryEmployee, req, timeout=self._timeout)

    # ------------- Batch Attendance (Streaming) -------------
    def add_batch_attendance(
        self,
        attendance_records: Iterator[attendance_pb2.AddAttendanceInput]
    ) -> Optional[attendance_pb2.AddAttendanceOutput]:
        """
        Ghi nhiều bản ghi điểm danh cùng lúc bằng streaming.
        Mặc định tăng timeout * 10 để tránh timeout sớm.
        """
        return self._call_with_retry(
            self._stub.AddBatchAttendance,
            attendance_records,
            timeout=self._timeout * 10
        )

    def service_add_batch_attendance(
        self,
        attendance_records: Iterator[attendance_pb2.ServiceAddBatchAttendanceInput]
    ) -> Optional[attendance_pb2.AddAttendanceOutput]:
        """
        Ghi nhiều bản ghi dùng service-level authentication (nếu proto định nghĩa).
        """
        return self._call_with_retry(
            self._stub.ServiceAddBatchAttendance,
            attendance_records,
            timeout=self._timeout * 10
        )

    # ------------- Generator Helpers -------------
    @staticmethod
    def build_add_attendance_input(
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        session: attendance_pb2.SessionInfo
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
            session=session
        )

    @staticmethod
    def build_service_add_attendance_input(
        company_id: str,
        employee_id: str,
        device_id: str,
        record_time: int,
        verification_method: str,
        verification_score: float,
        face_image_url: str,
        location_coordinates: str,
        service_token: str,
        session: attendance_pb2.SessionInfo
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
            service_token=service_token,
            session=session
        )

    @staticmethod
    def iter_add_attendance_records(records: Iterable[dict], session: attendance_pb2.SessionInfo
                                    ) -> Generator[attendance_pb2.AddAttendanceInput, None, None]:
        """
        records: Iterable[dict] với các key:
            company_id, employee_id, device_id, record_time,
            verification_method, verification_score,
            face_image_url, location_coordinates
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
                session=session
            )

    @staticmethod
    def iter_service_add_attendance_records(
        records: Iterable[dict],
        session: attendance_pb2.SessionInfo,
        service_token: str
    ) -> Generator[attendance_pb2.ServiceAddBatchAttendanceInput, None, None]:
        """
        records: Iterable[dict] với các key như trên + (tuỳ chọn override verification_method,...)
        service_token: token dành cho service-level auth.
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
                service_token=service_token,
                session=session
            )


def get_client(
    target: Optional[str] = None,
    timeout: float = 5.0,
    *,
    max_retries: int = 0,
    retry_backoff: float = 0.5
) -> AttendanceClient:
    """Factory tạo AttendanceClient."""
    return AttendanceClient(
        target=target,
        timeout=timeout,
        max_retries=max_retries,
        retry_backoff=retry_backoff
    )