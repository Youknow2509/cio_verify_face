from __future__ import annotations

import threading
import time
import logging
from typing import List, Optional, Sequence, Tuple, Callable

from app.grpc_generated import attendance_pb2
from app.grpc.client.attendance_client import AttendanceClient
from app.models.attendance_batching_models import RawAttendanceRecord

_LOGGER = logging.getLogger(__name__)


class AttendanceBatchingService:
    """
    Dịch vụ gom và ghi batch attendance nền.

    Chức năng:
    - enqueue_record(r): đưa record vào hàng đợi batch.
    - send_immediate(r, session_user): ghi ngay lập tức bằng add_attendance.
    - flush(): cưỡng bức flush hiện tại.
    - close(): dừng luồng nền + flush phần còn lại.

    Điều kiện flush tự động:
    - Đạt max_batch_size
    - Hoặc quá flush_interval giây từ lần flush gần nhất.
    """

    def __init__(
        self,
        client: AttendanceClient,
        service_session: attendance_pb2.ServiceSessionInfo,
        *,
        max_batch_size: int = 500,
        flush_interval: float = 5.0,
        max_pending_records: Optional[int] = None,
        metadata: Optional[Sequence[Tuple[str, str]]] = None,
        on_before_flush: Optional[Callable[[int], None]] = None,
        on_after_flush: Optional[Callable[[int, Optional[int], Optional[str]], None]] = None,
        raise_on_immediate_error: bool = False,
    ):
        """
        Args:
            client: AttendanceClient dùng chung.
            service_session: ServiceSessionInfo cho batch streaming.
            max_batch_size: số lượng kỷ lục để flush ngay khi đạt.
            flush_interval: thời gian tối đa giữa 2 lần flush (dù chưa đủ size).
            max_pending_records: giới hạn buffer (None = không giới hạn).
            metadata: metadata truyền khi flush batch.
            on_before_flush: callback(count_in_batch)
            on_after_flush: callback(count_in_batch, status_code, message)
            raise_on_immediate_error: nếu True thì send_immediate ném exception nếu RPC lỗi.
        """
        self._client = client
        self._service_session = service_session
        self._max_batch_size = max_batch_size
        self._flush_interval = flush_interval
        self._max_pending_records = max_pending_records
        self._metadata = metadata
        self._on_before_flush = on_before_flush
        self._on_after_flush = on_after_flush
        self._raise_on_immediate_error = raise_on_immediate_error

        self._buffer: List[RawAttendanceRecord] = []
        self._lock = threading.Lock()
        self._cond = threading.Condition(self._lock)
        self._last_flush_time = time.monotonic()
        self._running = True
        self._thread = threading.Thread(target=self._worker_loop, name="AttendanceBatchWorker", daemon=True)
        self._thread.start()

        self._total_flushed_batches = 0
        self._total_flushed_records = 0
        self._total_immediate = 0
        self._total_immediate_fail = 0

    # ---------------- Public API ----------------
    def enqueue_record(self, record: RawAttendanceRecord) -> bool:
        """
        Đưa record vào buffer.
        Returns:
            True nếu thành công, False nếu bị từ chối (full).
        """
        with self._lock:
            if not self._running:
                _LOGGER.warning("enqueue_record: service stopped")
                return False
            if self._max_pending_records is not None and len(self._buffer) >= self._max_pending_records:
                _LOGGER.error("enqueue_record: buffer full (%d)", len(self._buffer))
                return False
            self._buffer.append(record)
            # Notify worker để kiểm tra điều kiện flush.
            self._cond.notify()
            return True

    def send_immediate(
        self,
        record: RawAttendanceRecord,
        session_user: attendance_pb2.SessionInfo,
        *,
        metadata: Optional[Sequence[Tuple[str, str]]] = None
    ):
        """
        Ghi ngay lập tức một bản ghi bằng add_attendance (dành cho tác vụ ưu tiên / real-time).
        """
        resp = self._client.add_attendance(
            company_id=record.company_id,
            employee_id=record.employee_id,
            device_id=record.device_id,
            record_time=record.record_time,
            verification_method=record.verification_method,
            verification_score=record.verification_score,
            face_image_url=record.face_image_url,
            location_coordinates=record.location_coordinates,
            session=session_user,
            metadata=metadata,
        )
        self._total_immediate += 1
        if not resp or resp.status_code != 200:
            self._total_immediate_fail += 1
            _LOGGER.error("send_immediate failed record=%s status=%s", record, getattr(resp, "status_code", None))
            if self._raise_on_immediate_error:
                raise RuntimeError(f"Immediate write failed status={getattr(resp, 'status_code', None)}")
        return resp

    def flush(self) -> None:
        """
        Cưỡng bức flush hiện tại nếu buffer có dữ liệu.
        """
        self._internal_flush(force=True)

    def stats(self) -> dict:
        """
        Trả về thống kê đơn giản.
        """
        with self._lock:
            return {
                "pending_buffer": len(self._buffer),
                "total_flushed_batches": self._total_flushed_batches,
                "total_flushed_records": self._total_flushed_records,
                "total_immediate": self._total_immediate,
                "total_immediate_fail": self._total_immediate_fail,
            }

    def close(self, *, flush_final: bool = True, timeout: Optional[float] = None):
        """
        Dừng dịch vụ. Option flush phần còn lại trước khi thoát.
        """
        with self._lock:
            self._running = False
            self._cond.notify_all()

        if timeout is not None:
            deadline = time.time() + timeout
        else:
            deadline = None

        # Join worker
        while self._thread.is_alive():
            self._thread.join(timeout=0.2)
            if deadline and time.time() > deadline:
                _LOGGER.warning("close: timeout waiting worker join")
                break

        # Final flush nếu cần và worker chưa làm
        if flush_final:
            self._internal_flush(force=True)

    # ---------------- Worker Loop ----------------
    def _worker_loop(self):
        _LOGGER.info("AttendanceBatchingService worker started")
        while True:
            with self._lock:
                if not self._running and not self._buffer:
                    break

                # Chờ cho tới khi có bản ghi hoặc hết interval
                wait_time = max(self._flush_interval - (time.monotonic() - self._last_flush_time), 0.0)
                # Nếu đã đủ batch size thì flush ngay không chờ
                if len(self._buffer) >= self._max_batch_size:
                    self._internal_flush()
                    continue

                notified = self._cond.wait(timeout=wait_time)

                # Sau wait: kiểm tra điều kiện
                if len(self._buffer) >= self._max_batch_size:
                    self._internal_flush()
                elif (time.monotonic() - self._last_flush_time) >= self._flush_interval and self._buffer:
                    self._internal_flush()

                # Loop tiếp để xem trạng thái _running

        _LOGGER.info("AttendanceBatchingService worker exiting")

    # ---------------- Internal Flush ----------------
    def _internal_flush(self, force: bool = False):
        """
        Thực tế thực thi flush (Phải gọi dưới lock).
        """
        if not self._buffer:
            return
        batch = self._buffer
        self._buffer = []  # swap nhanh để giảm lock thời gian dài
        count = len(batch)
        self._last_flush_time = time.monotonic()

        if self._on_before_flush:
            try:
                self._on_before_flush(count)
            except Exception:
                _LOGGER.exception("on_before_flush callback error")

        try:
            iterator = self._build_service_iterator(batch)
            _LOGGER.debug("Calling service_add_batch_attendance with %d records", count)
            resp = self._client.service_add_batch_attendance(iterator, metadata=self._metadata)
            
            if resp is None:
                _LOGGER.error("service_add_batch_attendance returned None - possible connection issue")
                # Đưa batch lại vào buffer để retry (tuỳ chọn)
                # with self._lock:
                #     self._buffer = batch + self._buffer
                status_code = None
                message = "Response is None - connection/network error"
            else:
                status_code = resp.status_code
                message = resp.message
                
                if status_code in (0, 300):
                    self._total_flushed_batches += 1
                    self._total_flushed_records += count
                    _LOGGER.info("Flushed batch count=%d status=%s", count, status_code)
                else:
                    _LOGGER.error("Flush failed count=%d status=%s message=%s", count, status_code, message)
                    
        except Exception as e:
            _LOGGER.exception("Exception during flush count=%d error=%s", count, str(e))
            status_code = None
            message = f"Exception: {type(e).__name__}: {str(e)}"

        if self._on_after_flush:
            try:
                self._on_after_flush(count, status_code, message)
            except Exception:
                _LOGGER.exception("on_after_flush callback error")

    # ---------------- Iterator builder ----------------
    def _build_service_iterator(self, batch: List[RawAttendanceRecord]):
        return (self._to_service_message(r) for r in batch)

    def _to_service_message(self, r: RawAttendanceRecord) -> attendance_pb2.ServiceAddBatchAttendanceInput:
        """Convert RawAttendanceRecord to protobuf message"""
        try:
            # Convert record_time to Unix timestamp (int64)
            if isinstance(r.record_time, int):
                record_time_ts = r.record_time
            elif isinstance(r.record_time, str):
                from datetime import datetime
                dt = datetime.fromisoformat(r.record_time.replace('Z', '+00:00'))
                record_time_ts = int(dt.timestamp())
            elif hasattr(r.record_time, 'timestamp'):
                record_time_ts = int(r.record_time.timestamp())
            else:
                raise ValueError(f"Invalid record_time type: {type(r.record_time)}")
            
            _LOGGER.debug(
                "Converting record: employee_id=%s time=%s (ts=%d)",
                r.employee_id, r.record_time, record_time_ts
            )
            
            return attendance_pb2.ServiceAddBatchAttendanceInput(
                company_id=str(r.company_id),
                employee_id=str(r.employee_id),
                device_id=str(r.device_id),
                record_time=record_time_ts,  # int64 Unix timestamp
                verification_method=str(r.verification_method),
                verification_score=float(r.verification_score),
                face_image_url=str(r.face_image_url or ""),
                location_coordinates=str(r.location_coordinates or ""),
                session=self._service_session,
            )
        except Exception as e:
            _LOGGER.error(
                "Error converting record: %s\nRecord data: %s",
                str(e), r, exc_info=True
            )
            raise        
# get client 
def get_client() -> AttendanceBatchingService:
    from app.main import app
    return getattr(app.state, "attendance_batching_service", None)