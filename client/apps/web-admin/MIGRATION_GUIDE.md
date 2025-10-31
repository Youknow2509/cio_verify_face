# 🔄 Migration Guide: Mock API → Real API

Hướng dẫn chuyển từ Mock API sang Real API Backend.

## Tổng Quan

### Mock API (Development)
- ✅ Không cần backend
- ✅ Phát triển frontend nhanh
- ✅ Dữ liệu fake

**Bật:** `VITE_ENABLE_MOCK_API=true` trong `.env`

### Real API (Integration)
- ✅ Kết nối backend thực
- ✅ Dữ liệu thực
- ✅ Xác thực người dùng

**Bật:** `VITE_ENABLE_MOCK_API=false` trong `.env`

## Step-by-Step Migration

### 1. Start Backend Services
```bash
cd server
docker-compose up -d
```

Verify:
```bash
curl http://localhost:8080/api/v1/ping
```

### 2. Update Environment Variables
```bash
# .env
VITE_API_BASE_URL=http://localhost:8080
VITE_ENABLE_MOCK_API=false
VITE_API_TIMEOUT=10000
```

### 3. Clear Browser Cache
```bash
# Clear localStorage
localStorage.clear()

# Clear sessionStorage
sessionStorage.clear()

# Restart dev server
npm run dev
```

### 4. Test Basic Flows
- ✅ Login page loads
- ✅ Login với credentials
- ✅ Dashboard loads
- ✅ API calls appear in Network tab

### 5. Fix Common Issues

#### Issue: 401 Unauthorized
```
Solution: 
- Check token in localStorage: auth_token
- Verify token format in browser DevTools
- Check backend auth service
```

#### Issue: 403 Forbidden
```
Solution:
- Check user role/permissions
- Verify tenant ID matches
- Check CORS settings
```

#### Issue: CORS Error
```
Solution:
- Backend must allow origin
- Check CORS configuration
- Verify credentials mode: 'include'
```

#### Issue: 500 Server Error
```
Solution:
- Check backend logs
- Verify database connection
- Check data validation
```

## API Endpoints Reference

### Authentication
```typescript
// src/services/api/auth.api.ts
loginAPI(email, password)
logoutAPI()
refreshTokenAPI()
getCurrentUserAPI()
activateDeviceAPI(deviceCode, deviceSecret)
```

### Employee Management
```typescript
// src/services/api/employees.api.ts
getEmployeesAPI(filter)
getEmployeeAPI(userId)
createEmployeeAPI(data)
updateEmployeeAPI(userId, data)
deleteEmployeeAPI(userId)
uploadFaceDataAPI(userId, file)
```

### Attendance
```typescript
// src/services/api/attendance.api.ts
checkInAPI(faceImage)
checkOutAPI(faceImage)
getAttendanceRecordsAPI(filter)
getMyAttendanceHistoryAPI(filter)
```

### Reports
```typescript
// src/services/api/reports.api.ts
getDailyReportAPI(date)
getSummaryReportAPI(startDate, endDate)
exportReportAPI(params)
```

Full list: [src/services/API_GUIDE.md](src/services/API_GUIDE.md)

## Testing Checklist

- [ ] Backend services running
- [ ] Environment variables updated
- [ ] Login works
- [ ] Can view employees list
- [ ] Can create new employee
- [ ] Can upload face data
- [ ] Check-in/check-out works
- [ ] Reports load correctly
- [ ] Filters & search work
- [ ] Export functionality works

## Debugging Tips

### Check Network Requests
1. Open DevTools: `F12`
2. Go to **Network** tab
3. Click **XHR** filter
4. Perform action and check request/response

### Check Browser Storage
1. DevTools → **Application** tab
2. **Local Storage** → Check `auth_token`
3. **Session Storage** → Check temporary data

### Check Backend Logs
```bash
# If using Docker
docker logs -f [container_id]

# If running locally
tail -f logs/app.log
```

### Enable Debug Mode
```typescript
// src/services/http.ts
const DEBUG = true; // Logs all API calls
```

## Rollback to Mock API

If issues occur:

```bash
# .env
VITE_ENABLE_MOCK_API=true
VITE_API_BASE_URL=http://localhost:3003
```

Then:
```bash
npm run dev
```

## Performance Considerations

- Mock API: Fast, instant responses
- Real API: Network delay (50-200ms typical)
- Consider adding loading indicators
- Implement proper error handling

## Security Notes

- ✅ Tokens stored in localStorage
- ✅ HTTPS required in production
- ✅ CORS properly configured
- ✅ Sensitive data not logged

---

**Need Help?** 
- Check logs in browser console
- Verify backend is running
- Check `.env` configuration
- Create issue with error message
