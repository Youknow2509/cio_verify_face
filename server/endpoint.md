# Service ai
- Put       /api/v1/face/profile/:profile_id/upload
- Delete    /api/v1/face/profile/:profile_id
- Get       /api/v1/face/profiles/:user_id
- Post      /api/v1/face/enroll/upload
- Post      /api/v1/face/verify/upload
- Post      /api/v1/face/cleanup/profiles

# Service analytic
- GET    /swagger/*any             
- GET    /health                   
- GET    /api/v1/reports/daily     
- GET    /api/v1/reports/summary   
- POST   /api/v1/reports/export    
- GET    /api/v1/reports/download/:filename 
- GET    /api/v1/attendance-records 
- GET    /api/v1/attendance-records/range 
- GET    /api/v1/attendance-records/employee/:employee_id 
- GET    /api/v1/attendance-records/user/:employee_id 
- GET    /api/v1/daily-summaries   
- GET    /api/v1/daily-summaries/user/:employee_id 
- GET    /api/v1/audit-logs        
- GET    /api/v1/audit-logs/range  
- POST   /api/v1/audit-logs        
- GET    /api/v1/face-enrollment-logs 
- GET    /api/v1/face-enrollment-logs/employee/:employee_id 
- GET    /api/v1/attendance-records-no-shift 
- GET    /api/v1/company/daily-attendance-status 
- GET    /api/v1/company/attendance-status/range 
- GET    /api/v1/company/monthly-summary 
- POST   /api/v1/company/export-daily-status 
- POST   /api/v1/company/export-monthly-summary 
- GET    /api/v1/employee/my-attendance-records 
- GET    /api/v1/employee/my-attendance-records/range 
- GET    /api/v1/employee/my-daily-summaries 
- GET    /api/v1/employee/my-daily-summary/:date 
- GET    /api/v1/employee/my-stats 
- GET    /api/v1/employee/my-daily-status 
- GET    /api/v1/employee/my-status/range 
- GET    /api/v1/employee/my-monthly-summary 
- POST   /api/v1/employee/export-daily-status 
- POST   /api/v1/employee/export-monthly-summary 

# Service attendance
- GET    /swagger/*any             
- POST   /api/v1/attendance/       
- POST   /api/v1/attendance/records 
- POST   /api/v1/attendance/records/summary/daily 
- POST   /api/v1/attendance/records/employee 
- POST   /api/v1/attendance/records/employee/summary/daily 

# Service auth
- GET    /swagger/*any             
- POST   /api/v1/auth/login        
- POST   /api/v1/auth/login/admin  
- POST   /api/v1/auth/refresh      
- POST   /api/v1/auth/logout       
- GET    /api/v1/auth/me           
- POST   /api/v1/auth/device       
- DELETE /api/v1/auth/device   
      
# Service device
- GET    /swagger/*any             
- GET    /api/v1/device            
- POST   /api/v1/device            
- GET    /api/v1/device/:device_id 
- GET    /api/v1/device/token/:device_id 
- POST   /api/v1/device/token/refresh/:device_id 
- PUT    /api/v1/device/:device_id 
- DELETE /api/v1/device/:device_id 
- POST   /api/v1/device/location   
- POST   /api/v1/device/name       
- POST   /api/v1/device/info       
- POST   /api/v1/device/status     

# Service identity
- GET       /api/v1/companies
- POST      /api/v1/companies
- GET       /api/v1/companies/:company_id
- PUT       /api/v1/companies/:company_id
- DELETE    /api/v1/companies/:company_id

- GET       /api/v1/users
- POST      /api/v1/users
- GET       /api/v1/users/:user_id
- PUT       /api/v1/users/:user_id
- DELETE    /api/v1/users/:user_id
- POST      /api/v1/users/:user_id/face-data
- GET       /api/v1/users/:user_id/face-data
- POST      /api/v1/users/:user_id/face-data/upload
- DELETE    /api/v1/users/:user_id/face-data/:fid
- PUT       /api/v1/users/:user_id/face-data/:fid/primary
- 
# Service notify

# Service profile update
- GET    /swagger/*any             
- GET    /health                   
- POST   /api/v1/profile-update/requests 
- GET    /api/v1/profile-update/requests/me 
- GET    /api/v1/profile-update/requests/pending 
- POST   /api/v1/profile-update/requests/:id/approve 
- POST   /api/v1/profile-update/requests/:id/reject 
- GET    /api/v1/profile-update/token/validate 
- POST   /api/v1/profile-update/face 
- POST   /api/v1/password/reset    

# Service signature

# Service workforce
- GET    /swagger/*any           
- GET    /api/v1/shift             
- POST   /api/v1/shift             
- GET    /api/v1/shift/:id         
- POST   /api/v1/shift/edit        
- DELETE /api/v1/shift/:id         
- POST   /api/v1/shift/status      
- POST   /api/v1/employee/shift    
- POST   /api/v1/employee/shift/edit/effective 
- POST   /api/v1/employee/shift/enable 
- POST   /api/v1/employee/shift/disable 
- POST   /api/v1/employee/shift/delete 
- POST   /api/v1/employee/shift/add 
- POST   /api/v1/employee/shift/add/list 
- POST   /api/v1/employee/shift/not_in 
- POST   /api/v1/employee/shift/in 

# Service ws delivery
- GET    /ws                       
- GET    /api/health               
- GET    /api/health/details       