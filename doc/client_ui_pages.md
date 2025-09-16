# TÀI LIỆU MÔ TẢ GIAO DIỆN CLIENT - HỆ THỐNG CHẤM CÔNG KHUÔN MẶT

---

## **TỔNG QUAN HỆ THỐNG GIAO DIỆN**

Hệ thống Face Attendance SaaS bao gồm 3 giao diện client chính:

1. **Web App (Admin)** - Giao diện quản trị dành cho công ty
2. **Device App (IoT)** - Giao diện thiết bị chấm công 
3. **System Admin Interface** - Giao diện quản trị hệ thống

---

## **1. WEB APP (ADMIN) - GIAO DIỆN QUẢN TRỊ CÔNG TY**

### **1.1. Trang Đăng nhập (Login Page)**

**Đường dẫn:** `/login`  
**Chức năng:** Xác thực tài khoản quản trị viên công ty

**Thành phần giao diện:**
- Logo hệ thống
- Form đăng nhập:
  - Email/Tên đăng nhập
  - Mật khẩu  
  - Nút "Đăng nhập"
  - Link "Quên mật khẩu?"
- Lựa chọn ngôn ngữ (Tiếng Việt/English)
- Thông báo lỗi nếu sai thông tin

**API tương ứng:** `POST /api/v1/auth/login`

---

### **1.2. Dashboard - Trang Tổng Quan**

**Đường dẫn:** `/dashboard`  
**Chức năng:** Hiển thị thông tin tổng quan về tình hình chấm công

**Thành phần giao diện:**
- **Header:**
  - Logo công ty
  - Menu điều hướng chính
  - Thông tin tài khoản (avatar, tên, logout)
  - Thông báo realtime

- **Thống kê tổng quan (Cards):**
  - Tổng số nhân viên active
  - Số lượt chấm công hôm nay  
  - Tỷ lệ đi trễ trong tháng
  - Số thiết bị đang hoạt động

- **Biểu đồ:**
  - Biểu đồ cột: Lượt chấm công 7 ngày gần nhất
  - Biểu đồ tròn: Phân bố chấm công theo ca
  - Timeline: Lịch sử chấm công trong ngày

- **Hoạt động gần đây:**
  - Danh sách 10 lượt chấm công mới nhất
  - Trạng thái thiết bị
  - Cảnh báo nếu có

**API tương ứng:** 
- `GET /api/v1/reports/daily`
- `GET /api/v1/attendance/records`

---

### **1.3. Quản lý Nhân viên (Employee Management)**

#### **1.3.1. Danh sách Nhân viên**

**Đường dẫn:** `/employees`  
**Chức năng:** Hiển thị và quản lý danh sách nhân viên

**Thành phần giao diện:**
- **Thanh công cụ:**
  - Nút "Thêm nhân viên mới"
  - Ô tìm kiếm (theo tên, mã NV)
  - Bộ lọc (theo phòng ban, trạng thái)
  - Nút xuất Excel

- **Bảng danh sách:**
  - Avatar + Tên nhân viên
  - Mã nhân viên
  - Email/SĐT
  - Phòng ban/Chức vụ
  - Trạng thái (Active/Inactive)
  - Số ảnh khuôn mặt đã đăng ký
  - Thao tác (Sửa, Xóa, Xem chi tiết)

- **Phân trang:** Previous, Next, số trang

**API tương ứng:** `GET /api/v1/users`

#### **1.3.2. Thêm/Sửa Nhân viên**

**Đường dẫn:** `/employees/add` hoặc `/employees/{id}/edit`  
**Chức năng:** Form nhập/chỉnh sửa thông tin nhân viên

**Thành phần giao diện:**
- **Thông tin cá nhân:**
  - Ảnh đại diện (upload)
  - Họ và tên (*)
  - Mã nhân viên (tự sinh hoặc nhập)
  - Email (*)
  - Số điện thoại
  - Ngày sinh
  - Giới tính

