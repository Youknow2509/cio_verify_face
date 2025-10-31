# 🎯 Face Attendance Web Admin

Hệ thống quản lý chấm công bằng khuôn mặt - Giao diện quản trị công ty.

## 🚀 Quick Start

```bash
# 1. Cài đặt
npm install

# 2. Chạy development server
npm run dev

# 3. Mở trình duyệt
# http://localhost:3003
```

## 📚 Documentation

👉 **[Xem Tất Cả Tài Liệu →](DOCS_INDEX.md)**

### ⚡ Bắt Đầu Nhanh

- **[QUICK_START.md](QUICK_START.md)** - 3 bước cài đặt
- **[SETUP_GUIDE.md](SETUP_GUIDE.md)** - Cấu hình chi tiết

### 👨‍💻 Phát Triển

- **[src/services/API_GUIDE.md](src/services/API_GUIDE.md)** - API Services
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Code Standards
- **[MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)** - Mock → Real API

### 📝 Tham Khảo

- **[CHANGELOG.md](CHANGELOG.md)** - Lịch sử & Thay đổi
- **[DOCS_INDEX.md](DOCS_INDEX.md)** - Danh sách đầy đủ

## ✨ Features

### Dashboard

- Thống kê nhân sự (tổng số, hôm nay check-in, muộn giờ, thiết bị online)
- Biểu đồ chấm công theo ngày
- Hoạt động gần đây

### Quản lý Nhân viên

- Danh sách nhân viên với phân trang, tìm kiếm
- Thêm/sửa/xóa nhân viên
- Quản lý ảnh khuôn mặt

### Chấm công

- Check-in/Check-out bằng camera
- Xem lịch sử chấm công
- Xuất báo cáo

### Báo cáo

- Báo cáo hàng ngày
- Thống kê tổng hợp
- Xuất Excel/PDF

### Ca & Lịch

- Quản lý ca làm việc
- Lịch làm việc cho nhân viên

### Cài đặt

- Quản lý công ty (tên, múi giờ, format ngày)
- Cấu hình toàn hệ thống

## 🔌 API Services

Tất cả API endpoints sử dụng prefix `/api/v1/` và được tổ chức theo services:

- **Auth** - Đăng nhập, token, kích hoạt thiết bị
- **Users** - Quản lý nhân viên, ảnh khuôn mặt
- **Devices** - Quản lý thiết bị
- **Attendance** - Check-in, check-out, lịch sử
- **Shifts & Schedules** - Ca làm việc và lịch
- **Reports** - Báo cáo
- **Signatures** - Chữ ký

Xem chi tiết: [src/services/API_GUIDE.md](src/services/API_GUIDE.md)

## 🏗️ Cấu trúc Folder

```
src/
├── pages/          # React pages
├── components/     # Reusable components
├── services/       # API & business logic
│   ├── api/        # API endpoint functions
│   ├── http.ts     # HTTP client
│   └── error-handler.ts
├── hooks/          # Custom hooks
├── types/          # TypeScript types
├── styles/         # CSS/SCSS
└── utils/          # Helper functions
```

## 🛠️ Development

### Scripts

```bash
npm run dev       # Start dev server
npm run build     # Build for production
npm run preview   # Preview production build
npm run lint      # Check code quality
npm run type-check # Check TypeScript types
```

### Environment

Sao chép `.env.example` thành `.env` và cấu hình:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=10000
```

## 🔐 Roles & Permissions

- **Company Admin** - Quản lý toàn bộ công ty (mặc định)
- **Manager** - Quản lý nhân viên, báo cáo
- **Staff** - Xem chấm công cá nhân

## 📝 Notes

- Ứng dụng sử dụng **React 18** + **TypeScript**
- State management với **Context API** hoặc **Zustand** (tùy chọn)
- UI Components từ **React Bootstrap**
- Chart từ **Recharts**

## 📞 Support

Liên hệ: [support@example.com](mailto:support@example.com)

## 📄 License

MIT License - xem [LICENSE](../../LICENSE)
