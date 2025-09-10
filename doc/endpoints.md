# API Endpoint Overview for Face Attendance SaaS System

## 1. Auth Service (Xác thực người dùng & thiết bị)
| Method | Endpoint                          | Mô tả                                       |
| ------ | --------------------------------- | ------------------------------------------- |
| POST   | /api/v1/auth/login                | Đăng nhập tài khoản (user/device)           |
| POST   | /api/v1/auth/logout               | Đăng xuất                                   |
| POST   | /api/v1/auth/refresh              | Làm mới access token                        |
| POST   | /api/v1/auth/device/activate      | Kích hoạt thiết bị chấm công                |
| GET    | /api/v1/auth/me                   | Lấy thông tin tài khoản hiện tại            |

---

## 2. Identity & Organization Service (Quản lý công ty, user, nhân viên)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| GET    | /api/v1/companies                         | Danh sách công ty (System admin)                     |
| POST   | /api/v1/companies                         | Tạo mới công ty                                      |
| GET    | /api/v1/companies/{company_id}            | Xem thông tin công ty                                |
| PUT    | /api/v1/companies/{company_id}            | Sửa thông tin công ty                                |
| DELETE | /api/v1/companies/{company_id}            | Xóa hoặc khóa công ty                                |
| GET    | /api/v1/users                             | Danh sách user / nhân viên trong công ty             |
| POST   | /api/v1/users                             | Thêm mới nhân viên                                   |
| GET    | /api/v1/users/{user_id}                   | Xem thông tin nhân viên                              |
| PUT    | /api/v1/users/{user_id}                   | Sửa thông tin nhân viên                              |
| DELETE | /api/v1/users/{user_id}                   | Vô hiệu hóa/xóa nhân viên                            |
| POST   | /api/v1/users/{user_id}/face-data         | Đăng ký thêm ảnh khuôn mặt (upload)                  |
| GET    | /api/v1/users/{user_id}/face-data         | Lấy danh sách ảnh khuôn mặt của nhân viên            |
| DELETE | /api/v1/users/{user_id}/face-data/{fid}   | Xoá ảnh khuôn mặt                                    |

---

## 3. Device Management Service (Quản lý thiết bị)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| GET    | /api/v1/devices                           | Danh sách thiết bị chấm công của công ty             |
| POST   | /api/v1/devices                           | Thêm mới thiết bị                                    |
| GET    | /api/v1/devices/{device_id}               | Xem thông tin thiết bị                               |
| PUT    | /api/v1/devices/{device_id}               | Sửa thông tin thiết bị                               |
| DELETE | /api/v1/devices/{device_id}               | Xóa/vô hiệu hóa thiết bị                             |

---

## 4. Workforce Service (Quản lý ca làm việc, lịch trình)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| GET    | /api/v1/shifts                            | Danh sách ca làm việc                                |
| POST   | /api/v1/shifts                            | Tạo ca làm việc mới                                  |
| GET    | /api/v1/shifts/{shift_id}                 | Xem thông tin ca làm việc                            |
| PUT    | /api/v1/shifts/{shift_id}                 | Sửa ca làm việc                                      |
| DELETE | /api/v1/shifts/{shift_id}                 | Xoá ca làm việc                                      |
| GET    | /api/v1/schedules                         | Lịch trình chấm công (theo nhân viên/thiết bị)       |
| POST   | /api/v1/schedules                         | Phân ca/lịch làm việc                                |
| DELETE | /api/v1/schedules/{schedule_id}           | Xoá lịch trình                                       |

---

## 5. Attendance Service (Chấm công)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| POST   | /api/v1/attendance/check-in               | Chấm công vào bằng khuôn mặt                         |
| POST   | /api/v1/attendance/check-out              | Chấm công ra bằng khuôn mặt                          |
| GET    | /api/v1/attendance/records                | Lấy lịch sử chấm công (theo user/device/ngày)        |
| GET    | /api/v1/attendance/records/{record_id}    | Chi tiết một lượt chấm công                          |
| GET    | /api/v1/attendance/history/my             | Nhân viên xem lịch sử chấm công cá nhân              |

---

## 6. Analytics & Reporting Service (Báo cáo, thống kê)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| GET    | /api/v1/reports/daily                     | Báo cáo chi tiết ngày (theo công ty/địa điểm)        |
| GET    | /api/v1/reports/summary                   | Báo cáo tổng hợp tháng                               |
| GET    | /api/v1/reports/export                    | Xuất báo cáo (Excel/PDF)                             |

---

## 7. Signature Upload Service (Chữ ký)
| Method | Endpoint                                  | Mô tả                                               |
| ------ | ----------------------------------------- | --------------------------------------------------- |
| POST   | /api/v1/signatures                        | Tải lên chữ ký điện tử                               |
| GET    | /api/v1/signatures/{user_id}              | Lấy danh sách chữ ký của nhân viên                   |
| DELETE | /api/v1/signatures/{signature_id}         | Xóa chữ ký                                           |

---

## 8. WebSocket Events (Realtime Notification)
| Channel/Topic                  | Event name                       | Mô tả                                           |
| ------------------------------ | -------------------------------- | ----------------------------------------------- |
| ws://.../ws                    | attendance_result                | Đẩy kết quả chấm công (success/fail) về thiết bị|
| ws://.../ws                    | device_status                    | Trạng thái hoạt động của thiết bị                |
| ws://.../ws                    | admin_alert                      | Cảnh báo cho admin (nếu cần)                    |

---

## 9. System Admin APIs (Quản trị hệ thống)
| Method | Endpoint                          | Mô tả                                       |
| ------ | --------------------------------- | ------------------------------------------- |
| GET    | /api/v1/admin/companies           | Danh sách công ty khách hàng                |
| POST   | /api/v1/admin/companies           | Tạo mới công ty                             |
| PUT    | /api/v1/admin/companies/{id}/lock | Khóa công ty                                |
| DELETE | /api/v1/admin/companies/{id}      | Xoá công ty                                 |
| GET    | /api/v1/admin/audit-logs          | Nhật ký thao tác hệ thống                   |

---

## 10. Common/Utility APIs
| Method | Endpoint                          | Mô tả                                       |
| ------ | --------------------------------- | ------------------------------------------- |
| GET    | /api/v1/ping                      | Kiểm tra trạng thái hệ thống (health check) |
| GET    | /api/v1/time                      | Lấy thời gian server                        |

---