- **Thông tin công việc:**
  - Phòng ban/Bộ phận
  - Chức vụ
  - Ngày vào làm
  - Loại hợp đồng
  - Ca làm việc mặc định

- **Nút thao tác:**
  - "Lưu"
  - "Lưu và thêm mới"
  - "Hủy"

**API tương ứng:** 
- `POST /api/v1/users` (thêm mới)
- `PUT /api/v1/users/{user_id}` (cập nhật)

#### **1.3.3. Đăng ký Khuôn mặt**

**Đường dẫn:** `/employees/{id}/face-data`  
**Chức năng:** Quản lý dữ liệu khuôn mặt của nhân viên

**Thành phần giao diện:**
- **Thông tin nhân viên:** Avatar, tên, mã NV
- **Upload ảnh:**
  - Khu vực drag & drop
  - Nút "Chọn ảnh" (multiple)
  - Hướng dẫn chụp ảnh (các góc khác nhau)
  - Preview ảnh đã chọn

- **Danh sách ảnh đã lưu:**
  - Thumbnail ảnh
  - Ngày upload
  - Chất lượng (Good/Average/Poor)
  - Nút xóa

- **Thao tác:**
  - "Upload ảnh mới"
  - "Xóa ảnh đã chọn"
  - "Quay lại"

**API tương ứng:**
- `POST /api/v1/users/{user_id}/face-data`
- `GET /api/v1/users/{user_id}/face-data`
- `DELETE /api/v1/users/{user_id}/face-data/{fid}`

---

### **1.4. Quản lý Thiết bị (Device Management)**

#### **1.4.1. Danh sách Thiết bị**

**Đường dẫn:** `/devices`  
**Chức năng:** Quản lý các thiết bị chấm công

**Thành phần giao diện:**
- **Thanh công cụ:**
  - Nút "Thêm thiết bị mới"
  - Tìm kiếm theo tên/mã thiết bị
  - Bộ lọc theo địa điểm, trạng thái

- **Card view thiết bị:**
  - Tên thiết bị
  - Mã thiết bị/Serial
  - Địa điểm đặt
  - Trạng thái (Online/Offline/Error)
  - Thời gian online cuối
  - Số lượt chấm công hôm nay
  - Nút "Chi tiết", "Cấu hình"

**API tương ứng:** `GET /api/v1/devices`

#### **1.4.2. Cấu hình Thiết bị**

**Đường dẫn:** `/devices/{id}/config`  
**Chức năng:** Cấu hình và kích hoạt thiết bị

**Thành phần giao diện:**
- **Thông tin cơ bản:**
  - Tên thiết bị
  - Vị trí/Địa điểm
  - Mô tả
  - IP Address
  - MAC Address

- **Cấu hình chấm công:**
  - Cho phép check-in/check-out
  - Thời gian timeout
  - Độ nhạy nhận diện (threshold)
  - Âm thanh thông báo

- **Bảo mật:**
  - Device Token (hiển thị/tạo mới)
  - QR Code kích hoạt
  - Reset thiết bị

**API tương ứng:**
- `GET /api/v1/devices/{device_id}`
- `PUT /api/v1/devices/{device_id}`
- `POST /api/v1/auth/device/activate`

---

### **1.5. Quản lý Ca làm việc (Shift Management)**

**Đường dẫn:** `/shifts`  
**Chức năng:** Quản lý các ca làm việc và phân ca

**Thành phần giao diện:**
- **Danh sách ca làm việc:**
  - Tên ca
  - Giờ bắt đầu - Kết thúc  
  - Thời gian nghỉ trưa
  - Ngày áp dụng (T2-CN)
  - Số nhân viên được phân
  - Thao tác (Sửa, Xóa, Phân ca)

- **Form thêm/sửa ca:**
  - Tên ca làm việc
  - Giờ vào - Giờ ra
  - Giờ nghỉ trưa (từ - đến)
  - Thời gian chấm công hợp lệ
  - Ngày trong tuần áp dụng
  - Mô tả ca

