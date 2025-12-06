# API Documentation - Service Analytic

## Tổng quan

Service Analytic cung cấp các API để quản lý và truy xuất dữ liệu phân tích chấm công, báo cáo tổng hợp, và thống kê cho hệ thống chấm công nhận diện khuôn mặt.

**Base URL:** `/api/v1`

**Authentication:** Tất cả các endpoints yêu cầu Bearer Token trong header `Authorization`

---

## Phân quyền (Authorization)

Hệ thống có 3 loại role:

-   **SystemAdmin**: Toàn quyền truy cập tất cả dữ liệu
-   **CompanyAdmin**: Truy cập dữ liệu công ty của mình
-   **Employee**: Chỉ truy cập dữ liệu cá nhân của mình

---

## 1. Health Check

### GET `/health`

**Mô tả:** Kiểm tra trạng thái hoạt động của service

**Input:** Không

**Output:**

```json
{
    "success": true,
    "data": {
        "status": "healthy",
        "timestamp": "2024-01-15T10:30:00Z"
    }
}
```

**Sử dụng:** Monitoring, health check cho load balancer

**Phân quyền:** Không yêu cầu authentication

---

## 2. Reports - Báo cáo

### 2.1. GET `/reports/daily`

**Mô tả:** Lấy báo cáo chấm công chi tiết theo ngày

**Input (Query Parameters):**

-   `date` (required): Ngày báo cáo, format `YYYY-MM-DD` (ví dụ: `2024-01-15`)
-   `company_id` (required): UUID của công ty
-   `device_id` (optional): UUID của thiết bị chấm công

**Output:**

```json
{
  "success": true,
  "data": {
    "date": "2024-01-15",
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_employees": 100,
    "present": 95,
    "absent": 5,
    "late": 10,
    "early_leave": 3,
    "records": [...]
  }
}
```

**Sử dụng:** Xem báo cáo chấm công hàng ngày của công ty, theo dõi tình hình đi làm của nhân viên

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 2.2. GET `/reports/summary`

**Mô tả:** Lấy báo cáo tổng hợp theo tháng với phân tích theo tuần

**Input (Query Parameters):**

-   `month` (required): Tháng báo cáo, format `YYYY-MM` (ví dụ: `2024-01`)
-   `company_id` (required): UUID của công ty

**Output:**

```json
{
    "success": true,
    "data": {
        "month": "2024-01",
        "company_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_working_days": 22,
        "weekly_breakdown": [
            {
                "week": 1,
                "present_rate": 95.5,
                "late_rate": 10.2
            }
        ],
        "summary": {
            "total_present": 2090,
            "total_absent": 110,
            "average_attendance_rate": 95.0
        }
    }
}
```

**Sử dụng:** Xem báo cáo tổng hợp tháng, phân tích xu hướng chấm công theo tuần

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 2.3. POST `/reports/export`

**Mô tả:** Export báo cáo theo khoảng thời gian sang file (Excel/PDF/CSV)

**Input (Request Body):**

```json
{
    "start_date": "2024-01-01",
    "end_date": "2024-01-31",
    "format": "excel",
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@example.com"
}
```

**Các trường:**

-   `start_date` (required): Ngày bắt đầu `YYYY-MM-DD`
-   `end_date` (required): Ngày kết thúc `YYYY-MM-DD`
-   `format` (required): Định dạng file, giá trị: `excel`, `pdf`, `csv`
-   `company_id` (required): UUID công ty
-   `email` (optional): Email để gửi file, nếu không có sẽ trả về download link

**Output:**

```json
{
    "success": true,
    "data": {
        "file_id": "report_20240115_103000.xlsx",
        "download_url": "/api/v1/reports/download/report_20240115_103000.xlsx",
        "expires_at": "2024-01-16T10:30:00Z"
    }
}
```

**Sử dụng:** Export báo cáo để lưu trữ, gửi email hoặc in ấn

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 2.4. POST `/reports/daily/export`

**Mô tả:** Export báo cáo chi tiết theo ngày

**Input (Request Body):**

```json
{
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "date": "2024-01-15",
    "format": "excel",
    "email": "admin@example.com"
}
```

**Output:** Tương tự 2.3

