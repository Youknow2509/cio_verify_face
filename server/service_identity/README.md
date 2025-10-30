# Identity & Organization Microservice

Microservice chuyên biệt xử lý định danh và quản lý tổ chức trong hệ thống microservices.

## Chức năng chính

- **Quản lý công ty**: Tạo, sửa, xóa, xem danh sách công ty với các gói dịch vụ
- **Quản lý người dùng**: Tạo, sửa, xóa, xem danh sách nhân viên và quản trị viên
- **Quản lý dữ liệu khuôn mặt**: Upload và quản lý ảnh khuôn mặt của nhân viên
- **Xác thực và phân quyền**: JWT authentication với các role khác nhau
- **Health Check**: Monitoring và health check endpoints
- **Service Discovery**: Tích hợp với service registry
- **Caching**: Redis cache để tối ưu hiệu suất

## Công nghệ sử dụng

- **NestJS**: Framework Node.js
- **TypeORM**: ORM cho PostgreSQL
- **PostgreSQL**: Cơ sở dữ liệu chính
- **Redis**: Cache và session storage
- **JWT**: Xác thực người dùng
- **Multer**: Upload file
- **Swagger**: API documentation

## Cài đặt

1. Clone repository:
```bash
git clone <repository-url>
cd identity-organization-service
```

2. Cài đặt dependencies:
```bash
npm install
```

3. Cấu hình environment:
```bash
cp env.example .env
```

4. Cấu hình database và Redis trong file `.env`

5. Chạy migration (nếu cần):
```bash
npm run migration:run
```

6. Chạy ứng dụng:
```bash
# Development
npm run start:dev

# Production
npm run build
npm run start:prod
```

## API Endpoints

### Companies
- `GET /api/v1/companies` - Danh sách công ty
- `POST /api/v1/companies` - Tạo mới công ty
- `GET /api/v1/companies/{id}` - Xem thông tin công ty
- `PUT /api/v1/companies/{id}` - Sửa thông tin công ty
- `DELETE /api/v1/companies/{id}` - Xóa công ty

### Users
- `GET /api/v1/users` - Danh sách nhân viên
- `POST /api/v1/users` - Thêm mới nhân viên
- `GET /api/v1/users/{id}` - Xem thông tin nhân viên
- `PUT /api/v1/users/{id}` - Sửa thông tin nhân viên
- `DELETE /api/v1/users/{id}` - Xóa nhân viên

### Face Data
- `POST /api/v1/users/{userId}/face-data` - Upload ảnh khuôn mặt
- `GET /api/v1/users/{userId}/face-data` - Xem danh sách ảnh khuôn mặt
- `DELETE /api/v1/users/{userId}/face-data/{faceId}` - Xóa ảnh khuôn mặt

### Authentication
- `POST /auth/login` - Đăng nhập
- `POST /auth/profile` - Xem thông tin profile

### Health Check
- `GET /health` - Tổng quan health check
- `GET /health/database` - Kiểm tra database
- `GET /health/redis` - Kiểm tra Redis
- `GET /health/memory` - Kiểm tra memory usage
- `GET /health/disk` - Kiểm tra disk usage

## Swagger Documentation

Sau khi chạy ứng dụng, truy cập: `http://localhost:3000/api/docs`

## Cấu trúc Database

Microservice sử dụng PostgreSQL với các bảng chính:
- `companies` - Thông tin công ty (bao gồm plan, settings)
- `companies_secret` - Bí mật công ty
- `users` - Thông tin người dùng
- `employees` - Thông tin nhân viên (bao gồm permissions, manager_id)
- `face_data` - Dữ liệu khuôn mặt

**Lưu ý**: Các bảng liên quan đến attendance, devices, work shifts sẽ được xử lý bởi các microservice khác.

## Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_DATABASE=identity_service

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Application
PORT=3000
NODE_ENV=development

# File Upload
MAX_FILE_SIZE=5242880
UPLOAD_PATH=./uploads

# Face Processing
FACE_PROCESSING_ENABLED=true
MAX_FACE_IMAGES_PER_USER=5
```

## Testing

```bash
# Unit tests
npm run test

# E2E tests
npm run test:e2e

# Test coverage
npm run test:cov
```

## License

MIT
