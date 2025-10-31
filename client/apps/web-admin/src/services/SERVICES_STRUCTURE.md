// src/services/SERVICES_STRUCTURE.md

# Services Structure Guide

## 📁 Overview

```
src/services/
├── api/                      # API Service Layer (NEW)
│   ├── auth.api.ts          # Authentication endpoints
│   ├── employees.api.ts     # Employee management
│   ├── devices.api.ts       # Device management
│   ├── attendance.api.ts    # Attendance tracking
│   ├── shifts.api.ts        # Shift management
│   ├── account.api.ts       # Account management
│   └── index.ts             # Centralized exports
│
├── mock/                     # Mock Services (Legacy)
│   ├── auth.ts
│   ├── employees.ts
│   ├── devices.ts
│   ├── attendance.ts
│   ├── shifts.ts
│   ├── reports.ts
│   └── index.ts
│
├── http.ts                   # HTTP Client (Enhanced)
├── http-interceptor.ts       # Interceptors (NEW)
├── error-handler.ts          # Error Handling (NEW)
├── api-helpers.ts            # Helper Utilities (NEW)
├── index.ts                  # Main exports
└── API_GUIDE.md             # Usage guide
```

---

## 🔄 Request Flow

```
Component
    ↓
import { getEmployeesAPI } from '@/services/api'
    ↓
getEmployeesAPI(filters)
    ↓
http.get() → createRequestConfig()
    ↓
add Authorization header
    ↓
fetch(url, config)
    ↓
Response handling
    ↓
return { data?, error? }
    ↓
Component handles result
```

---

## 📋 API Service Pattern

Each API service file follows this pattern:

```typescript
// src/services/api/[module].api.ts

import { http } from "../http";
import type { ApiResponse, DataType } from "@/types";

const API_PREFIX = "/[module]";

/**
 * Description of what this function does
 * @param params Parameter descriptions
 * @returns Return type
 */
export async function functionNameAPI(params): Promise<ApiResponse<T>> {
  return http.method(`${API_PREFIX}/endpoint`, params);
}
```

### Key Points

- All functions end with `API` suffix
- Always return `ApiResponse<T>` type
- Document with JSDoc comments
- Use environment-based API_PREFIX
- Consistent naming conventions

---

## 🎯 Using Services in Components

### 1. Import

```typescript
import { getEmployeesAPI } from "@/services/api";
import {
  handleApiError,
  getUserFriendlyMessage,
} from "@/services/error-handler";
```

### 2. Use in Effect

```typescript
useEffect(() => {
  const loadData = async () => {
    setLoading(true);

    const result = await getEmployeesAPI({ page: 1, limit: 10 });

    if (result.error) {
      const errorInfo = handleApiError(new Error(result.error));
      setError(getUserFriendlyMessage(errorInfo));
    } else {
      setEmployees(result.data.data);
    }

    setLoading(false);
  };

  loadData();
}, []);
```

### 3. Use in Handler

```typescript
const handleCreate = async (formData) => {
  setLoading(true);

  const result = await createEmployeeAPI(formData);

  if (result.error) {
    setError(result.error);
  } else {
    showSuccess("Created successfully");
    onCreated(result.data);
  }

  setLoading(false);
};
```

---

## 🔐 Authentication

### Set Token

```typescript
import { setAuthToken } from "@/services/http";

const loginResult = await loginAPI(email, password);
if (!loginResult.error) {
  setAuthToken(loginResult.data.token);
}
```

### Clear Token

```typescript
import { clearAuthToken } from "@/services/http";

await logoutAPI();
clearAuthToken();
```

Token is automatically added to all subsequent requests:

```
Authorization: Bearer {token}
```

---

## ⚠️ Error Handling

### Error Types

```typescript
import {
  isAuthError,
  isValidationError,
  isNotFoundError,
  isServerError,
} from "@/services/error-handler";

const errorInfo = handleApiError(error);

if (isAuthError(errorInfo)) {
  // Redirect to login
} else if (isValidationError(errorInfo)) {
  // Show validation errors
} else if (isNotFoundError(errorInfo)) {
  // Show 404 message
} else if (isServerError(errorInfo)) {
  // Show server error message
}
```

### User-Friendly Messages

```typescript
import { getUserFriendlyMessage } from "@/services/error-handler";

const message = getUserFriendlyMessage(errorInfo);
// "Authentication failed. Please login again."
// "Invalid data. Please check your input."
// etc.
```

---

## 🔧 Helper Functions

### Form Data Builders