**Sử dụng:** Export báo cáo ngày cụ thể với thông tin chi tiết từng nhân viên

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 2.5. GET `/reports/download/:filename`

**Mô tả:** Download file báo cáo đã được tạo

**Input (Path Parameter):**

-   `filename`: Tên file cần download (từ response của các API export)

**Output:** File binary (Excel/PDF/CSV)

**Sử dụng:** Download file báo cáo đã được export

**Phân quyền:** Authenticated user

---

## 3. Attendance Records - Bản ghi chấm công

### 3.1. GET `/attendance-records`

**Mô tả:** Lấy danh sách bản ghi chấm công theo công ty và tháng

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `year_month` (required): Tháng năm, format `YYYY-MM`
-   `limit` (optional): Số lượng bản ghi, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "record_id": "uuid",
            "employee_id": "uuid",
            "employee_name": "Nguyễn Văn A",
            "check_in_time": "2024-01-15T08:30:00Z",
            "check_out_time": "2024-01-15T17:30:00Z",
            "device_id": "uuid",
            "location": "Main Office",
            "status": "on_time"
        }
    ]
}
```

**Sử dụng:** Xem lịch sử chấm công của công ty theo tháng

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 3.2. GET `/attendance-records/range`

**Mô tả:** Lấy bản ghi chấm công trong khoảng thời gian cụ thể

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `year_month` (required): Tháng năm `YYYY-MM`
-   `start_time` (required): Thời gian bắt đầu, RFC3339 format
-   `end_time` (required): Thời gian kết thúc, RFC3339 format

**Output:** Tương tự 3.1

**Sử dụng:** Query chấm công trong khoảng thời gian cụ thể (ví dụ: từ 8h-10h sáng)

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 3.3. GET `/attendance-records/employee/:employee_id`

**Mô tả:** Lấy bản ghi chấm công của một nhân viên cụ thể

**Input:**

-   Path: `employee_id` (required): UUID nhân viên
-   Query: `company_id` (required): UUID công ty
-   Query: `year_month` (required): Tháng năm `YYYY-MM`
-   Query: `limit` (optional): Số lượng bản ghi, mặc định 100

**Output:** Tương tự 3.1

**Sử dụng:** Xem lịch sử chấm công của một nhân viên cụ thể

**Phân quyền:**

-   CompanyAdmin: Xem nhân viên trong công ty của mình
-   Employee: Chỉ xem của chính mình
-   SystemAdmin: Tất cả

---

### 3.4. GET `/attendance-records/user/:employee_id`

**Mô tả:** Lấy bản ghi chấm công theo user (optimized query từ bảng by_user)

**Input:** Tương tự 3.3

**Output:** Tương tự 3.1

**Sử dụng:** Truy vấn nhanh hơn cho dữ liệu của user cụ thể (dùng partition key tối ưu)

**Phân quyền:** Tương tự 3.3

---

### 3.5. GET `/attendance-records-no-shift`

**Mô tả:** Lấy các bản ghi chấm công không có ca làm việc

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `year_month` (required): Tháng năm `YYYY-MM`
-   `limit` (optional): Số lượng bản ghi, mặc định 100

**Output:** Tương tự 3.1

**Sử dụng:** Phát hiện các bản ghi chấm công bất thường (chấm công ngoài ca)

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

## 4. Daily Summaries - Tổng hợp theo ngày

### 4.1. GET `/daily-summaries`

**Mô tả:** Lấy tổng hợp chấm công theo ngày

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `date` (required): Ngày, format `YYYY-MM-DD`
-   `limit` (optional): Số lượng, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "employee_id": "uuid",
            "employee_name": "Nguyễn Văn A",
            "date": "2024-01-15",
            "total_hours": 8.5,
            "check_in": "08:30:00",
            "check_out": "17:30:00",
            "late_minutes": 0,
            "early_leave_minutes": 0,
            "status": "present"
        }
    ]
}
```

**Sử dụng:** Xem tổng hợp nhanh tình trạng chấm công trong ngày

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 4.2. POST `/daily-summaries/details`

**Mô tả:** Lấy chi tiết báo cáo ngày với pagination

**Input (Request Body):**

```json
{
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "page_size": 50,
    "page_state": "base64_encoded_page_token"
}
```

**Output:**

