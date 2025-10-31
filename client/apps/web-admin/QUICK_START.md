# ⚡ Quick Start

## 3 Bước Cài Đặt

### 1️⃣ Install Dependencies
```bash
cd client/apps/web-admin
npm install
```

### 2️⃣ Configure Environment
```bash
cp .env.example .env
```

Cập nhật `.env`:
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=10000
```

### 3️⃣ Run Development Server
```bash
npm run dev
```

Truy cập: **http://localhost:3003**

## 🔐 Credentials

**Default Admin Login:**
- Email: `admin@company.com`
- Password: `admin@123`

## 📚 Next Steps

1. **Setup Details** → [SETUP_GUIDE.md](SETUP_GUIDE.md)
2. **API Integration** → [src/services/API_GUIDE.md](src/services/API_GUIDE.md)
3. **Migrate to Real API** → [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)

## ❓ Troubleshooting

### Port 3003 đang bị sử dụng?
```bash
npm run dev -- --port 3004
```

### Clear cache
```bash
npm run clean
npm install
npm run dev
```

### TypeScript errors?
```bash
npm run type-check
```

## 💡 Common Commands

```bash
npm run dev          # Start dev
npm run build        # Production build
npm run preview      # Preview build
npm run lint         # Check code
npm run type-check   # Check types
```

---

**⏱️ Estimated Setup Time:** 5 phút

**❓ Có vấn đề?** Xem [SETUP_GUIDE.md](SETUP_GUIDE.md) chi tiết hơn