- **Phân ca cho nhân viên:**
  - Danh sách nhân viên (checkbox)
  - Ngày áp dụng từ - đến
  - Loại phân ca (Cố định/Tạm thời)

**API tương ứng:**
- `GET /api/v1/shifts`
- `POST /api/v1/shifts`
- `POST /api/v1/schedules`

---

### **1.6. Báo cáo Chấm công (Attendance Reports)**

#### **1.6.1. Báo cáo Hàng ngày**

**Đường dẫn:** `/reports/daily`  
**Chức năng:** Xem chi tiết chấm công trong ngày

**Thành phần giao diện:**
- **Bộ lọc:**
  - Chọn ngày
  - Chọn địa điểm/thiết bị
  - Chọn ca làm việc
  - Loại báo cáo (Tất cả/Đi trễ/Về sớm/Thiếu chấm công)

- **Bảng dữ liệu:**
  - STT
  - Tên nhân viên
  - Ca làm việc
  - Giờ vào (thực tế vs quy định)
  - Giờ ra (thực tế vs quy định)
  - Tổng giờ làm
  - Trạng thái (Đúng giờ/Trễ/Sớm/Thiếu)
  - Thiết bị chấm công
  - Ghi chú

- **Tổng kết:**
  - Tổng nhân viên có mặt
  - Số lượt đi trễ
  - Số lượt về sớm
  - Tỷ lệ chấm công đúng

- **Xuất báo cáo:** PDF, Excel

**API tương ứng:** `GET /api/v1/reports/daily`

#### **1.6.2. Báo cáo Tổng hợp**

**Đường dẫn:** `/reports/summary`  
**Chức năng:** Báo cáo thống kê theo tháng/quý

**Thành phần giao diện:**
- **Bộ lọc thời gian:**
  - Từ ngày - Đến ngày
  - Chọn nhanh (Tháng này, Tháng trước, Quý này)
  - Nhóm theo (Nhân viên/Phòng ban/Thiết bị)

- **Biểu đồ thống kê:**
  - Biểu đồ cột: Số giờ làm việc theo ngày
  - Biểu đồ tròn: Tỷ lệ đúng giờ vs đi trễ
  - Trend line: Xu hướng chấm công

- **Bảng tổng hợp:**
  - Tên nhân viên/Phòng ban
  - Tổng ngày làm việc
  - Tổng giờ làm việc
  - Số lần đi trễ
  - Số lần về sớm
  - Tỷ lệ tuân thủ (%)

**API tương ứng:** `GET /api/v1/reports/summary`

---

### **1.7. Cài đặt Công ty (Company Settings)**

**Đường dẫn:** `/settings`  
**Chức năng:** Cấu hình các thông số của công ty

**Thành phần giao diện:**

#### **1.7.1. Thông tin Công ty**
- Logo công ty (upload)
- Tên công ty
- Địa chỉ
- Số điện thoại
- Email liên hệ
- Website

#### **1.7.2. Cài đặt Chấm công**
- Thời gian chấm công hợp lệ (trước/sau ca bao lâu)
- Cho phép chấm công offline
- Độ chính xác nhận diện (threshold)
- Thời gian timeout nhận diện

#### **1.7.3. Cài đặt Thông báo**
- Email thông báo hàng ngày
- Cảnh báo khi thiết bị offline
- Thông báo khi có chấm công bất thường

#### **1.7.4. Quản lý Tài khoản**
- Danh sách admin phụ
- Phân quyền truy cập
- Thay đổi mật khẩu
- Lịch sử đăng nhập

---

## **2. DEVICE APP (IoT) - GIAO DIỆN THIẾT BỊ CHẤM CÔNG**

### **2.1. Màn hình Chính (Main Screen)**

**Chức năng:** Giao diện chính cho việc chấm công