```json
{
  "success": true,
  "data": {
    "items": [...],
    "page_state": "next_page_token",
    "has_more": true
  }
}
```

**Sử dụng:** Lấy chi tiết báo cáo với pagination cho dataset lớn

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 4.3. GET `/daily-summaries/user/:employee_id`

**Mô tả:** Lấy tổng hợp theo ngày của một user cụ thể

**Input:**

-   Path: `employee_id` (required): UUID nhân viên
-   Query: `company_id` (required): UUID công ty
-   Query: `year_month` (required): Tháng năm `YYYY-MM`
-   Query: `limit` (optional): Số lượng, mặc định 100

**Output:** Tương tự 4.1

**Sử dụng:** Xem tổng hợp chấm công của một nhân viên trong tháng

**Phân quyền:**

-   CompanyAdmin: Xem nhân viên trong công ty
-   Employee: Chỉ xem của chính mình
-   SystemAdmin: Tất cả

---

## 5. Audit Logs - Nhật ký hệ thống

### 5.1. GET `/audit-logs`

**Mô tả:** Lấy nhật ký audit

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `limit` (optional): Số lượng, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "log_id": "uuid",
            "timestamp": "2024-01-15T10:30:00Z",
            "user_id": "uuid",
            "action": "UPDATE_EMPLOYEE",
            "resource_type": "employee",
            "resource_id": "uuid",
            "details": "Updated employee information",
            "ip_address": "192.168.1.100"
        }
    ]
}
```

**Sử dụng:** Theo dõi các hành động trong hệ thống, audit trail

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 5.2. GET `/audit-logs/range`

**Mô tả:** Lấy audit logs trong khoảng thời gian

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `start_time` (required): Thời gian bắt đầu, RFC3339
-   `end_time` (required): Thời gian kết thúc, RFC3339
-   `limit` (optional): Số lượng, mặc định 100

**Output:** Tương tự 5.1

**Sử dụng:** Tìm kiếm audit logs trong khoảng thời gian cụ thể

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 5.3. POST `/audit-logs`

**Mô tả:** Tạo audit log mới

**Input (Request Body):**

```json
{
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "uuid",
    "action": "DELETE_EMPLOYEE",
    "resource_type": "employee",
    "resource_id": "uuid",
    "details": "Deleted employee record",
    "ip_address": "192.168.1.100"
}
```

**Output:**

```json
{
    "success": true,
    "data": {
        "log_id": "generated_uuid",
        "created_at": "2024-01-15T10:30:00Z"
    }
}
```

**Sử dụng:** Ghi lại các hành động quan trọng trong hệ thống

**Phân quyền:** Authenticated users (ghi log của chính mình)

---

## 6. Face Enrollment Logs - Nhật ký đăng ký khuôn mặt

### 6.1. GET `/face-enrollment-logs`

**Mô tả:** Lấy nhật ký đăng ký khuôn mặt

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `limit` (optional): Số lượng, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "log_id": "uuid",
            "employee_id": "uuid",
            "employee_name": "Nguyễn Văn A",
            "enrolled_at": "2024-01-15T10:30:00Z",
            "status": "success",
            "quality_score": 95.5,
            "enrolled_by": "admin_user_id"
        }
    ]
}
```

**Sử dụng:** Theo dõi quá trình đăng ký khuôn mặt của nhân viên

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 6.2. GET `/face-enrollment-logs/employee/:employee_id`

**Mô tả:** Lấy lịch sử đăng ký khuôn mặt của một nhân viên

**Input:**

-   Path: `employee_id` (required): UUID nhân viên
-   Query: `company_id` (required): UUID công ty
-   Query: `limit` (optional): Số lượng, mặc định 100

**Output:** Tương tự 6.1

**Sử dụng:** Xem lịch sử đăng ký và cập nhật khuôn mặt của một nhân viên

**Phân quyền:**

-   CompanyAdmin: Xem nhân viên trong công ty
-   Employee: Chỉ xem của chính mình
-   SystemAdmin: Tất cả

---

## 7. Company Admin - Advanced Analytics

Các endpoints dành cho quản lý công ty để xem báo cáo nâng cao.

### 7.1. GET `/company/daily-attendance-status`

