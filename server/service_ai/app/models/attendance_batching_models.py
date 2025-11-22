from dataclasses import dataclass
from typing import Optional, Union
from datetime import datetime


@dataclass
class RawAttendanceRecord:
    """Raw attendance record for batching"""
    company_id: str
    employee_id: str
    device_id: str
    record_time: Union[int, str, datetime]  # Accepts int (timestamp), ISO string, or datetime
    verification_method: str
    verification_score: float
    face_image_url: Optional[str] = None
    location_coordinates: Optional[str] = None