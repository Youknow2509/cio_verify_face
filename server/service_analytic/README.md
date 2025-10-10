# Service Analytics - Báo cáo Chấm công

## Mô tả
Service analytics cho hệ thống chấm công, cung cấp các API để tạo báo cáo chi tiết và tổng hợp về tình hình chấm công của nhân viên.

## Các tính năng chính

### 1. Báo cáo chi tiết ngày (`GET /api/v1/reports/daily`)
- Báo cáo chi tiết về chấm công trong một ngày cụ thể
- Thống kê theo công ty/địa điểm
- Phân tích theo phòng ban và ca làm việc

### 2. Báo cáo tổng hợp tháng (`GET /api/v1/reports/summary`)
- Báo cáo tổng hợp về chấm công trong một tháng
- Thống kê theo tuần và top nhân viên

### 3. Xuất báo cáo (`POST /api/v1/reports/export`)
- Xuất báo cáo ra file Excel/PDF/CSV
- Hỗ trợ xuất theo khoảng thời gian

## Cài đặt và chạy

### 1. Cài đặt dependencies
```bash
npm install
```

### 2. Cấu hình environment
Tạo file `.env` với nội dung:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_NAME=attendance_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRES_IN=1h

# Application Configuration
PORT=3000
NODE_ENV=development
```

### 3. Chạy ứng dụng
```bash
# Development
npm run start:dev

# Production
npm run build
npm run start:prod
```

## API Documentation
Truy cập Swagger UI tại: http://localhost:3000/api/docs

## Health Check
Kiểm tra trạng thái service: http://localhost:3000/health

## Công nghệ sử dụng
- **Framework**: NestJS
- **Database**: PostgreSQL với TypeORM
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Export**: ExcelJS, PDFKit