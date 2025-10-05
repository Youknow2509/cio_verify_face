# 🤝 Contributing Guide

Hướng dẫn đóng góp vào dự án Face Attendance Web Admin.

## 📋 Quy trình làm việc

### 1. Setup môi trường

```bash
# Clone repository
git clone <repository-url>
cd client/apps/web-admin

# Cài đặt
npm install

# Tạo branch mới
git checkout -b feature/ten-tinh-nang
```

### 2. Development

```bash
# Chạy dev server
npm run dev

# Mở http://localhost:3003
```

### 3. Code & Test

- Viết code
- Test thủ công trên browser
- Lint code: `npm run lint`
- Build test: `npm run build`

### 4. Commit

```bash
# Add files
git add .

# Commit với message rõ ràng
git commit -m "feat: thêm chức năng xyz"

# Push
git push origin feature/ten-tinh-nang
```

### 5. Create Pull Request

- Tạo PR trên GitHub/GitLab
- Mô tả rõ ràng thay đổi
- Request review từ team

---

## 📝 Commit Message Convention

Sử dụng format:

```
<type>(<scope>): <subject>

<body>
```

### Types

- `feat`: Tính năng mới
- `fix`: Sửa lỗi
- `refactor`: Refactor code
- `style`: Thay đổi style/UI
- `docs`: Cập nhật documentation
- `test`: Thêm/sửa tests
- `chore`: Các thay đổi khác

### Ví dụ

```
feat(employees): thêm form upload ảnh nhân viên

- Thêm component UploadImageModal
- Thêm validation cho file ảnh
- Integrate với API upload
```

```
fix(dashboard): sửa lỗi chart không hiển thị

- Fix data format cho Recharts
- Thêm fallback khi không có dữ liệu
```

---

## 🎨 Code Style Guidelines

### TypeScript

```tsx
// ✅ Good
interface User {
  id: string;
  name: string;
  email: string;
}

function getUser(id: string): User {
  // ...
}

// ❌ Bad
function getUser(id: any): any {
  // ...
}
```

### React Components

```tsx
// ✅ Good - Functional component with TypeScript
interface ButtonProps {
  label: string;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
}

export function Button({ label, onClick, variant = 'primary' }: ButtonProps) {
  return (
    <button className={styles[variant]} onClick={onClick}>
      {label}
    </button>
  );
}

// ❌ Bad - No types
export function Button(props) {
  return <button onClick={props.onClick}>{props.label}</button>;
}
```

### SCSS Modules

```scss
// ✅ Good
@import '../../styles/tokens';
@import '../../styles/mixins';

.container {
  padding: 1rem;
  background: var(--bg-primary);
  
  @include respond(sm) {
    padding: 0.5rem;
  }
}

// ❌ Bad - Hardcoded values
.container {
  padding: 16px;
  background: #ffffff;
}
```

### File Naming

- Components: `PascalCase.tsx` (e.g., `UserCard.tsx`)
- Pages: `PascalCase.tsx` (e.g., `Dashboard.tsx`)
- Utils: `camelCase.ts` (e.g., `formatDate.ts`)
- Styles: `ComponentName.module.scss`

---

## 🏗 Component Structure

### Tạo component mới

```
src/components/UserCard/
├── UserCard.tsx           # Component logic
├── UserCard.module.scss   # Styles
└── index.ts               # Export (optional)
```

### Template component

```tsx
// UserCard.tsx
import styles from './UserCard.module.scss';

interface UserCardProps {
  user: {
    id: string;
    name: string;
    email: string;
    avatar?: string;
  };
  onEdit?: (id: string) => void;
  onDelete?: (id: string) => void;
}

export function UserCard({ user, onEdit, onDelete }: UserCardProps) {
  return (
    <div className={styles.card}>
      {user.avatar && (
        <img src={user.avatar} alt={user.name} className={styles.avatar} />
      )}
      <div className={styles.info}>
        <h3 className={styles.name}>{user.name}</h3>
        <p className={styles.email}>{user.email}</p>
      </div>
      <div className={styles.actions}>
        {onEdit && (
          <button onClick={() => onEdit(user.id)}>Edit</button>
        )}
        {onDelete && (
          <button onClick={() => onDelete(user.id)}>Delete</button>
        )}
      </div>
    </div>
  );
}
```

