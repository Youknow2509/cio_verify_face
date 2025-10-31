#!/usr/bin/env bash
# 🚀 Web-Admin API Integration - Quick Start Guide

echo "======================================"
echo "🎯 Web-Admin API Integration Setup"
echo "======================================"
echo ""

# Step 1: Copy environment template
echo "📝 Step 1: Setting up environment variables..."
if [ ! -f ".env.local" ]; then
  cp .env.example .env.local
  echo "✅ Created .env.local from .env.example"
  echo "   Please edit .env.local with your API configuration"
else
  echo "✅ .env.local already exists"
fi
echo ""

# Step 2: Check Node modules
echo "📦 Step 2: Checking dependencies..."
if [ ! -d "node_modules" ]; then
  echo "⚠️  node_modules not found. Installing dependencies..."
  npm install
  echo "✅ Dependencies installed"
else
  echo "✅ Dependencies already installed"
fi
echo ""

# Step 3: Check TypeScript
echo "🔍 Step 3: Running type check..."
npm run type-check
if [ $? -eq 0 ]; then
  echo "✅ Type checking passed"
else
  echo "❌ Type checking failed - please fix errors"
  exit 1
fi
echo ""

# Step 4: Print summary
echo "======================================"
echo "✅ Setup Complete!"
echo "======================================"
echo ""
echo "📚 Documentation:"
echo "  - API Guide: src/services/API_GUIDE.md"
echo "  - Services Structure: src/services/SERVICES_STRUCTURE.md"
echo "  - Migration Guide: MIGRATION_GUIDE.md"
echo "  - Integration Checklist: API_INTEGRATION_CHECKLIST.md"
echo ""
echo "🔧 Configuration:"
echo "  - API Base URL: Check .env.local (VITE_API_BASE_URL)"
echo "  - API Timeout: ${VITE_API_TIMEOUT:-30000}ms"
echo "  - Logging: ${VITE_ENABLE_API_LOGGING:-true}"
echo ""
echo "🚀 Start development:"
echo "  npm run dev"
echo ""
echo "📖 Next Steps:"
echo "  1. Read SERVICES_STRUCTURE.md to understand the API layer"
echo "  2. Read API_GUIDE.md for detailed usage"
echo "  3. Follow MIGRATION_GUIDE.md to migrate features"
echo "  4. Use API_INTEGRATION_CHECKLIST.md to track progress"
echo ""
