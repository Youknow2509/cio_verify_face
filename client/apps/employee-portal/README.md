# Employee Portal

Cổng thông tin nhân viên cho hệ thống chấm công nhận dạng khuôn mặt CIO Verify Face.

## Tính năng

### 🔐 Xác thực
- Đăng nhập bằng email/mã nhân viên và mật khẩu
- Tự động làm mới token
- Bảo vệ các route yêu cầu xác thực

### 📊 Bảng điều khiển
- Thống kê chấm công hôm nay
- Tổng số ngày làm việc trong tháng
- Số ngày có mặt
- Danh sách ca làm việc hiện tại

### ⏰ Chấm công
- Xem lịch sử chấm công theo tháng
- Hiển thị thời gian vào/ra
- Phương thức xác thực (khuôn mặt)
- Điểm số xác thực
- Trạng thái đồng bộ

### 📅 Tổng hợp theo ngày
- Xem tổng hợp chấm công hàng ngày
- Giờ vào/ra thực tế
- Thời gian đi muộn/về sớm
- Tổng giờ làm việc

### 🕐 Ca làm việc
- Xem danh sách ca làm việc được gán
- Thông tin giờ làm việc
- Thời gian hiệu lực
- Trạng thái ca làm việc

### 👤 Hồ sơ cá nhân
- Xem thông tin cá nhân
- Tạo yêu cầu cập nhật khuôn mặt
- Upload ảnh khuôn mặt mới (khi được duyệt)
- Theo dõi trạng thái yêu cầu

### 📄 Xuất báo cáo
- Xuất báo cáo chấm công theo tháng
- Hỗ trợ định dạng: Excel, PDF, CSV
- Gửi báo cáo qua email

## Cài đặt

### Yêu cầu
- Node.js >= 18
- pnpm (hoặc npm)

### Cài đặt dependencies
```bash
cd client
pnpm install
```

### Chạy development server
```bash
cd client/apps/employee-portal
npm run dev
```

Ứng dụng sẽ chạy tại: http://localhost:3003

### Build production
```bash
cd client/apps/employee-portal
npm run build
```

## Cấu hình

### Biến môi trường
Tạo file `.env` trong thư mục `client/apps/employee-portal`:

```env
VITE_API_URL=http://localhost:8080
```

## Cấu trúc thư mục

```
src/
├── components/          # Shared components
│   ├── layouts/        # Layout components (AppBar, Sidebar, MainLayout)
│   └── ProtectedRoute.tsx
├── features/           # Feature modules
│   ├── auth/          # Authentication
│   ├── dashboard/     # Dashboard
│   ├── attendance/    # Attendance & Export
│   ├── shifts/        # Shifts
│   └── profile/       # Profile & Face Update
├── routes/            # Route configuration
├── services/          # API services
├── stores/            # State management (Zustand)
├── theme/             # MUI theme configuration
├── App.tsx
└── main.tsx
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/refresh` - Làm mới token
- `GET /api/v1/auth/me` - Lấy thông tin người dùng

### Profile Update
- `POST /api/v1/profile-update/requests` - Tạo yêu cầu cập nhật khuôn mặt
- `GET /api/v1/profile-update/requests/me` - Lấy trạng thái yêu cầu
- `POST /api/v1/profile-update/face` - Upload khuôn mặt mới

### Shifts
- `GET /api/v1/shift/employee` - Lấy danh sách ca làm việc

### Attendance
- `GET /api/v1/employee/my-attendance-records` - Lịch sử chấm công
- `GET /api/v1/employee/my-daily-summaries` - Tổng hợp theo ngày
- `POST /api/v1/employee/export-monthly-summary` - Xuất báo cáo

## Công nghệ sử dụng

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Material-UI v5** - UI components
- **React Router v6** - Routing
- **Zustand** - State management
- **Axios** - HTTP client
- **Date-fns** - Date formatting
- **Vite** - Build tool

## Tính năng nổi bật

### 🎨 Giao diện
- Dark theme với gradient màu đẹp mắt
- Responsive design cho mobile và desktop
- Animations mượt mà
- Icons từ Material Icons

### 🔒 Bảo mật
- JWT token authentication
- Tự động làm mới token
- Protected routes
- Token lưu trong localStorage

### 📱 Responsive
- Sidebar ẩn/hiện trên mobile
- Tables responsive
- Cards layout linh hoạt

### ⚡ Performance
- Code splitting
- Lazy loading
- Optimized bundle size

## Hướng dẫn sử dụng

### Đăng nhập
1. Mở trang http://localhost:3003/login
2. Nhập username (email hoặc mã nhân viên) và mật khẩu
3. Click "Đăng nhập"

### Xem chấm công
1. Vào menu "Chấm công"
2. Chọn tháng cần xem
3. Xem danh sách các lần chấm công

### Xuất báo cáo
1. Vào menu "Xuất báo cáo"
2. Chọn tháng và định dạng file
3. Click "Xuất báo cáo"
4. Báo cáo sẽ được gửi qua email

### Cập nhật khuôn mặt
1. Vào menu "Hồ sơ cá nhân"
2. Click "Tạo yêu cầu mới"
3. Nhập lý do và gửi yêu cầu
4. Sau khi được duyệt, sử dụng token để upload ảnh mới

## Troubleshooting

### Lỗi kết nối API
- Kiểm tra biến môi trường `VITE_API_URL`
- Đảm bảo backend đang chạy
- Kiểm tra CORS settings trên backend

### Lỗi build
- Xóa `node_modules` và `pnpm-lock.yaml`
- Chạy lại `pnpm install`
- Xóa thư mục `dist` nếu có

## License

Copyright © 2025 CIO Verify Face