**Mô tả:** Lấy trạng thái chấm công chi tiết theo ngày (bao gồm thống kê về đúng giờ, muộn, về sớm, v.v.)

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `date` (required): Ngày, format `YYYY-MM-DD`

**Output:**

```json
{
    "success": true,
    "data": {
        "date": "2024-01-15",
        "statistics": {
            "total_employees": 100,
            "checked_in": 95,
            "not_checked_in": 5,
            "on_time": 85,
            "late": 10,
            "early_leave": 3,
            "overtime": 15
        },
        "employees": [
            {
                "employee_id": "uuid",
                "name": "Nguyễn Văn A",
                "check_in": "08:30:00",
                "check_out": "17:30:00",
                "status": "on_time",
                "late_minutes": 0,
                "total_hours": 8.5
            }
        ]
    }
}
```

**Sử dụng:** Dashboard quản lý xem tổng quan chi tiết tình trạng chấm công trong ngày

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 7.2. GET `/company/attendance-status/range`

**Mô tả:** Lấy trạng thái chấm công trong khoảng thời gian

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `start_date` (required): Ngày bắt đầu `YYYY-MM-DD`
-   `end_date` (required): Ngày kết thúc `YYYY-MM-DD`

**Output:**

```json
{
    "success": true,
    "data": {
        "start_date": "2024-01-01",
        "end_date": "2024-01-31",
        "daily_statistics": [
            {
                "date": "2024-01-01",
                "total_employees": 100,
                "present": 95,
                "absent": 5,
                "late": 10
            }
        ],
        "overall_summary": {
            "average_attendance_rate": 95.0,
            "average_late_rate": 10.5,
            "total_working_days": 22
        }
    }
}
```

**Sử dụng:** Phân tích xu hướng chấm công trong một khoảng thời gian

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 7.3. GET `/company/monthly-summary`

**Mô tả:** Lấy báo cáo tổng hợp chi tiết theo tháng

**Input (Query Parameters):**

-   `company_id` (required): UUID công ty
-   `month` (required): Tháng, format `YYYY-MM`

**Output:**

```json
{
    "success": true,
    "data": {
        "month": "2024-01",
        "total_working_days": 22,
        "employees_summary": [
            {
                "employee_id": "uuid",
                "name": "Nguyễn Văn A",
                "present_days": 20,
                "absent_days": 2,
                "late_days": 3,
                "total_hours": 176.0,
                "average_hours_per_day": 8.8,
                "overtime_hours": 10.0
            }
        ],
        "statistics": {
            "average_attendance_rate": 95.0,
            "total_late_instances": 150,
            "total_overtime_hours": 500.0
        }
    }
}
```

**Sử dụng:** Báo cáo tổng hợp tháng với thống kê chi tiết từng nhân viên

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 7.4. POST `/company/export-daily-status`

**Mô tả:** Export báo cáo trạng thái chấm công theo ngày

**Input (Request Body):**

```json
{
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "date": "2024-01-15",
    "format": "excel",
    "email": "admin@example.com"
}
```

**Output:** Tương tự section 2.3

**Sử dụng:** Export báo cáo trạng thái ngày sang file

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

### 7.5. POST `/company/export-monthly-summary`

**Mô tả:** Export báo cáo tổng hợp tháng

**Input (Request Body):**

```json
{
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "month": "2024-01",
    "format": "excel",
    "email": "admin@example.com"
}
```

**Output:** Tương tự section 2.3

**Sử dụng:** Export báo cáo tháng chi tiết sang file

**Phân quyền:** CompanyAdmin (chỉ công ty của mình), SystemAdmin (tất cả)

---

## 8. Employee Self-Service - Nhân viên tự tra cứu

Các endpoints dành cho nhân viên để xem dữ liệu của chính mình.

### 8.1. GET `/employee/my-attendance-records`

**Mô tả:** Nhân viên xem bản ghi chấm công của mình

**Input (Query Parameters):**

