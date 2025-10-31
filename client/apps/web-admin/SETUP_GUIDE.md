# 📖 Setup Guide

## Prerequisites

- **Node.js** >= 18.x
- **npm** >= 9.x hoặc **pnpm**
- **Git**

## Installation

### 1. Clone Project
```bash
git clone https://github.com/youknow2509/cio_verify_face.git
cd cio_verify_face/client/apps/web-admin
```

### 2. Install Dependencies
```bash
npm install
# hoặc nếu dùng pnpm
pnpm install
```

### 3. Environment Configuration

```bash
cp .env.example .env
```

**Cấu hình `.env`:**
```env
# API Server
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=10000

# Feature Flags
VITE_ENABLE_MOCK_API=false

# App
VITE_APP_TITLE=Face Attendance System
VITE_APP_ENV=development
```

### 4. Start Development
```bash
npm run dev
```

Ứng dụng sẽ chạy tại: **http://localhost:3003**

## Build for Production

```bash
npm run build
npm run preview  # Preview production build
```

Output folder: `dist/`

## Project Structure

```
src/
├── pages/                 # Page components
├── components/           # Reusable UI components
├── services/
│   ├── api/             # API endpoint functions
│   ├── http.ts          # HTTP client with interceptors
│   ├── error-handler.ts # Error handling
│   └── api-helpers.ts   # Helper utilities
├── hooks/               # Custom React hooks
├── types/               # TypeScript type definitions
├── styles/              # Global styles
└── utils/               # Utility functions
```

## Database & Backend Setup

### Start Backend Services (Docker)

```bash
cd ../../.. # Go to project root
docker-compose -f server/docker-compose.yml up -d
```

Services started:
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- Minio: `localhost:9000`
- API Gateway: `localhost:8080`

### Check Backend Status
```bash
curl http://localhost:8080/api/v1/ping
```

## Development Commands

```bash
npm run dev          # Start dev server with HMR
npm run build        # Production build
npm run preview      # Preview production build
npm run lint         # ESLint check
npm run type-check   # TypeScript check
npm run format       # Format code with Prettier
```

## Debugging

### Browser DevTools
1. Open Chrome DevTools: `F12`
2. Check **Network** tab for API calls
3. Check **Console** for errors

### VS Code Debugging
Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "chrome",
      "request": "launch",
      "name": "Launch Chrome",
      "url": "http://localhost:3003",
      "webRoot": "${workspaceFolder}/client/apps/web-admin",
      "sourceMapPathOverride": {
        "/src/*": "${webspaceFolder}/src/*"
      }
    }
  ]
}
```

## Common Issues

### Port 3003 Already in Use
```bash
# Change port
npm run dev -- --port 3004

# Or find process using port
lsof -i :3003  # macOS/Linux
netstat -ano | findstr :3003  # Windows
```

### API Connection Error
1. Verify backend is running: `curl http://localhost:8080/api/v1/ping`
2. Check `VITE_API_BASE_URL` in `.env`
3. Check CORS settings in backend

### Module Not Found
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

### TypeScript Errors
```bash
npm run type-check
npm run build
```

## Environment Presets

### Development
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_ENABLE_MOCK_API=true
VITE_APP_ENV=development
```

### Staging
```env
VITE_API_BASE_URL=https://api-staging.example.com
VITE_ENABLE_MOCK_API=false
VITE_APP_ENV=staging
```

### Production
```env
VITE_API_BASE_URL=https://api.example.com
VITE_ENABLE_MOCK_API=false
VITE_APP_ENV=production
```

## Performance Tips

- Use **dev tools** to monitor performance
- Enable **mock API** for frontend-only development
- Use **React DevTools** extension
- Profile with **Lighthouse** in Chrome

## Support

For issues, check:
1. [QUICK_START.md](QUICK_START.md) - Quick setup
2. [API_GUIDE.md](src/services/API_GUIDE.md) - API integration
3. Create GitHub issue with details
