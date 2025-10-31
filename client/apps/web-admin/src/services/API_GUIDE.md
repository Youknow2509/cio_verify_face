# 🔌 API Services Guide

Hướng dẫn sử dụng API Services layer trong web-admin.

## 📁 Cấu trúc

```
src/services/
├── api/                    # API endpoint functions
│   ├── auth.api.ts        # Auth & Device
│   ├── employees.api.ts   # Users/Employees
│   ├── devices.api.ts     # Device management
│   ├── attendance.api.ts  # Attendance (check-in/out)
│   ├── shifts.api.ts      # Shifts & Schedules
│   ├── reports.api.ts     # Reports
│   ├── signatures.api.ts  # Signatures
│   ├── account.api.ts     # Account settings
│   └── index.ts           # Central exports
├── http.ts                # HTTP client
├── error-handler.ts       # Error handling
└── api-helpers.ts         # Helper utilities
```

## 🚀 Quick Usage

### Import API Functions
```typescript
import {
  loginAPI,
  getEmployeesAPI,
  checkInAPI,
  getDailyReportAPI,
} from '@/services/api';
```

### Call API
```typescript
// Login
const result = await loginAPI('user@example.com', 'password');
if (!result.error) {
  const { data } = result;
  localStorage.setItem('token', data.token);
}

// Get employees
const employees = await getEmployeesAPI({ page: 1, limit: 10 });
if (!employees.error) {
  console.log(`Total: ${employees.data.total}`);
}

// Check in
const checkIn = await checkInAPI(faceImage);
if (!checkIn.error) {
  console.log('Check in successful');
}
```

## 📝 API Services Overview

### Auth Service (`/api/v1/auth`)
```typescript
loginAPI(email, password)
logoutAPI()
refreshTokenAPI()
getCurrentUserAPI()
activateDeviceAPI(deviceCode, secret)
```

### Users Service (`/api/v1/users`)
```typescript
getEmployeesAPI(filter)              // List
getEmployeeAPI(userId)               // Detail
createEmployeeAPI(data)              // Create
updateEmployeeAPI(userId, data)      // Update
deleteEmployeeAPI(userId)            // Delete
uploadFaceDataAPI(userId, file)      // Face data
getEmployeeFaceDataAPI(userId)       // Get faces
deleteFaceDataAPI(userId, faceId)    // Delete face
```

### Devices Service (`/api/v1/devices`)
```typescript
getDevicesAPI(page, limit, search)
getDeviceAPI(deviceId)
createDeviceAPI(data)
updateDeviceAPI(deviceId, data)
deleteDeviceAPI(deviceId)
syncDeviceAPI(deviceId)
```

### Attendance Service (`/api/v1/attendance`)
```typescript
checkInAPI(faceImage)
checkOutAPI(faceImage)
getAttendanceRecordsAPI(filter)
getAttendanceRecordAPI(recordId)
getMyAttendanceHistoryAPI(filter)
```

### Shifts Service (`/api/v1/shifts` & `/api/v1/schedules`)
```typescript
// Shifts
getShiftsAPI(page, limit, active)
getShiftAPI(shiftId)
createShiftAPI(data)
updateShiftAPI(shiftId, data)
deleteShiftAPI(shiftId)

// Schedules
getSchedulesAPI(page, limit, userId)
getScheduleAPI(scheduleId)
createScheduleAPI(data)
updateScheduleAPI(scheduleId, data)
deleteScheduleAPI(scheduleId)
```

### Reports Service (`/api/v1/reports`)
```typescript
getDailyReportAPI(date, department)
getSummaryReportAPI(startDate, endDate, department)
exportReportAPI(params)
```

### Signatures Service (`/api/v1/signatures`)
```typescript
uploadSignatureAPI(file)
getUserSignatureAPI(userId)
deleteSignatureAPI(signatureId)
```

### Account Service
```typescript
getAccountProfileAPI()            // /api/v1/auth/me
updateAccountProfileAPI(data)     // /api/v1/users/{id}
changeAccountPasswordAPI(payload) // /api/v1/auth/change-password
uploadAccountAvatarAPI(file)      // /api/v1/users/{id}/avatar
```

## 🔐 Authentication

```typescript
import { loginAPI, setAuthToken } from '@/services/api';

// Login
const response = await loginAPI(email, password);
if (response.data?.token) {
  // Token automatically added to all requests
  setAuthToken(response.data.token);
  localStorage.setItem('token', response.data.token);
}

// Logout
await logoutAPI();
localStorage.removeItem('token');
```

Token is automatically included in `Authorization: Bearer <token>` header.

## ⚠️ Error Handling

```typescript
import { handleApiError, isAuthError } from '@/services/error-handler';

try {
  const result = await getEmployeesAPI();
  
  if (result.error) {
    const errorInfo = handleApiError(new Error(result.error));
    
    if (isAuthError(errorInfo)) {
      // Redirect to login
      window.location.href = '/login';
    }
  }
} catch (error) {
  const errorInfo = handleApiError(error);
  console.error('API Error:', errorInfo.message);
}
```

## 🔧 Configuration

`.env` file:
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=10000
VITE_ENABLE_MOCK_API=false
```

## 📊 Response Format

```typescript
// Success
{
  data: { /* actual data */ },
  error: null,
  message: 'Success'
}

// Error
{
  data: null,
  error: 'Error message',
  message: 'Failed'
}

// Paginated
{
  data: {
    data: [ /* items */ ],
    total: 100,
    page: 1,
    limit: 10,
    totalPages: 10
  }
}
```

## 💡 Best Practices

1. **Always check error field:**
   ```typescript
   if (result.error) {
     // Handle error
   } else {
     // Use result.data
   }
   ```

2. **Use TypeScript types:**
   ```typescript
   const employees = await getEmployeesAPI();
   // Type is inferred: ApiResponse<PaginatedResponse<Employee>>
   ```

3. **Handle API timeout:**
   ```typescript
   try {
     const result = await getEmployeesAPI();
   } catch (error) {
     // Timeout or network error
   }
   ```

4. **Don't hardcode URLs:**
   ```typescript
   // ❌ Bad
   const response = await fetch('http://localhost:8080/api/v1/users');
   
   // ✅ Good
   const response = await getEmployeesAPI();
   ```

5. **Filter and pagination:**
   ```typescript
   const result = await getEmployeesAPI({
     page: 1,
     limit: 20,
     search: 'John',
     department: 'Sales'
   });
   ```

## 🔗 Related Files

- [QUICK_START.md](../../QUICK_START.md) - Getting started
- [MIGRATION_GUIDE.md](../../MIGRATION_GUIDE.md) - Mock → Real API
- [CONTRIBUTING.md](../../CONTRIBUTING.md) - Code standards