```scss
// UserCard.module.scss
@import '../../styles/tokens';
@import '../../styles/mixins';

.card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  transition: all 0.2s ease;
  
  &:hover {
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  }
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.info {
  flex: 1;
}

.name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.25rem 0;
}

.email {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.actions {
  display: flex;
  gap: 0.5rem;
}
```

---

## 🎯 Best Practices

### 1. Component Design

- **Single Responsibility**: Mỗi component làm 1 việc
- **Reusability**: Thiết kế để tái sử dụng
- **Props Interface**: Luôn định nghĩa TypeScript interface
- **Default Props**: Có giá trị mặc định khi cần

### 2. State Management

- **Local State**: Dùng `useState` cho state component
- **Shared State**: Dùng Context API
- **Avoid Prop Drilling**: Dùng Context thay vì truyền props nhiều cấp

### 3. Performance

- **Memoization**: Dùng `useMemo`, `useCallback` khi cần
- **Lazy Loading**: Code splitting cho pages lớn
- **Virtualization**: Dùng cho lists dài (>100 items)

### 4. Styling

- **CSS Modules**: Scope styles cho component
- **Design Tokens**: Dùng CSS variables từ `_tokens.scss`
- **Responsive**: Mobile-first approach
- **Mixins**: Tái sử dụng patterns từ `_mixins.scss`

### 5. TypeScript

- **Strict Mode**: Bật strict TypeScript
- **No Any**: Tránh dùng `any` type
- **Type Inference**: Để TypeScript tự infer khi có thể
- **Interfaces**: Định nghĩa rõ ràng cho objects

---

## 📦 Thêm Dependencies

### Trước khi thêm

1. Kiểm tra xem đã có library tương tự chưa
2. Kiểm tra license
3. Kiểm tra bundle size
4. Đọc documentation

### Cài đặt

```bash
# Development dependency
npm install -D package-name

# Production dependency
npm install package-name
```

### Cập nhật package.json

Thêm comment giải thích tại sao cần package đó.

---

## 🧪 Testing Guidelines

### Manual Testing Checklist

Trước khi commit, kiểm tra:

- [ ] Component hiển thị đúng
- [ ] Responsive trên mobile/tablet/desktop
- [ ] Form validation hoạt động
- [ ] Error handling đúng
- [ ] Loading states hiển thị
- [ ] No console errors
- [ ] Lint pass: `npm run lint`
- [ ] Build success: `npm run build`

### Browser Testing

Test trên:
- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

---

## 🔍 Code Review

### Reviewer Checklist

- [ ] Code style đúng convention
- [ ] TypeScript types đầy đủ
- [ ] No hardcoded values
- [ ] Responsive design
- [ ] Performance optimization
- [ ] Error handling
- [ ] Comments cho logic phức tạp
- [ ] No console.log trong production code

### Submitter Checklist

- [ ] Self-review code
- [ ] Test thoroughly
- [ ] Update documentation
- [ ] Clear commit messages
- [ ] PR description đầy đủ

---

## 📚 Resources

### Design System
- Material Design 3: https://m3.material.io/
- Design Tokens: `src/styles/_tokens.scss`
- Components: `src/components/`

### TypeScript
- Handbook: https://www.typescriptlang.org/docs/
- React TypeScript: https://react-typescript-cheatsheet.netlify.app/

### React
- Docs: https://react.dev/
- Hooks: https://react.dev/reference/react

### SCSS
- Documentation: https://sass-lang.com/documentation/
- Mixins: `src/styles/_mixins.scss`

---

## 🎓 Learning Path

Cho developers mới:

1. Đọc [SETUP_GUIDE.md](./SETUP_GUIDE.md)
2. Đọc [QUICK_START.md](./QUICK_START.md)
3. Khám phá `src/components/` để hiểu các component cơ bản
4. Xem `src/pages/Dashboard/` như ví dụ page hoàn chỉnh
5. Thử tạo component đơn giản
6. Đọc `doc/frontend_development_guide.md`

---

## ❓ Questions?

Nếu có thắc mắc:

1. Kiểm tra documentation trong `doc/`
2. Hỏi team members
3. Tạo issue trên repository

---

**Thank you for contributing! 🙏**