**Thành phần giao diện:**
- **Header:**
  - Logo công ty
  - Tên địa điểm/thiết bị
  - Thời gian hiện tại (lớn, nổi bật)
  - Trạng thái kết nối (Online/Offline)

- **Khu vực chính:**
  - Camera preview (live video)
  - Khung nhận diện khuôn mặt (highlight khi detect)
  - Hướng dẫn: "Đưa mặt vào khung để chấm công"

- **Panel bên:**
  - Nút "Check In" / "Check Out"
  - Thông tin ca làm việc hiện tại
  - Số lượt chấm công hôm nay
  - Danh sách 5 lượt chấm công gần nhất

- **Footer:**
  - Nút "Cài đặt" (admin)
  - Trạng thái kết nối server
  - Phiên bản app

### **2.2. Kết quả Chấm công (Result Screen)**

**Chức năng:** Hiển thị kết quả sau khi chấm công

**Thành phần giao diện:**

#### **2.2.1. Chấm công Thành công**
- **Icon thành công** (checkmark xanh lá)
- **Thông tin nhân viên:**
  - Ảnh đại diện
  - Tên nhân viên
  - Mã nhân viên
- **Chi tiết chấm công:**
  - Loại: "Chấm công VÀO" / "Chấm công RA"
  - Thời gian chấm công
  - Trạng thái: "Đúng giờ" / "Trễ 15 phút" / "Sớm 30 phút"
- **Âm thanh:** Tiếng beep thành công
- **Tự động quay lại màn hình chính sau 3 giây**

#### **2.2.2. Chấm công Thất bại**
- **Icon thất bại** (X đỏ)
- **Thông báo lỗi:**
  - "Không nhận diện được khuôn mặt"
  - "Khuôn mặt chưa được đăng ký"
  - "Thiết bị đang offline"
  - "Đã chấm công rồi (nếu trong thời gian chờ)"
- **Hướng dẫn khắc phục**
- **Nút "Thử lại"**
- **Âm thanh:** Tiếng beep lỗi

### **2.3. Màn hình Cài đặt (Settings Screen)**

**Chức năng:** Cấu hình thiết bị (chỉ admin)

**Bảo mật:** Yêu cầu mật khẩu admin hoặc QR Code

**Thành phần giao diện:**
- **Thông tin thiết bị:**
  - Tên thiết bị
  - Mã thiết bị/Serial
  - Địa chỉ IP
  - Phiên bản firmware

- **Cấu hình mạng:**
  - WiFi SSID
  - Trạng thái kết nối
  - Server URL
  - Test kết nối

- **Cấu hình camera:**
  - Độ phân giải
  - Độ sáng/Tương phản
  - Tần số quét (FPS)
  - Test camera

- **Cài đặt chấm công:**
  - Timeout nhận diện
  - Độ nhạy (threshold)
  - Âm thanh bật/tắt

- **Đồng bộ dữ liệu:**
  - Tải danh sách nhân viên
  - Tải dữ liệu khuôn mặt
  - Đồng bộ thời gian
  - Upload log

---

## **3. SYSTEM ADMIN INTERFACE - GIAO DIỆN QUẢN TRỊ HỆ THỐNG**

### **3.1. Dashboard Quản trị Hệ thống**

**Đường dẫn:** `/system/dashboard`  
**Chức năng:** Tổng quan về toàn bộ hệ thống

**Thành phần giao diện:**
- **Thống kê tổng thể:**
  - Tổng số công ty khách hàng
  - Tổng số nhân viên trong hệ thống
  - Tổng số thiết bị đang hoạt động
  - Tổng lượt chấm công hôm nay

- **Biểu đồ:**
  - Xu hướng đăng ký công ty mới theo tháng
  - Biểu đồ tải hệ thống (CPU, RAM, Database)
  - Top 10 công ty có lượng chấm công cao nhất

- **Cảnh báo hệ thống:**
  - Các công ty sắp hết hạn dịch vụ
  - Thiết bị offline lâu
  - Lỗi hệ thống cần xử lý

