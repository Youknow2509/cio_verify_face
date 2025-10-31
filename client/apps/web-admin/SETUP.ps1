#!/usr/bin/env pwsh
# 🚀 Web-Admin API Integration - Quick Start Guide (PowerShell)

Write-Host ""
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "🎯 Web-Admin API Integration Setup" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Copy environment template
Write-Host "📝 Step 1: Setting up environment variables..." -ForegroundColor Yellow
if (-not (Test-Path ".env.local")) {
    Copy-Item ".env.example" ".env.local"
    Write-Host "✅ Created .env.local from .env.example" -ForegroundColor Green
    Write-Host "   Please edit .env.local with your API configuration" -ForegroundColor Gray
}
else {
    Write-Host "✅ .env.local already exists" -ForegroundColor Green
}
Write-Host ""

# Step 2: Check Node modules
Write-Host "📦 Step 2: Checking dependencies..." -ForegroundColor Yellow
if (-not (Test-Path "node_modules")) {
    Write-Host "⚠️  node_modules not found. Installing dependencies..." -ForegroundColor Yellow
    & npm install
    Write-Host "✅ Dependencies installed" -ForegroundColor Green
}
else {
    Write-Host "✅ Dependencies already installed" -ForegroundColor Green
}
Write-Host ""

# Step 3: Check TypeScript
Write-Host "🔍 Step 3: Running type check..." -ForegroundColor Yellow
& npm run type-check
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Type checking failed - please fix errors" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Type checking passed" -ForegroundColor Green
Write-Host ""

# Step 4: Print summary
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "✅ Setup Complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "📚 Documentation:" -ForegroundColor Cyan
Write-Host "  - API Guide: src/services/API_GUIDE.md"
Write-Host "  - Services Structure: src/services/SERVICES_STRUCTURE.md"
Write-Host "  - Migration Guide: MIGRATION_GUIDE.md"
Write-Host "  - Integration Checklist: API_INTEGRATION_CHECKLIST.md"
Write-Host ""

Write-Host "🔧 Configuration:" -ForegroundColor Cyan
Write-Host "  - Check .env.local for API_BASE_URL and other settings"
Write-Host "  - API Timeout: 30000ms (by default)"
Write-Host "  - Logging: Enabled (by default)"
Write-Host ""

Write-Host "🚀 Start development:" -ForegroundColor Green
Write-Host "  npm run dev"
Write-Host ""

Write-Host "📖 Next Steps:" -ForegroundColor Cyan
Write-Host "  1. Read SERVICES_STRUCTURE.md to understand the API layer"
Write-Host "  2. Read API_GUIDE.md for detailed usage"
Write-Host "  3. Follow MIGRATION_GUIDE.md to migrate features"
Write-Host "  4. Use API_INTEGRATION_CHECKLIST.md to track progress"
Write-Host ""
