#!/usr/bin/env node
/**
 * 📋 Web-Admin API Integration - Final Checklist
 * 
 * Run this file to verify everything is set up correctly
 * Usage: node CHECK_SETUP.js
 */

const fs = require('fs');
const path = require('path');

const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function checkFile(filePath, description) {
  const fullPath = path.join(__dirname, filePath);
  const exists = fs.existsSync(fullPath);
  const status = exists ? '✅' : '❌';
  const color = exists ? 'green' : 'red';
  log(`${status} ${description}`, color);
  return exists;
}

function checkFiles(files, category) {
  log(`\n📁 ${category}`, 'cyan');
  log('─'.repeat(50), 'cyan');
  
  let passed = 0;
  let failed = 0;
  
  files.forEach(({ path: filePath, description }) => {
    if (checkFile(filePath, description)) {
      passed++;
    } else {
      failed++;
    }
  });
  
  return { passed, failed };
}

function main() {
  log('\n', 'blue');
  log('╔════════════════════════════════════════════════════╗', 'blue');
  log('║  📋 Web-Admin API Integration - Setup Checklist    ║', 'blue');
  log('╚════════════════════════════════════════════════════╝', 'blue');
  
  let totalPassed = 0;
  let totalFailed = 0;
  
  // API Services
  const apiFiles = [
    { path: 'src/services/api/auth.api.ts', description: 'Auth API service' },
    { path: 'src/services/api/employees.api.ts', description: 'Employees API service' },
    { path: 'src/services/api/devices.api.ts', description: 'Devices API service' },
    { path: 'src/services/api/attendance.api.ts', description: 'Attendance API service' },
    { path: 'src/services/api/shifts.api.ts', description: 'Shifts API service' },
    { path: 'src/services/api/account.api.ts', description: 'Account API service' },
    { path: 'src/services/api/index.ts', description: 'API services index' }
  ];
  let result = checkFiles(apiFiles, 'API Services (7 files)');
  totalPassed += result.passed;
  totalFailed += result.failed;
  
  // Infrastructure
  const infraFiles = [
    { path: 'src/services/http.ts', description: 'HTTP client' },
    { path: 'src/services/http-interceptor.ts', description: 'HTTP interceptor' },
    { path: 'src/services/error-handler.ts', description: 'Error handler' },
    { path: 'src/services/api-helpers.ts', description: 'API helpers' },
    { path: 'src/config/api.config.ts', description: 'API configuration' }
  ];
  result = checkFiles(infraFiles, 'Infrastructure Files (5 files)');
  totalPassed += result.passed;
  totalFailed += result.failed;
  
  // Mock Services
  const mockFiles = [
    { path: 'src/services/mock/index.ts', description: 'Mock services index' }
  ];
  result = checkFiles(mockFiles, 'Mock Services (1 file)');
  totalPassed += result.passed;
  totalFailed += result.failed;
  
  // Documentation
  const docFiles = [
    { path: 'src/services/API_GUIDE.md', description: 'API usage guide' },
    { path: 'src/services/SERVICES_STRUCTURE.md', description: 'Services structure guide' },
    { path: 'MIGRATION_GUIDE.md', description: 'Migration guide' },
    { path: 'API_INTEGRATION_CHECKLIST.md', description: 'Integration checklist' },
    { path: 'WEB_ADMIN_REFACTORING_SUMMARY.md', description: 'Refactoring summary' },
    { path: 'NEW_FILES_SUMMARY.md', description: 'New files summary' },
    { path: 'COMPLETION_SUMMARY.md', description: 'Completion summary' },
    { path: 'DOCUMENTATION_INDEX.md', description: 'Documentation index' }
  ];
  result = checkFiles(docFiles, 'Documentation (8 files)');
  totalPassed += result.passed;
  totalFailed += result.failed;
  
  // Configuration
  const configFiles = [
    { path: '.env.example', description: 'Environment template' },
    { path: 'SETUP.sh', description: 'Setup script (bash)' },
    { path: 'SETUP.bat', description: 'Setup script (batch)' },
    { path: 'SETUP.ps1', description: 'Setup script (powershell)' }
  ];
  result = checkFiles(configFiles, 'Configuration (4 files)');
  totalPassed += result.passed;
  totalFailed += result.failed;
  
  // Environment check
  log('\n🔧 Environment Variables', 'cyan');
  log('─'.repeat(50), 'cyan');
  const envLocal = fs.existsSync(path.join(__dirname, '.env.local'));
  if (envLocal) {
    log('✅ .env.local exists', 'green');
  } else {
    log('⚠️  .env.local not found (copy from .env.example)', 'yellow');
  }
  
  // Dependencies check
  log('\n📦 Dependencies', 'cyan');
  log('─'.repeat(50), 'cyan');
  const nodeModules = fs.existsSync(path.join(__dirname, 'node_modules'));
  if (nodeModules) {
    log('✅ node_modules exists', 'green');
  } else {
    log('❌ node_modules not found (run npm install)', 'red');
    totalFailed++;
  }
  
  // Summary
  log('\n📊 Summary', 'cyan');
  log('─'.repeat(50), 'cyan');
  log(`✅ Passed: ${totalPassed}`, 'green');
  log(`❌ Failed: ${totalFailed}`, totalFailed > 0 ? 'red' : 'green');
  
  if (totalFailed === 0) {
    log('\n🎉 All checks passed! Ready for integration.', 'green');
    log('\n📖 Next Steps:', 'blue');
    log('1. Read DOCUMENTATION_INDEX.md', 'yellow');
    log('2. Review SERVICES_STRUCTURE.md', 'yellow');
    log('3. Read API_GUIDE.md for usage', 'yellow');
    log('4. Follow MIGRATION_GUIDE.md to integrate features', 'yellow');
    log('5. Track progress with API_INTEGRATION_CHECKLIST.md', 'yellow');
    log('\n🚀 Start development: npm run dev', 'green');
  } else {
    log('\n⚠️  Some files are missing. Please run setup scripts.', 'yellow');
    log('\n🔧 Run one of:', 'blue');
    log('   Linux/Mac: ./SETUP.sh', 'yellow');
    log('   Windows:   SETUP.bat', 'yellow');
    log('   PowerShell: .\\SETUP.ps1', 'yellow');
  }
  
  log('\n', 'reset');
}

main();