### **3.2. Quản lý Công ty Khách hàng**

#### **3.2.1. Danh sách Công ty**

**Đường dẫn:** `/system/companies`  
**Chức năng:** Quản lý tất cả công ty khách hàng

**Thành phần giao diện:**
- **Thanh công cụ:**
  - Nút "Thêm công ty mới"
  - Tìm kiếm theo tên công ty, email
  - Bộ lọc (Trạng thái, Gói dịch vụ, Ngày hết hạn)

- **Bảng danh sách:**
  - Logo + Tên công ty
  - Email liên hệ
  - Số nhân viên / Giới hạn
  - Gói dịch vụ
  - Ngày đăng ký - Hết hạn
  - Trạng thái (Active/Suspended/Expired)
  - Dung lượng sử dụng
  - Thao tác (Xem, Sửa, Khóa/Mở khóa, Xóa)

#### **3.2.2. Chi tiết Công ty**

**Đường dẫn:** `/system/companies/{id}`  
**Chức năng:** Xem chi tiết và thống kê của một công ty

**Thành phần giao diện:**

**Tab Thông tin:**
- Thông tin cơ bản công ty
- Thông tin thanh toán
- Lịch sử thay đổi gói dịch vụ
- Danh sách admin của công ty

**Tab Thống kê:**
- Số lượng nhân viên theo thời gian
- Số lượt chấm công hàng ngày
- Tỷ lệ sử dụng dịch vụ
- Thống kê thiết bị

**Tab Audit Log:**
- Lịch sử thao tác của admin công ty
- Log đăng nhập/đăng xuất  
- Thay đổi cấu hình quan trọng

### **3.3. Giám sát Hệ thống (System Monitoring)**

**Đường dẫn:** `/system/monitoring`  
**Chức năng:** Theo dõi hiệu suất và trạng thái hệ thống

**Thành phần giao diện:**

#### **3.3.1. Server Status**
- **Microservices:**
  - Auth Service (CPU, RAM, Response time)
  - Identity & Org Service
  - Device Management Service  
  - Attendance Service
  - Analytics & Reporting Service
  - Trạng thái: Running/Error/Stopped

- **Database & Cache:**
  - PostgreSQL Cluster (Connections, Query time)
  - ScyllaDB/Cassandra (Read/Write latency)
  - Redis Cluster (Hit ratio, Memory usage)
  - Minio Cluster (Storage usage, Bandwidth)

#### **3.3.2. Performance Metrics**
- **API Response Times:** Biểu đồ realtime theo từng endpoint
- **Error Rates:** Tỷ lệ lỗi 4xx, 5xx theo service
- **Database Performance:** Query execution time, connection pool
- **Queue Status:** Kafka consumer lag, message throughput

#### **3.3.3. Alert Management**
- **Active Alerts:** Danh sách cảnh báo đang active
- **Alert Rules:** Cấu hình ngưỡng cảnh báo
- **Notification Settings:** Email, Slack, SMS notifications

---

## **4. RESPONSIVE DESIGN & UX/UI GUIDELINES**

### **4.1. Thiết kế Đáp ứng (Responsive Design)**

**Breakpoints:**
- **Desktop:** ≥ 1200px - Full layout với sidebar
- **Tablet:** 768px - 1199px - Collapsed sidebar, responsive tables
- **Mobile:** < 768px - Mobile-first navigation, card layouts

**Grid System:**
- Bootstrap 5 grid system
- Flexbox cho component alignment
- CSS Grid cho dashboard layouts

### **4.2. Color Scheme & Typography**

**Primary Colors:**
- Primary: #2563eb (Blue)
- Success: #10b981 (Green) 
- Warning: #f59e0b (Orange)
- Danger: #ef4444 (Red)
- Info: #06b6d4 (Cyan)

**Typography:**
- Font chính: Inter, -apple-system, BlinkMacSystemFont
- Headings: 600-700 font-weight
- Body text: 400 font-weight
- Code/Monospace: 'Fira Code', monospace

