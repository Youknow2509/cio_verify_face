# ScyllaDB Schema & Query Guide (Dual Write Model)

Tài liệu này mô tả cấu trúc dữ liệu và các câu lệnh CQL chuẩn để thao tác với hệ thống chấm công.
Mô hình sử dụng: **Dual Writes** (Ghi song song vào 2 bảng để tối ưu hoá việc đọc theo Company và theo User).

---

## 1. Attendance Records (Log Chấm Công)

### Bảng A: `attendance_records_by_company`
*Mục đích: Admin xem logs chấm công của toàn công ty theo tháng.*

### Bảng B: `attendance_records_by_user`
*Mục đích: User xem lịch sử chấm công của cá nhân theo tháng.*

### Check-in / Check-out (INSERT)
Sử dụng **BATCH** để đảm bảo dữ liệu được ghi vào cả 2 bảng cùng lúc.

```sql
BEGIN BATCH
    -- 1. Ghi vào bảng Company
    INSERT INTO attendance_records_by_company (
        company_id, year_month, record_time, employee_id, 
        device_id, record_type, verification_method, verification_score, 
        face_image_url, location_coordinates, metadata, sync_status, created_at
    ) VALUES (
        uuid_company, '2023-10', '2023-10-25 08:00:00', uuid_employee, 
        uuid_device, 0, 'FACE', 0.98, 
        'http://minio/img.jpg', '10.7,106.6', {'ip': '1.2.3.4'}, 'synced', toTimestamp(now())
    );

    -- 2. Ghi vào bảng User
    INSERT INTO attendance_records_by_user (
        company_id, employee_id, year_month, record_time, 
        device_id, record_type, verification_method, verification_score, 
        face_image_url, location_coordinates, metadata, sync_status, created_at
    ) VALUES (
        uuid_company, uuid_employee, '2023-10', '2023-10-25 08:00:00', 
        uuid_device, 0, 'FACE', 0.98, 
        'http://minio/img.jpg', '10.7,106.6', {'ip': '1.2.3.4'}, 'synced', toTimestamp(now())
    );
APPLY BATCH;
```

### Lấy dữ liệu (SELECT)

**Admin xem logs tháng 10:**
```sql
SELECT * FROM attendance_records_by_company 
WHERE company_id = uuid_company AND year_month = '2023-10';
```

**User xem logs tháng 10 của mình:**
```sql
SELECT * FROM attendance_records_by_user 
WHERE company_id = uuid_company AND employee_id = uuid_employee AND year_month = '2023-10';
```

### Xóa Logs (DELETE)
*Lưu ý: Hạn chế xóa logs chấm công. Nếu cần, phải xóa cả 2 bảng.*

```sql
BEGIN BATCH
    DELETE FROM attendance_records_by_company 
    WHERE company_id = uuid_company AND year_month = '2023-10' 
    AND record_time = '2023-10-25 08:00:00' AND employee_id = uuid_employee;

    DELETE FROM attendance_records_by_user 
    WHERE company_id = uuid_company AND employee_id = uuid_employee 
    AND year_month = '2023-10' AND record_time = '2023-10-25 08:00:00';
APPLY BATCH;
```

---

## 2. Daily Summaries (Tổng Hợp Ngày)

### Bảng C: `daily_summaries_by_company`
*Mục đích: Dashboard Công ty (Ai vắng, ai muộn hôm nay).*

### Bảng D: `daily_summaries_by_user`
*Mục đích: Bảng công cá nhân (Timesheet tháng).*

### Cập nhật trạng thái cuối ngày (UPDATE/INSERT)
Khi tính toán công xong (ví dụ: lúc Check-out hoặc Cronjob chạy cuối ngày), cập nhật cả 2 bảng.

