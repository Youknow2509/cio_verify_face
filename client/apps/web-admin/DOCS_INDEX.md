# 📑 Web-Admin Documentation Index

Danh sách các tài liệu hướng dẫn trong dự án web-admin.

## 🎯 Getting Started

Start here if you're new to the project:

1. **[README.md](README.md)** ⭐

   - Project overview
   - Feature list
   - Quick links to all docs
   - **⏱️ 5 minutes**

2. **[QUICK_START.md](QUICK_START.md)** ⚡
   - 3-step installation
   - Default credentials
   - Common issues
   - **⏱️ 5 minutes**

## 🔧 Installation & Setup

For detailed setup instructions:

- **[SETUP_GUIDE.md](SETUP_GUIDE.md)**
  - Prerequisites
  - Full installation steps
  - Environment configuration
  - Docker setup
  - Debugging tips
  - **⏱️ 15 minutes**

## 🔌 API Integration

For working with APIs:

- **[src/services/API_GUIDE.md](src/services/API_GUIDE.md)**

  - API services overview
  - All 8 services documented
  - Usage examples
  - Error handling
  - Best practices
  - **⏱️ 20 minutes**

- **[MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)**
  - Switching from mock to real API
  - Testing checklist
  - Common issues
  - Debugging tips
  - **⏱️ 10 minutes**

## 👨‍💻 Development

For contributing to the project:

- **[CONTRIBUTING.md](CONTRIBUTING.md)**
  - Code standards
  - TypeScript conventions
  - Git commit format
  - PR process
  - Testing guidelines
  - **⏱️ 10 minutes**

## 📝 Reference

- **[CHANGELOG.md](CHANGELOG.md)**
  - Version history
  - New features per version
  - Breaking changes
  - Migration guides

## 📁 Project Structure

```
web-admin/
├── src/
│   ├── pages/              # React pages
│   ├── components/         # Reusable components
│   ├── services/
│   │   ├── api/           # API functions (8 services)
│   │   ├── http.ts        # HTTP client
│   │   ├── error-handler.ts
│   │   ├── api-helpers.ts
│   │   └── API_GUIDE.md   # 📖
│   ├── types/             # TypeScript types
│   ├── hooks/             # Custom hooks
│   ├── styles/            # Global styles
│   └── utils/             # Helpers
├── README.md              # 📖
├── QUICK_START.md         # 📖
├── SETUP_GUIDE.md         # 📖
├── CONTRIBUTING.md        # 📖
├── MIGRATION_GUIDE.md     # 📖
├── CHANGELOG.md           # 📖
├── .env.example
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 🔍 Finding What You Need

### I want to...

**...start developing**
→ [QUICK_START.md](QUICK_START.md) + [SETUP_GUIDE.md](SETUP_GUIDE.md)

**...understand the API**
→ [src/services/API_GUIDE.md](src/services/API_GUIDE.md)

**...connect to real backend**
→ [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)

**...contribute code**
→ [CONTRIBUTING.md](CONTRIBUTING.md)

**...check what's new**
→ [CHANGELOG.md](CHANGELOG.md)

**...see feature overview**
→ [README.md](README.md)

**...troubleshoot issues**
→ [SETUP_GUIDE.md](SETUP_GUIDE.md) (Debugging section)

## 📚 Services Documentation

All API services in `src/services/api/`:

| Service           | Endpoint                             | Purpose                  |
| ----------------- | ------------------------------------ | ------------------------ |
| auth.api.ts       | `/api/v1/auth`                       | Authentication & Device  |
| employees.api.ts  | `/api/v1/users`                      | User/Employee Management |
| devices.api.ts    | `/api/v1/devices`                    | Device Management        |
| attendance.api.ts | `/api/v1/attendance`                 | Check-in/Check-out       |
| shifts.api.ts     | `/api/v1/shifts` `/api/v1/schedules` | Shifts & Schedules       |
| reports.api.ts    | `/api/v1/reports`                    | Reports & Analytics      |
| signatures.api.ts | `/api/v1/signatures`                 | Signature Upload         |
| account.api.ts    | Auth + Users                         | Account Settings         |

👉 See [API_GUIDE.md](src/services/API_GUIDE.md) for full details

## 🆘 Quick Help

### Installation Issues?

- Check [SETUP_GUIDE.md](SETUP_GUIDE.md) - Troubleshooting section

### API Errors?

- Check [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Common Issues section
- Check [API_GUIDE.md](src/services/API_GUIDE.md) - Error Handling section

### Code Questions?

- Check [CONTRIBUTING.md](CONTRIBUTING.md) - Code Standards section

### Feature Questions?

- Check [README.md](README.md) - Features section

## 📞 Support

For issues:

1. Check relevant documentation above
2. Search GitHub issues
3. Create new issue with:
   - Steps to reproduce
   - Error message/screenshot
   - Environment details

## 🎓 Learning Path

**Beginner:**

1. README.md
2. QUICK_START.md
3. SETUP_GUIDE.md

**Developer:** 4. API_GUIDE.md 5. CONTRIBUTING.md 6. Review src/services/

**Advanced:** 7. MIGRATION_GUIDE.md 8. CHANGELOG.md 9. Source code exploration

---

**Last Updated:** 2025-10-31
**Version:** 1.1.0
**Status:** Documentation cleaned up and simplified ✅
