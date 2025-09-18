# TÀI LIỆU BỔ SUNG HOÀN THIỆN HỆ THỐNG CHẤM CÔNG KHUÔN MẶT

---

## 1. QUY TRÌNH SAO LƯU & PHỤC HỒI (BACKUP & DISASTER RECOVERY)

### 1.1. Chính sách sao lưu (Backup Policy)
- **CSDL quan hệ (PostgreSQL):**  
  - Full backup hàng ngày lúc 02:00 AM (GMT+7), lưu trữ 30 bản gần nhất.
  - Binlog/WAL backup mỗi 1 giờ.
- **Timeseries DB (ScyllaDB/Cassandra):**
  - Snapshot hàng ngày, lưu trữ 7 ngày.
- **File Storage (Minio):**
  - Replication sang S3/remote Minio.
  - Snapshot định kỳ 1 lần/ngày.

### 1.2. Quy trình phục hồi (Restore Procedure)
- Tài liệu hướng dẫn thao tác restore từ backup lên môi trường mới/hiện tại.
- Checklist kiểm tra tính toàn vẹn dữ liệu sau restore.
- Thời gian mục tiêu phục hồi (RTO): < 4h cho toàn hệ thống.

### 1.3. Kịch bản thảm hoạ (Disaster Scenarios)
- Lỗi phần cứng, lỗi DB, mất toàn bộ region.
- Quy trình failover sang site DR (nếu có).
- Kiểm tra backup định kỳ mỗi tuần (test restore).

---

## 2. KẾ HOẠCH KIỂM THỬ (TEST PLAN MẪU)

### 2.1. Phạm vi kiểm thử
- Đăng nhập/đăng xuất các vai trò (Admin, Employee, Device).
- Chấm công thành công/thất bại.
- Thao tác CRUD nhân viên, thiết bị, ca làm việc.
- Báo cáo, xuất file, lọc/sort/search.
- Quyền truy cập/phân quyền.
- Tải lên/xóa dữ liệu khuôn mặt.
- Realtime notification (WebSocket).

### 2.2. Test Case mẫu
| Test Case ID | Chức năng | Bước thực hiện | Kết quả kỳ vọng |
| ------------ | --------- | -------------- | --------------- |
| TC-001 | Đăng nhập Admin | Nhập đúng email/mật khẩu | Login thành công, chuyển Dashboard |
| TC-002 | Đăng nhập Admin | Nhập sai mật khẩu | Hiện lỗi, không login |
| TC-010 | Chấm công | Đưa mặt đúng vào camera | Hiện thông báo thành công, ghi nhận log |
| TC-011 | Chấm công | Đưa mặt lạ/không đăng ký | Báo lỗi "khuôn mặt chưa đăng ký" |
| TC-020 | Export báo cáo | Click Export Excel trên màn báo cáo | File .xlsx tải về đúng dữ liệu |
| ... | ... | ... | ... |

### 2.3. Kiểm thử hiệu năng (Performance)
- Đo thời gian phản hồi từ lúc quét mặt đến khi có kết quả (≤2s).
- Kiểm thử tải đồng thời (concurrent check-in).
- Test xuất báo cáo với dữ liệu lớn.

### 2.4. Kiểm thử bảo mật (Security)
- Thử brute force login.
- SQL Injection, XSS trên các form.
- Kiểm tra bảo vệ API (token expiry, role-based access).

---

## 3. MONITORING & ALERTING

### 3.1. Hệ thống giám sát (Monitoring)
- **Mục tiêu:** Phát hiện sớm lỗi, đo hiệu suất, cảnh báo downtime.
- **Công cụ đề xuất:** Prometheus + Grafana, Loki, Alertmanager.
- **Các metric cần giám sát:**
  - Uptime các service (Auth, Attendance, Device, Analytics…)
  - DB (PostgreSQL, ScyllaDB): Connection, query time, disk usage
  - Minio: Storage usage, latency
  - API Gateway: QPS, latency, error rate
  - Device offline rate

### 3.2. Quy tắc cảnh báo (Alert Rules)
- **CPU/RAM > 80%**: Warning alert
- **DB conn > 90%**: Critical alert
- **Kafka lag > 60s**: Warning
- **Device offline > 5 phút**: Gửi email + Push notification admin
- **API error rate > 5% trong 5 phút**: Trigger alert

### 3.3. Quy trình xử lý sự cố (Incident Response)
- Nhận alert → Xác nhận → Điều tra log → Gửi thông báo khách hàng nếu sự cố nghiêm trọng
- Escalate lên DevOps nếu không thể tự xử lý

---

## 4. DEPLOYMENT & INFRASTRUCTURE CHECKLIST

### 4.1. Chuẩn bị môi trường
- Đảm bảo các secret/config đều được quản lý qua vault hoặc secret manager.
- Tách biệt môi trường DEV/UAT/PROD.

### 4.2. CI/CD Pipeline
- Build, test, lint, scan image tự động.
- Triển khai qua blue-green hoặc rolling update.
- Health check endpoint trên tất cả service.

### 4.3. Đảm bảo High Availability
- PostgreSQL: Cluster/replica, auto failover
- ScyllaDB/Cassandra: Cluster multi-node
- Redis: Cluster mode, sentinel
- Minio: Distributed mode

### 4.4. Zero-downtime
- Hỗ trợ hot reload schema (DB migration an toàn).
- API Gateway có thể reload config không downtime.

---
