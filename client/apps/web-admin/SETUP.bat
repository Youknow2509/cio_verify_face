@echo off
REM 🚀 Web-Admin API Integration - Quick Start Guide (Windows)

setlocal enabledelayedexpansion

echo.
echo ======================================
echo 🎯 Web-Admin API Integration Setup
echo ======================================
echo.

REM Step 1: Copy environment template
echo 📝 Step 1: Setting up environment variables...
if not exist ".env.local" (
    copy .env.example .env.local
    echo ✅ Created .env.local from .env.example
    echo    Please edit .env.local with your API configuration
) else (
    echo ✅ .env.local already exists
)
echo.

REM Step 2: Check Node modules
echo 📦 Step 2: Checking dependencies...
if not exist "node_modules" (
    echo ⚠️  node_modules not found. Installing dependencies...
    call npm install
    echo ✅ Dependencies installed
) else (
    echo ✅ Dependencies already installed
)
echo.

REM Step 3: Check TypeScript
echo 🔍 Step 3: Running type check...
call npm run type-check
if errorlevel 1 (
    echo ❌ Type checking failed - please fix errors
    pause
    exit /b 1
)
echo ✅ Type checking passed
echo.

REM Step 4: Print summary
echo ======================================
echo ✅ Setup Complete!
echo ======================================
echo.
echo 📚 Documentation:
echo   - API Guide: src/services/API_GUIDE.md
echo   - Services Structure: src/services/SERVICES_STRUCTURE.md
echo   - Migration Guide: MIGRATION_GUIDE.md
echo   - Integration Checklist: API_INTEGRATION_CHECKLIST.md
echo.
echo 🔧 Configuration:
echo   - Check .env.local for API_BASE_URL and other settings
echo   - API Timeout: 30000ms (by default)
echo   - Logging: Enabled (by default)
echo.
echo 🚀 Start development:
echo   npm run dev
echo.
echo 📖 Next Steps:
echo   1. Read SERVICES_STRUCTURE.md to understand the API layer
echo   2. Read API_GUIDE.md for detailed usage
echo   3. Follow MIGRATION_GUIDE.md to migrate features
echo   4. Use API_INTEGRATION_CHECKLIST.md to track progress
echo.
pause
