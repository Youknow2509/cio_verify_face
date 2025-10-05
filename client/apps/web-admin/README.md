# 🎯 Face Attendance Web Admin

Hệ thống quản lý chấm công bằng khuôn mặt - Giao diện quản trị web.

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](CHANGELOG.md)
[![Status](https://img.shields.io/badge/status-production--ready-green.svg)](PHASE4_IMPLEMENTATION.md)
[![Progress](https://img.shields.io/badge/progress-85%25-yellow.svg)](PHASE4_IMPLEMENTATION.md)

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| **[QUICK_START.md](QUICK_START.md)** | ⚡ Bắt đầu trong 3 bước |
| **[SETUP_GUIDE.md](SETUP_GUIDE.md)** | 📖 Hướng dẫn cài đặt chi tiết |
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | 🤝 Hướng dẫn đóng góp |
| **[CHANGELOG.md](CHANGELOG.md)** | 📝 Lịch sử thay đổi |
| **[PHASE4_IMPLEMENTATION.md](PHASE4_IMPLEMENTATION.md)** | 📊 Tiến độ dự án |

---

## 🚀 Quick Start

```bash
# 1. Cài đặt
npm install

# 2. Chạy
npm run dev

# 3. Mở trình duyệt
# http://localhost:3003
```

👉 **Xem chi tiết**: [QUICK_START.md](QUICK_START.md)

---

## ✨ Features

### ✅ Đã hoàn thành

- ✅ **Dashboard**: Tổng quan với stat cards, charts, activities
- ✅ **Employee Management**: Quản lý nhân viên, form validation
- ✅ **Attendance Tracking**: Theo dõi chấm công, filters, export
- ✅ **Reports**: Báo cáo theo ngày/tuần/tháng, export Excel/PDF
- ✅ **Shift Management**: Quản lý ca làm việc
- ✅ **Settings**: Cài đặt hệ thống với 4 tabs

### 🚧 Đang phát triển

- 🚧 **Device Management**: Quản lý thiết bị (85% complete)

---

## 🛠 Tech Stack

### Core
- **Frontend**: React 18, TypeScript, Vite
- **Routing**: React Router v6
- **Styling**: SCSS Modules, Material Design 3
- **Charts**: Recharts
- **State**: React Hooks + Context API

### Development
- **Build Tool**: Vite
- **Linting**: ESLint
- **Type Checking**: TypeScript (strict mode)

---

## 📁 Project Structure

```
web-admin/
├── public/              # Static assets
├── src/
│   ├── components/      # ✅ Reusable components
│   │   ├── Badge/       # Status badges
│   │   ├── Card/        # Card container
│   │   ├── Header/      # App header
│   │   ├── Sidebar/     # Navigation
│   │   ├── Table/       # Data table
│   │   └── Toolbar/     # Search & filters
│   ├── pages/           # ✅ Application pages
│   │   ├── Dashboard/   # Main dashboard
│   │   ├── Employees/   # Employee management
│   │   ├── Attendance/  # Attendance tracking
│   │   ├── Reports/     # Reports & analytics
│   │   ├── Shifts/      # Shift management
│   │   ├── Settings/    # System settings
│   │   └── Devices/     # 🚧 Device management
│   ├── styles/          # Global styles
│   ├── utils/           # Utility functions
│   ├── hooks/           # Custom hooks
│   └── types/           # TypeScript types
├── doc/                 # Documentation
├── QUICK_START.md       # Quick start guide
├── SETUP_GUIDE.md       # Setup guide
├── CONTRIBUTING.md      # Contributing guide
└── CHANGELOG.md         # Version history
```

---

## 📊 Current Status

**Version**: 1.0.0  
**Progress**: 85%  
**Status**: Production Ready (Device pages in progress)

### Completed (100%)
- ✅ Dashboard with Material Design 3
- ✅ Employee Management (CRUD + Validation)
- ✅ Attendance Tracking (Filters + Export)
- ✅ Reports (Multiple formats)
- ✅ Shift Management
- ✅ Settings (4 tabs)
- ✅ Toolbar Component
- ✅ Utils & Hooks

### In Progress (85%)
- 🚧 Device Management pages

---

## 🎨 Screenshots

| Dashboard | Employees |
|-----------|-----------|
| Gradient stats, charts, activities | Professional filters, validation |

| Attendance | Reports |
|------------|---------|
| Date filters, export Excel | Multiple report types |

---

## 💡 Key Features

### Material Design 3
- Gradient stat cards
- Smooth animations
- Hover effects
- Responsive layouts

### Developer Experience
- TypeScript strict mode
- SCSS Modules
- Hot Module Replacement
- ESLint configured

### Performance
- Vite for fast builds
- Code splitting
- CSS optimization
- Tree shaking

---

## 📝 Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start dev server (port 3003) |
| `npm run build` | Build for production |
| `npm run preview` | Preview production build |
| `npm run lint` | Lint code |

---

## 🔗 Resources

- **Documentation**: See `doc/` folder
- **Setup Guide**: [SETUP_GUIDE.md](SETUP_GUIDE.md)
- **API Docs**: Coming soon
- **Design System**: Material Design 3

---

## 🤝 Contributing

Đọc [CONTRIBUTING.md](CONTRIBUTING.md) để biết cách đóng góp vào dự án.

---

## 📄 License

(License information here)

---

## 👥 Team

- Development Team
- UI/UX Design Team
- QA Team

---

**Built with ❤️ using React + TypeScript + Vite**
│   └── index.ts                # TypeScript type definitions
├── utils/
│   ├── csv.ts                  # CSV export utilities
│   └── format.ts               # Formatting helpers
└── main.tsx                    # App entry point
```

## Available Scripts

```bash
npm run dev          # Start development server
npm run build        # Build for production  
npm run preview      # Preview production build
npm run type-check   # TypeScript type checking
```

## License

MIT