-   `year_month` (optional): Tháng năm `YYYY-MM`, mặc định tháng hiện tại
-   `limit` (optional): Số lượng bản ghi, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "record_id": "uuid",
            "check_in_time": "2024-01-15T08:30:00Z",
            "check_out_time": "2024-01-15T17:30:00Z",
            "device_name": "Main Office - Floor 1",
            "status": "on_time",
            "total_hours": 8.5
        }
    ]
}
```

**Sử dụng:** Nhân viên tra cứu lịch sử chấm công của mình

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.2. GET `/employee/my-attendance-records/range`

**Mô tả:** Nhân viên xem chấm công trong khoảng thời gian

**Input (Query Parameters):**

-   `year_month` (required): Tháng năm `YYYY-MM`
-   `start_time` (required): Thời gian bắt đầu, RFC3339
-   `end_time` (required): Thời gian kết thúc, RFC3339

**Output:** Tương tự 8.1

**Sử dụng:** Nhân viên xem chấm công trong khoảng thời gian cụ thể

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.3. GET `/employee/my-daily-summaries`

**Mô tả:** Nhân viên xem tổng hợp theo ngày của mình

**Input (Query Parameters):**

-   `year_month` (optional): Tháng năm `YYYY-MM`, mặc định tháng hiện tại
-   `limit` (optional): Số lượng, mặc định 100

**Output:**

```json
{
    "success": true,
    "data": [
        {
            "date": "2024-01-15",
            "total_hours": 8.5,
            "check_in": "08:30:00",
            "check_out": "17:30:00",
            "status": "on_time",
            "late_minutes": 0
        }
    ]
}
```

**Sử dụng:** Nhân viên xem tổng hợp ngày làm việc của mình

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.4. GET `/employee/my-daily-summary/:date`

**Mô tả:** Nhân viên xem chi tiết một ngày cụ thể

**Input:**

-   Path: `date` (required): Ngày, format `YYYY-MM-DD`

**Output:**

```json
{
    "success": true,
    "data": {
        "date": "2024-01-15",
        "check_in": "08:30:00",
        "check_out": "17:30:00",
        "total_hours": 8.5,
        "break_time": 1.0,
        "status": "on_time",
        "late_minutes": 0,
        "early_leave_minutes": 0,
        "check_in_location": "Main Office",
        "check_out_location": "Main Office"
    }
}
```

**Sử dụng:** Nhân viên xem chi tiết chấm công một ngày

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.5. GET `/employee/my-stats`

**Mô tả:** Nhân viên xem thống kê chấm công của mình

**Input (Query Parameters):**

-   `month` (optional): Tháng `YYYY-MM`, mặc định tháng hiện tại

**Output:**

```json
{
    "success": true,
    "data": {
        "month": "2024-01",
        "total_working_days": 22,
        "present_days": 20,
        "absent_days": 2,
        "late_days": 3,
        "early_leave_days": 1,
        "total_hours": 168.0,
        "average_hours_per_day": 8.4,
        "attendance_rate": 90.9,
        "on_time_rate": 85.0
    }
}
```

**Sử dụng:** Nhân viên xem thống kê tổng quan của mình trong tháng

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.6. GET `/employee/my-daily-status`

**Mô tả:** Nhân viên xem trạng thái chi tiết của một ngày

**Input (Query Parameters):**

-   `date` (required): Ngày `YYYY-MM-DD`

**Output:**

```json
{
    "success": true,
    "data": {
        "date": "2024-01-15",
        "status": "on_time",
        "check_in_time": "08:30:00",
        "check_out_time": "17:30:00",
        "required_hours": 8.0,
        "actual_hours": 8.5,
        "overtime_hours": 0.5,
        "late_minutes": 0,
        "early_leave_minutes": 0
    }
}
```

**Sử dụng:** Nhân viên xem trạng thái chi tiết làm việc trong ngày

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.7. GET `/employee/my-status/range`

**Mô tả:** Nhân viên xem trạng thái trong khoảng thời gian

**Input (Query Parameters):**

-   `start_date` (required): Ngày bắt đầu `YYYY-MM-DD`
-   `end_date` (required): Ngày kết thúc `YYYY-MM-DD`

**Output:**

```json
{
    "success": true,
    "data": {
        "start_date": "2024-01-01",
        "end_date": "2024-01-31",
        "daily_status": [
            {
                "date": "2024-01-01",
                "status": "on_time",
                "total_hours": 8.5
            }
        ],
        "summary": {
            "total_days": 22,
            "present_days": 20,
            "on_time_days": 17,
            "late_days": 3
        }
    }
}
```

**Sử dụng:** Nhân viên xem tổng quan trong khoảng thời gian

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.8. GET `/employee/my-monthly-summary`

**Mô tả:** Nhân viên xem báo cáo tổng hợp tháng của mình

**Input (Query Parameters):**

-   `month` (optional): Tháng `YYYY-MM`, mặc định tháng hiện tại

**Output:**

```json
{
    "success": true,
    "data": {
        "month": "2024-01",
        "total_working_days": 22,
        "present_days": 20,
        "absent_days": 2,
        "late_days": 3,
        "early_leave_days": 1,
        "total_hours": 168.0,
        "overtime_hours": 10.0,
        "weekly_breakdown": [
            {
                "week": 1,
                "present_days": 5,
                "total_hours": 42.5
            }
        ]
    }
}
```

**Sử dụng:** Nhân viên xem báo cáo tháng chi tiết với phân tích theo tuần

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.9. POST `/employee/export-daily-status`

**Mô tả:** Nhân viên export trạng thái ngày của mình

**Input (Request Body):**

```json
{
    "date": "2024-01-15",
    "format": "excel",
    "email": "employee@example.com"
}
```

**Output:** Tương tự section 2.3

**Sử dụng:** Nhân viên export dữ liệu của mình

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

### 8.10. POST `/employee/export-monthly-summary`

**Mô tả:** Nhân viên export báo cáo tháng của mình

**Input (Request Body):**

```json
{
    "month": "2024-01",
    "format": "excel",
    "email": "employee@example.com"
}
```

**Output:** Tương tự section 2.3

**Sử dụng:** Nhân viên export báo cáo tháng của mình

**Phân quyền:** Employee (chỉ dữ liệu của chính mình)

---

## Error Responses

Tất cả các API đều có thể trả về các error responses sau:

### 400 Bad Request

```json
{
    "success": false,
    "error": {
        "code": "INVALID_INPUT",
        "message": "Invalid query parameters",
        "details": "date field is required"
    }
}
```

### 401 Unauthorized

```json
{
    "success": false,
    "error": {
        "code": "UNAUTHORIZED",
        "message": "Invalid or missing authentication token",
        "details": ""
    }
}
```

### 403 Forbidden

```json
{
    "success": false,
    "error": {
        "code": "FORBIDDEN",
        "message": "You don't have permission to access this resource",
        "details": "Employees can only access their own data"
    }
}
```

### 500 Internal Server Error

```json
{
    "success": false,
    "error": {
        "code": "QUERY_FAILED",
        "message": "Failed to retrieve data",
        "details": "Database connection timeout"
    }
}
```

---

## Các Error Codes thường gặp

| Code            | Mô tả                                 |
| --------------- | ------------------------------------- |
| `INVALID_INPUT` | Dữ liệu đầu vào không hợp lệ          |
| `INVALID_DATE`  | Định dạng ngày tháng không đúng       |
| `INVALID_UUID`  | UUID không hợp lệ                     |
| `UNAUTHORIZED`  | Chưa xác thực hoặc token không hợp lệ |
| `FORBIDDEN`     | Không có quyền truy cập               |
| `NOT_FOUND`     | Không tìm thấy dữ liệu                |
| `QUERY_FAILED`  | Lỗi khi truy vấn database             |
| `EXPORT_FAILED` | Lỗi khi export file                   |

---

## Notes

1. **Pagination**: Các API trả về danh sách lớn hỗ trợ pagination qua `page_state` (base64 encoded token)
2. **Date Format**: Luôn sử dụng `YYYY-MM-DD` cho date và `YYYY-MM` cho month
3. **Time Format**: Sử dụng RFC3339 cho timestamp với timezone
4. **UUID**: Tất cả ID đều là UUID version 4
5. **Rate Limiting**: APIs có thể bị giới hạn số lần request, check header `X-RateLimit-*`
6. **Caching**: Một số endpoints có cache, check header `X-Cache-Status`

---

## Testing với Swagger

Khi chạy ở mode `dev`, có thể truy cập Swagger UI tại:

```
http://localhost:PORT/swagger/index.html
```

---

## Liên hệ

Nếu có vấn đề hoặc câu hỏi, vui lòng liên hệ team phát triển.
