# 📝 Changelog

All notable changes to this project are documented here.

## [1.1.0] - 2025-10-31

### ✨ Added
- **API Services Refactoring**
  - Updated all API endpoints to use `/api/v1/` prefix
  - Added device activation endpoint
  - Added check-in/check-out endpoints for attendance
  - Added schedule management (shifts + schedules)
  - Added reports service (daily, summary, export)
  - Added signatures service

- **New Services**
  - `reports.api.ts` - Analytics & Reporting Service
  - `signatures.api.ts` - Signature Upload Service

- **Types**
  - Added `Schedule` interface for shift scheduling

- **Documentation**
  - Refactored and condensed documentation
  - Updated API integration guide
  - Added detailed setup instructions

### 🔄 Changed
- **API Endpoints**
  - Auth: `/auth/` → `/api/v1/auth/`
  - Users: `/employees/` → `/api/v1/users/`
  - Devices: `/devices/` → `/api/v1/devices/`
  - Attendance: `/attendance/` → `/api/v1/attendance/`
  - Shifts: `/shifts/` → `/api/v1/shifts/`

- **Services Structure**
  - `employees.api.ts` now maps to `/api/v1/users`
  - `attendance.api.ts` redesigned for face recognition
  - `shifts.api.ts` now includes schedule management

- **Account Service**
  - Now uses Auth Service endpoints
  - Profile updates via Users Service
  - Password change via Auth Service

### 🗑️ Removed
- Removed duplicate documentation files (11 files cleaned up)
- Removed old README_old.md
- Removed auth endpoints not in specification

### 🐛 Fixed
- Correct API endpoint paths throughout codebase
- Proper tenant isolation in API calls
- Type safety improvements

## [1.0.0] - 2025-10-30

### ✨ Initial Release
- Complete API services layer
- Infrastructure components (HTTP client, error handler)
- Comprehensive documentation
- Setup scripts for all platforms
- Mock API support for development
- Dashboard with charts and statistics
- Employee management CRUD
- Attendance tracking
- Reports generation
- Shift management

---

## Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** - Breaking changes
- **MINOR** - New features (backward compatible)
- **PATCH** - Bug fixes

## Migration Guides

- [1.0.0 → 1.1.0](MIGRATION_GUIDE.md) - API endpoint updates
