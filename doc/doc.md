
---

### **TÀI LIỆU PHÂN TÍCH YÊU CẦU HỆ THỐNG (SRS)**

*   **Tên dự án:** Hệ thống Chấm công bằng Khuôn mặt Đa nền tảng (SaaS)
*   **Phiên bản:** 1.0

---

### **1. GIỚI THIỆU**

#### **1.1. Mục đích**
Tài liệu này xác định và mô tả chi tiết các yêu cầu chức năng và phi chức năng cho Hệ thống Chấm công bằng Khuôn mặt. Mục tiêu là cung cấp một nền tảng chung, rõ ràng cho các bên liên quan bao gồm đội ngũ phát triển, kiểm thử, quản lý dự án và khách hàng để hiểu rõ về phạm vi và các tính năng của sản phẩm.

#### **1.2. Phạm vi hệ thống**
Hệ thống là một giải pháp Phần mềm như một Dịch vụ (SaaS), cho phép nhiều công ty khách hàng sử dụng chung một hạ tầng máy chủ để quản lý việc chấm công của nhân viên qua nhận diện khuôn mặt tại nhiều địa điểm khác nhau.

*   **Trong phạm vi (In-scope):**
    *   Xác thực và ghi nhận thời gian chấm công (vào/ra) bằng khuôn mặt.
    *   Quản lý dữ liệu nhân viên, ca làm việc, địa điểm cho từng công ty.
    *   Cung cấp báo cáo chấm công, đi trễ, về sớm.
    *   Phân quyền quản trị cho từng công ty và quản trị toàn hệ thống.
    *   Hỗ trợ hoạt động tại nhiều chi nhánh/địa điểm.

*   **Ngoài phạm vi (Out-of-scope):**
    *   Tính lương và tích hợp với hệ thống Nhân sự - Tiền lương (HRM/Payroll).
    *   Quản lý đơn xin nghỉ phép.
    *   Chấm công bằng các phương thức khác (vân tay, thẻ từ).

#### **1.3. Đối tượng sử dụng tài liệu**
*   **Quản lý dự án:** Để lập kế hoạch và theo dõi tiến độ.
*   **Đội ngũ phát triển (Development Team):** Để thiết kế kiến trúc và lập trình các tính năng.
*   **Đội ngũ kiểm thử (QA/QC Team):** Để xây dựng kịch bản kiểm thử (test cases).
*   **Khách hàng/Bên liên quan:** Để xác nhận các yêu cầu và đảm bảo hệ thống đáp ứng đúng nhu cầu kinh doanh.

#### **1.4. Định nghĩa và từ viết tắt**
*   **SaaS (Software as a Service):** Phần mềm như một Dịch vụ.
*   **Multi-tenant (Đa công ty/Đa nền tảng):** Kiến trúc phần mềm mà một phiên bản duy nhất của ứng dụng chạy trên một máy chủ, phục vụ cho nhiều nhóm người dùng (khách hàng) khác nhau.
*   **Thiết bị Chấm công:** Thiết bị vật lý (máy tính bảng, camera chuyên dụng) có khả năng quét khuôn mặt và kết nối internet.

---

### **2. MÔ TẢ TỔNG QUAN**

#### **2.1. Bối cảnh sản phẩm**
Đây là một sản phẩm độc lập, hoạt động theo mô hình SaaS. Các công ty khách hàng sẽ đăng ký tài khoản và trả phí định kỳ để sử dụng dịch vụ. Hệ thống cung cấp giao diện web cho quản trị viên và tương tác với các thiết bị phần cứng đặt tại văn phòng của khách hàng.

#### **2.2. Chức năng chính**
*   **Chấm công:** Nhân viên thực hiện chấm công nhanh chóng, không tiếp xúc.
*   **Quản lý Nhân sự (cơ bản):** Quản trị viên công ty quản lý thông tin nhân viên và dữ liệu khuôn mặt.
*   **Quản lý Tổ chức:** Quản lý các địa điểm, ca làm việc linh hoạt.
*   **Báo cáo và Thống kê:** Cung cấp cái nhìn tổng quan về tình hình tuân thủ giờ giấc của nhân viên.
*   **Quản trị Hệ thống:** Quản lý các tài khoản khách hàng và giám sát hoạt động của toàn bộ hệ thống.

#### **2.3. Đặc điểm người dùng (Actors)**
1.  **Nhân viên (Employee):** Người dùng cuối, chỉ thực hiện chấm công và xem lịch sử cá nhân.
2.  **Quản trị viên Công ty (Company Admin):** Người được ủy quyền bởi công ty khách hàng, chịu trách nhiệm quản lý dữ liệu chấm công của công ty mình.
3.  **Quản trị viên Hệ thống (System Admin):** Nhân sự của nhà cung cấp dịch vụ, có quyền cao nhất, quản lý toàn bộ hệ thống và các khách hàng.

