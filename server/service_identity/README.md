# Identity & Organization Service

Dịch vụ quản lý công ty, người dùng và nhân viên với API Express.js + PostgreSQL

## Tính năng

- ✅ Quản lý công ty (CRUD)
- ✅ Quản lý người dùng/nhân viên (CRUD)
- ✅ Quản lý dữ liệu khuôn mặt (Create, Read, Delete)
- ✅ Kết nối PostgreSQL với connection pool
- ✅ Error handling và validation
- ✅ TypeScript support
- ✅ **Swagger/OpenAPI documentation**

## API Endpoints

### Companies
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/api/v1/companies` | Danh sách công ty |
| POST | `/api/v1/companies` | Tạo mới công ty |
| GET | `/api/v1/companies/{company_id}` | Xem thông tin công ty |
| PUT | `/api/v1/companies/{company_id}` | Sửa thông tin công ty |
| DELETE | `/api/v1/companies/{company_id}` | Xóa công ty |

### Users
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/api/v1/users` | Danh sách user / nhân viên |
| POST | `/api/v1/users` | Thêm mới nhân viên |
| GET | `/api/v1/users/{user_id}` | Xem thông tin nhân viên |
| PUT | `/api/v1/users/{user_id}` | Sửa thông tin nhân viên |
| DELETE | `/api/v1/users/{user_id}` | Vô hiệu hóa/xóa nhân viên |

### Face Data
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/users/{user_id}/face-data` | Đăng ký ảnh khuôn mặt |
| GET | `/api/v1/users/{user_id}/face-data` | Lấy danh sách ảnh khuôn mặt |
| DELETE | `/api/v1/users/{user_id}/face-data/{fid}` | Xoá ảnh khuôn mặt |

## Cấu trúc Project

```
src/
├── config/           # Database configuration
├── controllers/      # API controllers
├── middleware/       # Express middleware
├── routes/          # API routes
├── services/        # Business logic
├── types/           # TypeScript types
├── utils/           # Utility functions
└── index.ts         # Entry point

sql/                 # Database migration files
```

## Setup & Running

### 1. Cài đặt dependencies

```bash
npm install
```

### 2. Cấu hình environment

Tạo file `.env` từ `.env.example`:

```bash
cp .env.example .env
```

Chỉnh sửa các biến môi trường:

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=identity_service
DB_USER=postgres
DB_PASSWORD=postgres
PORT=3001
NODE_ENV=development
```

### 3. Chạy migrations (nếu sử dụng goose)

```bash
goose -dir sql postgres "user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up
```

### 4. Khởi chạy server

#### Development mode (với auto-reload)

```bash
npm run dev
```

#### Production mode

```bash
npm run build
npm start
```

## Yêu cầu

- Node.js 16+
- PostgreSQL 12+
- npm hoặc yarn

## Dependencies

- **express**: Web framework
- **pg**: PostgreSQL client
- **uuid**: Generate UUID
- **dotenv**: Environment configuration
- **cors**: CORS middleware
- **helmet**: Security headers
- **swagger-ui-express**: Interactive API docs
- **swagger-jsdoc**: Swagger documentation
- **typescript**: Type safety
- **ts-node**: Run TypeScript directly

## Cách sử dụng API

### Test API qua Swagger UI

Khi server chạy, mở: **http://localhost:3001/api-docs**

Swagger UI cung cấp:
- ✅ Interactive API documentation
- ✅ Try-it-out functionality
- ✅ Request/response examples
- ✅ Schema definitions

Xem chi tiết: [SWAGGER.md](./SWAGGER.md)

### Tạo công ty

```bash
curl -X POST http://localhost:3001/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Tech",
    "email": "contact@acmetech.com",
    "phone": "+84-28-0000-0001",
    "status": 1
  }'
```

### Tạo nhân viên

```bash
curl -X POST http://localhost:3001/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@example.com",
    "phone": "0900000001",
    "password": "password123",
    "full_name": "John Doe",
    "role": 2,
    "company_id": "company-uuid"
  }'
```

### Tải ảnh khuôn mặt

```bash
curl -X POST http://localhost:3001/api/v1/users/user-uuid/face-data \
  -H "Content-Type: application/json" \
  -d '{
    "image_url": "https://example.com/image.jpg",
    "quality_score": 0.95
  }'
```

## Lưu ý

- SQL schema không được thay đổi (sử dụng migration files có sẵn)
- Password được hash với salt trước khi lưu
- Tất cả thời gian lưu ở múi giờ UTC (WITH TIME ZONE)

## License

MIT