```sql
BEGIN BATCH
    -- 1. Cập nhật bảng Company
    INSERT INTO daily_summaries_by_company (
        company_id, summary_month, work_date, employee_id, 
        shift_id, actual_check_in, actual_check_out, attendance_status, 
        late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
    ) VALUES (
        uuid_company, '2023-10', '2023-10-25', uuid_employee, 
        uuid_shift, '2023-10-25 08:00:00', '2023-10-25 17:30:00', 1, 
        15, 0, 480, 'Đi muộn do kẹt xe', toTimestamp(now())
    );

    -- 2. Cập nhật bảng User
    INSERT INTO daily_summaries_by_user (
        company_id, employee_id, summary_month, work_date, 
        shift_id, actual_check_in, actual_check_out, attendance_status, 
        late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
    ) VALUES (
        uuid_company, uuid_employee, '2023-10', '2023-10-25', 
        uuid_shift, '2023-10-25 08:00:00', '2023-10-25 17:30:00', 1, 
        15, 0, 480, 'Đi muộn do kẹt xe', toTimestamp(now())
    );
APPLY BATCH;
```

### Lấy dữ liệu (SELECT)

**Dashboard Công ty ngày 25/10 (Ai đi làm, ai vắng?):**
```sql
SELECT * FROM daily_summaries_by_company 
WHERE company_id = uuid_company AND summary_month = '2023-10' AND work_date = '2023-10-25';
```

**Bảng công tháng 10 của User:**
```sql
SELECT * FROM daily_summaries_by_user 
WHERE company_id = uuid_company AND employee_id = uuid_employee AND summary_month = '2023-10';
```

### Sửa ghi chú / Duyệt công (UPDATE)
Nếu HR sửa công tay, cần update cả 2 bảng. ScyllaDB INSERT đè lên bản ghi cũ sẽ hoạt động như UPDATE.

```sql
BEGIN BATCH
    -- Update note bảng Company
    UPDATE daily_summaries_by_company SET notes = 'Đã duyệt phép', attendance_status = 0
    WHERE company_id = uuid_company AND summary_month = '2023-10' 
    AND work_date = '2023-10-25' AND employee_id = uuid_employee;

    -- Update note bảng User
    UPDATE daily_summaries_by_user SET notes = 'Đã duyệt phép', attendance_status = 0
    WHERE company_id = uuid_company AND employee_id = uuid_employee 
    AND summary_month = '2023-10' AND work_date = '2023-10-25';
APPLY BATCH;
```

---

## 3. Face Enrollment Logs (Log Đăng Ký Khuôn Mặt)

### Bảng E: `face_enrollment_logs`
*Mục đích: Theo dõi lịch sử đăng ký, cập nhật khuôn mặt của nhân viên, giúp debug khi user báo lỗi không nhận diện được hoặc đăng ký thất bại.*

### Ghi Log Đăng Ký (INSERT)
Ghi lại kết quả mỗi khi user thực hiện đăng ký khuôn mặt.

INSERT INTO face_enrollment_logs (
    company_id, year_month, created_at, employee_id, 
    action_type, status, image_url, failure_reason, metadata
) VALUES (
    uuid_company, '2023-10', toTimestamp(now()), uuid_employee, 
    'ENROLL', 'SUCCESS', 'http://minio/face_123.jpg', '', {'device': 'mobile'}
);### Truy vấn lịch sử đăng ký (SELECT)
Xem lịch sử đăng ký của công ty trong tháng.

SELECT * FROM face_enrollment_logs 
WHERE company_id = uuid_company AND year_month = '2023-10';---

## 4. Audit Logs (Nhật Ký Hệ Thống)

### Bảng F: `audit_logs`
*Mục đích: Lưu vết các hành động quan trọng của Admin/User trên hệ thống (VD: Thêm nhân viên, Xóa thiết bị, Sửa công).*

### Ghi Audit Log (INSERT)

INSERT INTO audit_logs (
    company_id, year_month, created_at, actor_id, 
    action_category, action_name, resource_type, resource_id, 
    details, ip_address, user_agent, status
) VALUES (
    uuid_company, '2023-10', toTimestamp(now()), uuid_admin, 
    'HR_MANAGEMENT', 'UPDATE_ATTENDANCE', 'ATTENDANCE_RECORD', 'rec_001', 
    {'reason': 'Fixed forgot checkout', 'old_val': 'missing', 'new_val': '17:30'}, 
    '192.168.1.1', 'Mozilla/5.0...', 'SUCCESS'
);### Truy vấn Audit Log (SELECT)
Xem ai đã làm gì trong tháng này.

SELECT * FROM audit_logs 
WHERE company_id = uuid_company AND year_month = '2023-10';