### **4.3. Component Standards**

**Buttons:**
- Primary: Solid background với primary color
- Secondary: Outline style
- Ghost: Text style với hover effect
- Sizes: sm (32px), md (40px), lg (48px)

**Tables:**
- Zebra striping cho dễ đọc
- Sticky headers khi scroll
- Loading skeleton khi fetch data
- Empty states với illustration

**Forms:**
- Floating labels
- Inline validation với real-time feedback
- Required field indicators (*)
- Help text cho complex fields

**Navigation:**
- Breadcrumb cho nested pages
- Active state highlighting
- Collapsible sidebar trên mobile
- Tab navigation cho related content

---

## **5. TECHNICAL SPECIFICATIONS**

### **5.1. Frontend Technology Stack**

**Core Framework:**
- **React 18+** hoặc **Vue.js 3+**
- **TypeScript** cho type safety
- **Vite** hoặc **Next.js** cho build tool

**State Management:**
- **Redux Toolkit** (React) hoặc **Pinia** (Vue)
- **React Query** / **TanStack Query** cho server state

**UI Framework:**
- **Tailwind CSS** cho styling
- **Headless UI** / **Radix UI** cho components
- **Chart.js** / **Recharts** cho biểu đồ

**Utilities:**
- **Axios** cho HTTP requests
- **Socket.io-client** cho WebSocket
- **Date-fns** cho date manipulation
- **React-hook-form** / **VeeValidate** cho forms

### **5.2. PWA Features**

**Service Worker:**
- Cache strategy cho static assets
- Offline fallback pages
- Background sync cho attendance data

**Manifest:**
- App icons cho multiple resolutions
- Theme colors
- Start URL và display mode

**Push Notifications:**
- Browser notifications cho admin alerts
- Background sync results
- Device offline/online status

### **5.3. Performance Optimization**

**Code Splitting:**
- Route-based code splitting
- Component lazy loading
- Dynamic imports cho heavy components

**Caching Strategy:**
- HTTP caching headers
- Browser storage (localStorage, sessionStorage)
- IndexedDB cho large datasets

**Image Optimization:**
- WebP format với fallbacks
- Responsive images với srcset
- Lazy loading cho images

---

## **6. SECURITY CONSIDERATIONS**

### **6.1. Authentication & Authorization**

**Token Management:**
- JWT access tokens (15-30 phút expiry)
- Refresh token rotation
- Secure HTTP-only cookies
- XSS protection

**Role-Based Access:**
- System Admin > Company Admin > Employee
- Feature-level permissions
- API endpoint protection
- UI component conditional rendering

### **6.2. Data Protection**

**Sensitive Data:**
- Face embeddings encryption
- PII data masking trong UI
- Secure file uploads
- Input sanitization

**Communication:**
- HTTPS enforced
- WebSocket Secure (WSS)
- API rate limiting
- CORS configuration

---

## **7. ACCESSIBILITY (A11Y) STANDARDS**

### **7.1. WCAG 2.1 Compliance**

**Level AA Requirements:**
- Keyboard navigation support
- Screen reader compatibility
- High contrast mode
- Focus management
- Alt texts cho images
- Semantic HTML structure

**Aria Labels:**
- Form field descriptions
- Button purposes
- Table headers
- Dynamic content updates
- Loading states

### **7.2. Internationalization (i18n)**

**Multi-language Support:**
- Vietnamese (primary)
- English (secondary)
- Dynamic language switching
- RTL layout support (future)
- Date/time localization
- Number format localization

---

Tài liệu này cung cấp mô tả chi tiết về tất cả các trang giao diện client trong hệ thống Face Attendance SaaS, bao gồm chức năng, layout, và technical requirements. Mỗi trang được thiết kế để tối ưu trải nghiệm người dùng và đáp ứng các yêu cầu nghiệp vụ của hệ thống chấm công thông minh.