```typescript
import {
  createEmployeePayload,
  createDevicePayload,
  createShiftPayload,
} from "@/services/api-helpers";

const payload = createEmployeePayload(formData);
await createEmployeeAPI(payload);
```

### Query Builders

```typescript
import { buildQueryString, buildFilterQuery } from "@/services/api-helpers";

const query = buildQueryString({ page: 1, search: "john" });
// "page=1&search=john"

const params = buildFilterQuery(filters);
// URLSearchParams object
```

### File Validation

```typescript
import { isValidImageFile, isValidFileSize } from "@/services/api-helpers";

if (!isValidImageFile(file)) {
  setError("Invalid image file");
  return;
}

if (!isValidFileSize(file, 10)) {
  // 10 MB
  setError("File too large");
  return;
}

await uploadFaceDataAPI(employeeId, file);
```

### Date Formatting

```typescript
import { formatDateForAPI, formatDatetimeForAPI } from "@/services/api-helpers";

const date = formatDateForAPI(new Date()); // "2024-10-31"
const datetime = formatDatetimeForAPI(new Date()); // "2024-10-31T10:30:00Z"
```

---

## 📝 Adding New API Endpoints

### Step 1: Create API File (if new module)

```typescript
// src/services/api/reports.api.ts
import { http } from "../http";
import type { ApiResponse } from "@/types";

const API_PREFIX = "/reports";

export async function getReportsAPI(): Promise<ApiResponse<Report[]>> {
  return http.get(API_PREFIX);
}
```

### Step 2: Export from Index

```typescript
// src/services/api/index.ts
export { getReportsAPI } from "./reports.api";
```

### Step 3: Update Main Services Index

```typescript
// src/services/index.ts
export * from "./api";
```

### Step 4: Use in Component

```typescript
import { getReportsAPI } from "@/services/api";

const result = await getReportsAPI();
```

---

## 🧪 Testing with Mock Services

Mock services are still available for fallback:

```typescript
// Use mock in development
import { getEmployees as getEmployeesMock } from "@/services/mock/employees";

const mockResult = await getEmployeesMock();
// { data: [...], error?: string }
```

Or use environment flag:

```typescript
// .env.development.local
VITE_ENABLE_MOCK_MODE = true;
```

---

## 🔍 Logging

### Enable/Disable

```typescript
// .env.local
VITE_ENABLE_API_LOGGING = true;
```

### Configure Programmatically

```typescript
import { configureHttpInterceptor } from "@/services/http-interceptor";

configureHttpInterceptor({
  enableLogging: true,
  enableRetry: true,
  maxRetries: 3,
});
```

### Log Output

```
📤 GET /api/employees
Body: undefined

✅ GET /api/employees [200]
Response: { data: [...], total: 50 }

❌ POST /api/employees [400]
Error: Invalid data
```

---

## 📊 Response Patterns

### Single Resource

```typescript
interface ApiResponse<T> {
  data: T;
  error?: string;
}

const result = await getEmployeeAPI("123");
// result.data: Employee
// result.error: string | undefined
```

### Paginated Collection

```typescript
interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

const result = await getEmployeesAPI({ page: 1, limit: 10 });
// result.data: PaginatedResponse<Employee>
// result.data.data: Employee[]
// result.data.total: number
```

### File Upload

```typescript
const formData = new FormData();
formData.append("file", file);

const result = await uploadFaceDataAPI(employeeId, file);
// result.data: FaceData
// result.error: string | undefined
```

---

## 🎓 Best Practices

### ✅ DO

- ✅ Always check `error` field first
- ✅ Use TypeScript types
- ✅ Handle all error cases
- ✅ Show loading states
- ✅ Provide user feedback
- ✅ Use helper utilities
- ✅ Document complex logic
- ✅ Test API integration

### ❌ DON'T

- ❌ Hardcode API URLs
- ❌ Ignore error responses
- ❌ Forget to set auth token
- ❌ Use `any` types
- ❌ Mix mock and real API calls
- ❌ Forget pagination handling
- ❌ Silently fail on errors
- ❌ Forget cleanup in effects

---

## 🔗 References

- [API_GUIDE.md](./API_GUIDE.md) - Detailed API usage
- [MIGRATION_GUIDE.md](../MIGRATION_GUIDE.md) - From mock to API
- [api.config.ts](../config/api.config.ts) - Configuration
- [error-handler.ts](./error-handler.ts) - Error utilities
- [api-helpers.ts](./api-helpers.ts) - Helper functions

---

**Updated**: October 31, 2025
