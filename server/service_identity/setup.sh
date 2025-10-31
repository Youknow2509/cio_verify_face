#!/bin/bash

# Setup script for Identity & Organization Service

echo "🚀 Setting up Identity & Organization Service..."

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
  echo "📦 Installing dependencies..."
  npm install
else
  echo "✅ Dependencies already installed"
fi

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
  echo "📝 Creating .env file from .env.example..."
  cp .env.example .env
  echo "⚠️  Please configure .env file with your database credentials"
else
  echo "✅ .env file already exists"
fi

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Configure .env with your database credentials"
echo "2. Run migrations: goose -dir sql postgres \"{DB_CONNECTION_STRING}\" up"
echo "3. Start development server: npm run dev"
echo ""