#### **2.4. Ràng buộc chung**
*   Hệ thống phải được xây dựng trên nền tảng web (cho phần quản trị) và có API để các thiết bị chấm công giao tiếp.
*   Dữ liệu của mỗi công ty phải được cách ly và bảo mật tuyệt đối (Tenant Isolation).
*   Hệ thống phải tuân thủ các quy định về bảo vệ dữ liệu cá nhân (ví dụ: GDPR nếu có khách hàng ở Châu Âu), đặc biệt là dữ liệu sinh trắc học.

---

### **3. YÊU CẦU CHỨC NĂNG (Functional Requirements)**

| ID | Module | Yêu cầu | Chi tiết |
| :--- | :--- | :--- | :--- |
| **FR-ATT-001** | Chấm công | **Chấm công bằng khuôn mặt** | Nhân viên có thể chấm công vào/ra bằng cách để thiết bị quét khuôn mặt. |
| **FR-ATT-002** | Chấm công | **Phản hồi tức thì** | Thiết bị phải hiển thị thông báo thành công (kèm tên) hoặc thất bại ngay sau khi quét. |
| **FR-EMP-001** | Quản lý Nhân viên | **Quản lý thông tin nhân viên** | Quản trị viên Công ty có thể Thêm/Sửa/Vô hiệu hóa tài khoản nhân viên. |
| **FR-EMP-002** | Quản lý Nhân viên | **Đăng ký dữ liệu khuôn mặt** | Hệ thống cho phép đăng ký (chụp/tải lên) nhiều ảnh mẫu cho một nhân viên để tăng độ chính xác. |
| **FR-RPT-001** | Báo cáo | **Báo cáo chi tiết ngày** | Hệ thống cho phép xem danh sách tất cả các lượt chấm công trong ngày theo từng địa điểm. |
| **FR-RPT-002** | Báo cáo | **Báo cáo tổng hợp** | Quản trị viên Công ty có thể tạo báo cáo tổng hợp tháng về số giờ làm, số lần đi trễ/về sớm. |
| **FR-RPT-003** | Báo cáo | **Xuất báo cáo** | Báo cáo có thể được xuất ra các định dạng phổ biến như Excel (.xlsx) hoặc PDF. |
| **FR-CMP-001** | Quản lý Công ty | **Cách ly dữ liệu** | Dữ liệu của công ty A (nhân viên, báo cáo) phải hoàn toàn không thể truy cập bởi người dùng của công ty B. |
| **FR-SYS-001** | Quản trị Hệ thống | **Quản lý khách hàng** | Quản trị viên Hệ thống có thể tạo, khóa, hoặc xóa tài khoản của một công ty khách hàng. |

---

### **4. YÊU CẦU PHI CHỨC NĂNG (Non-Functional Requirements)**

| ID | Loại yêu cầu | Yêu cầu |
| :--- | :--- | :--- |
| **NFR-PER-001** | Hiệu năng (Performance) | Thời gian từ lúc quét mặt đến khi có phản hồi trên thiết bị phải dưới 2 giây trong điều kiện mạng ổn định. |
| **NFR-PER-002** | Hiệu năng (Performance) | Hệ thống web quản trị phải tải trang trong vòng 3 giây. |
| **NFR-SEC-001** | Bảo mật (Security) | Mật khẩu người dùng phải được mã hóa (hashed). Dữ liệu sinh trắc học phải được mã hóa cả khi lưu trữ và truyền tải. |
| **NFR-SEC-002** | Bảo mật (Security) | Giao tiếp giữa thiết bị và server phải sử dụng giao thức HTTPS. |
| **NFR-USA-001** | Tính khả dụng (Usability) | Giao diện web phải tương thích với các trình duyệt phổ biến (Chrome, Firefox, Safari). Giao diện phải trực quan, dễ sử dụng. |
| **NFR-REL-001** | Độ tin cậy (Reliability) | Hệ thống phải có độ sẵn sàng (uptime) 99.5%. |
| **NFR-REL-002** | Độ tin cậy (Reliability) | Phải có cơ chế sao lưu (backup) cơ sở dữ liệu hàng ngày. |
| **NFR-SCA-001** | Khả năng mở rộng (Scalability)| Kiến trúc hệ thống phải cho phép mở rộng để phục vụ số lượng lớn công ty và nhân viên mà không cần thiết kế lại từ đầu. |

---