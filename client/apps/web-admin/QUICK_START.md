# ⚡ Quick Start Guide

Hướng dẫn nhanh để chạy Face Attendance Web Admin.

## 🚀 Bắt đầu trong 3 bước

### 1️⃣ Cài đặt

```bash
# Clone repository
git clone <repository-url>
cd client/apps/web-admin

# Cài đặt dependencies
npm install
```

### 2️⃣ Chạy

```bash
npm run dev
```

### 3️⃣ Mở trình duyệt

```
http://localhost:3003
```

---

## 📋 Yêu cầu

- Node.js >= 16.0.0
- npm >= 7.0.0

Kiểm tra phiên bản:
```bash
node --version
npm --version
```

---

## 🎯 Pages đã hoàn thành

| Page | URL | Status |
|------|-----|--------|
| Dashboard | `/` | ✅ Done |
| Employees | `/employees` | ✅ Done |
| Attendance | `/attendance` | ✅ Done |
| Reports | `/reports` | ✅ Done |
| Shifts | `/shifts` | ✅ Done |
| Settings | `/settings` | ✅ Done |
| Devices | `/devices` | 🚧 In Progress |

---

## 🛠 Commands

```bash
npm run dev      # Development server
npm run build    # Production build
npm run preview  # Preview production
npm run lint     # Lint code
```

---

## 🐛 Gặp lỗi?

### Port đã sử dụng?
Vite sẽ tự động chọn port khác (3004, 3005...)

### Module not found?
```bash
rm -rf node_modules
npm install
```

### TypeScript errors?
Restart VS Code hoặc chạy:
```bash
npx tsc --noEmit
```

---

## 📚 Đọc thêm

- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Hướng dẫn chi tiết
- [PHASE4_IMPLEMENTATION.md](./PHASE4_IMPLEMENTATION.md) - Tiến độ dự án
- [doc/](./doc/) - Documentation đầy đủ

---

**Happy Coding! 🎉**
