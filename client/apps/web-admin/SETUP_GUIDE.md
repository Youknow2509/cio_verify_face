# 🎯 Face Attendance Web Admin

Hệ thống quản lý chấm công bằng khuôn mặt - Giao diện quản trị web.

## 📋 Mục lục

- [Giới thiệu](#giới-thiệu)
- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [Cài đặt](#cài-đặt)
- [Chạy ứng dụng](#chạy-ứng-dụng)
- [Cấu trúc dự án](#cấu-trúc-dự-án)
- [Tính năng](#tính-năng)
- [Công nghệ sử dụng](#công-nghệ-sử-dụng)
- [Hướng dẫn phát triển](#hướng-dẫn-phát-triển)
- [Troubleshooting](#troubleshooting)

---

## 🎨 Giới thiệu

Web Admin là giao diện quản trị cho hệ thống chấm công bằng khuôn mặt, được xây dựng với React 18 + TypeScript + Vite.

### ✨ Tính năng chính

- **Dashboard**: Tổng quan hệ thống với biểu đồ, thống kê thời gian thực
- **Quản lý nhân viên**: CRUD nhân viên, import/export dữ liệu
- **Chấm công**: Theo dõi giờ vào/ra, xuất báo cáo
- **Quản lý thiết bị**: Cấu hình và giám sát thiết bị chấm công
- **Báo cáo**: Tạo báo cáo theo ngày/tuần/tháng, xuất Excel/PDF
- **Ca làm việc**: Quản lý ca và lịch làm việc
- **Cài đặt**: Cấu hình hệ thống, phân quyền

---

## 💻 Yêu cầu hệ thống

### Bắt buộc

- **Node.js**: >= 16.0.0 (khuyến nghị 18.x hoặc 20.x)
- **npm**: >= 7.0.0 hoặc **yarn** >= 1.22.0
- **Git**: Để clone repository

### Kiểm tra phiên bản

```bash
node --version   # v18.x.x hoặc cao hơn
npm --version    # 7.x.x hoặc cao hơn
```

---

## 🚀 Cài đặt

### 1. Clone repository

```bash
git clone <repository-url>
cd client/apps/web-admin
```

### 2. Cài đặt dependencies

**Sử dụng npm:**
```bash
npm install
```

**Hoặc sử dụng yarn:**
```bash
yarn install
```

### 3. Cấu hình môi trường (tùy chọn)

Tạo file `.env` trong thư mục `web-admin`:

```env
VITE_API_URL=http://localhost:3000/api
VITE_WS_URL=ws://localhost:3000
```

---

## 🏃 Chạy ứng dụng

### Development mode (Chế độ phát triển)

```bash
npm run dev
```

Ứng dụng sẽ chạy tại: **http://localhost:3003**

> **Lưu ý**: Nếu port 3003 đã được sử dụng, Vite sẽ tự động chọn port khác (3004, 3005...)

### Production build (Build cho production)

```bash
npm run build
```

Output sẽ được tạo trong thư mục `dist/`

### Preview production build

```bash
npm run preview
```

### Lint code

```bash
npm run lint
```

---

## 📁 Cấu trúc dự án

```
web-admin/
├── public/              # Static assets
│   ├── favicon.ico
│   └── index.html
├── src/
│   ├── components/      # Reusable components
│   │   ├── Badge/       # Status badges
│   │   ├── Card/        # Card component
│   │   ├── Header/      # App header
│   │   ├── Sidebar/     # Navigation sidebar
│   │   ├── Table/       # Data table
│   │   └── Toolbar/     # Search & filter toolbar
│   ├── contexts/        # React contexts
│   ├── hooks/           # Custom hooks
│   │   └── useVirtualizedTable.ts
│   ├── layouts/         # Page layouts
│   │   └── Layout/
│   ├── pages/           # Application pages
│   │   ├── Dashboard/   # ✅ Main dashboard
│   │   ├── Employees/   # ✅ Employee management
│   │   ├── Attendance/  # ✅ Attendance tracking
│   │   ├── Devices/     # 🚧 Device management
│   │   ├── Reports/     # ✅ Reports & analytics
│   │   ├── Shifts/      # ✅ Shift management
│   │   └── Settings/    # ✅ System settings
│   ├── services/        # API services
│   │   └── mock/        # Mock data for development
│   ├── styles/          # Global styles
│   │   ├── _tokens.scss # Design tokens
│   │   ├── _mixins.scss # SCSS mixins
│   │   └── _globals.scss
│   ├── types/           # TypeScript types
│   ├── utils/           # Utility functions
│   │   └── csv.ts       # CSV export utility
│   ├── App.tsx          # Main App component
│   └── main.tsx         # Entry point
├── doc/                 # Documentation
├── package.json
├── tsconfig.json        # TypeScript config
├── vite.config.ts       # Vite config
└── README.md
```

---

## 🎨 Tính năng

### ✅ Dashboard
- 4 stat cards với gradient đẹp mắt
- Biểu đồ chấm công 7 ngày
- Danh sách hoạt động gần đây
- Quick actions panel
- Upcoming events calendar

### ✅ Quản lý nhân viên
- Tìm kiếm, lọc theo phòng ban, trạng thái
- Thêm/Sửa/Xóa nhân viên
- Form validation đầy đủ
- Upload ảnh khuôn mặt
- Export danh sách

### ✅ Chấm công
- Filter theo ngày, trạng thái
- Hiển thị giờ vào/ra, giờ công
- Badge trạng thái (Đúng giờ, Đi trễ, Về sớm)
- Export Excel

### ✅ Báo cáo
- Báo cáo theo ngày/tuần/tháng
- Tùy chỉnh khoảng thời gian
- Export Excel & PDF
- Thống kê theo phòng ban

### ✅ Quản lý ca
- Tạo/Sửa/Xóa ca làm việc
- Cấu hình giờ làm, giờ nghỉ
- Active/Inactive status
- Card-based layout

### ✅ Cài đặt
- 4 tabs: General, Attendance, Notification, Security
- Cấu hình giờ làm việc
- Thiết lập thông báo
- Bảo mật hệ thống

---

## 🛠 Công nghệ sử dụng

### Core
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool & dev server

### Styling
- **SCSS** - CSS preprocessor
- **CSS Modules** - Scoped styling
- **Material Design 3** - Design system

### Charts & Visualization
- **Recharts** - Chart library

### State Management
- **React Context API** - Global state
- **React Hooks** - Local state

### Code Quality
- **ESLint** - Code linting
- **TypeScript** - Static typing

---

## 👨‍💻 Hướng dẫn phát triển

### 1. Tạo component mới

```bash
# Tạo folder cho component
mkdir src/components/MyComponent

# Tạo files
touch src/components/MyComponent/MyComponent.tsx
touch src/components/MyComponent/MyComponent.module.scss
```

**MyComponent.tsx:**
```tsx
import styles from './MyComponent.module.scss';

interface MyComponentProps {
  title: string;
}

export function MyComponent({ title }: MyComponentProps) {
  return (
    <div className={styles.container}>
      <h2>{title}</h2>
    </div>
  );
}
```

**MyComponent.module.scss:**
```scss
@import '../../styles/tokens';
@import '../../styles/mixins';

.container {
  padding: 1rem;
  background: var(--bg-primary);
  border-radius: 8px;
}
```

### 2. Tạo page mới

```bash
# Tạo folder
mkdir src/pages/MyPage

# Tạo files
touch src/pages/MyPage/MyPage.tsx
touch src/pages/MyPage/MyPage.module.scss
```

### 3. Sử dụng Toolbar component

```tsx
import { Toolbar, ToolbarSection, SearchBox } from '../../components/Toolbar/Toolbar';

function MyPage() {
  const [search, setSearch] = useState('');
  
  return (
    <Toolbar>
      <ToolbarSection>
        <SearchBox value={search} onChange={setSearch} />
      </ToolbarSection>
      <ToolbarSection align="right">
        <button>Add New</button>
      </ToolbarSection>
    </Toolbar>
  );
}
```

### 4. Sử dụng Badge component

```tsx
import { Badge } from '../../components/Badge/Badge';

<Badge variant="success">Active</Badge>
<Badge variant="warning">Pending</Badge>
<Badge variant="error">Offline</Badge>
```

### 5. Export CSV

```tsx
import { exportToCSV } from '../../utils/csv';

const data = [
  { name: 'John', age: 30 },
  { name: 'Jane', age: 25 }
];

exportToCSV(data, 'employees.csv');
```

---

## 🐛 Troubleshooting

### Port đã được sử dụng

**Lỗi:**
```
Port 3003 is in use, trying another one...
```

**Giải pháp:**
- Vite sẽ tự động chọn port khác (3004, 3005...)
- Hoặc đóng ứng dụng đang dùng port 3003
- Hoặc cấu hình port khác trong `vite.config.ts`

### Module not found

**Lỗi:**
```
Cannot find module '@/components/...'
```

**Giải pháp:**
```bash
# Xóa node_modules và cài lại
rm -rf node_modules
npm install
```

### TypeScript errors

**Giải pháp:**
```bash
# Check TypeScript errors
npx tsc --noEmit

# Restart VS Code TypeScript server
# Ctrl+Shift+P > TypeScript: Restart TS Server
```

### SCSS compilation errors

**Lỗi:**
```
Undefined mixin 'respond'
```

**Giải pháp:**
- Kiểm tra import trong file SCSS:
```scss
@import '../../styles/tokens';
@import '../../styles/mixins';
```

### Hot reload không hoạt động

**Giải pháp:**
```bash
# Restart dev server
# Ctrl+C để stop
npm run dev
```

---

## 📚 Tài liệu tham khảo

### Documentation
- Xem folder `doc/` để biết thêm chi tiết về:
  - Architecture
  - API endpoints
  - Component structure
  - UI workflows

### Quan trọng
- `doc/frontend_development_guide.md` - Hướng dẫn phát triển
- `doc/component_architecture_map.md` - Cấu trúc components
- `PHASE4_IMPLEMENTATION.md` - Tiến độ implementation

---

## 🔗 Links

- **Development**: http://localhost:3003
- **API Documentation**: (Sẽ cập nhật)
- **Design System**: Material Design 3

---

## 📝 Scripts

| Command | Mô tả |
|---------|-------|
| `npm run dev` | Chạy dev server |
| `npm run build` | Build production |
| `npm run preview` | Preview production build |
| `npm run lint` | Lint code |

---

## ✅ Checklist cho máy mới

- [ ] Cài đặt Node.js >= 16.0.0
- [ ] Cài đặt Git
- [ ] Clone repository
- [ ] Chạy `npm install`
- [ ] Chạy `npm run dev`
- [ ] Mở http://localhost:3003
- [ ] Kiểm tra tất cả pages hoạt động
- [ ] Đọc documentation trong folder `doc/`

---

## 🎉 Hoàn thành!

Bây giờ bạn đã sẵn sàng để phát triển. Happy coding! 🚀

---

## 📞 Liên hệ & Hỗ trợ

Nếu gặp vấn đề, vui lòng:
1. Kiểm tra phần [Troubleshooting](#troubleshooting)
2. Xem documentation trong folder `doc/`
3. Liên hệ team leader

---

**Version**: 1.0.0  
**Last Updated**: October 5, 2025  
**Status**: ✅ Production Ready (85% - Device pages đang phát triển)
