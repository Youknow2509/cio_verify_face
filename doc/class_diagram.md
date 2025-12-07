# Class Diagram - Hệ thống Chấm công Khuôn mặt SaaS

## Mô tả

Biểu đồ lớp (Class Diagram) cho Chapter 2, thể hiện cấu trúc các lớp thực thể (Entity Classes) và các lớp Service trong hệ thống.

---

## Figure 2.1 – Domain Entity Class Diagram (Mermaid)

```mermaid
classDiagram
    direction TB
    
    %% ==================== ENTITY CLASSES ====================
    
    class Company {
        -UUID companyId
        -String name
        -String address
        -String phone
        -CompanySettings settings
        -DateTime createdAt
        -DateTime updatedAt
        +getId() UUID
        +getName() String
        +getSettings() CompanySettings
        +updateInfo(name, address) void
    }
    
    class CompanySecret {
        -UUID secretId
        -UUID companyId
        -String salt
        -String secretHash
        +validate(secret) Boolean
        +regenerate() String
    }
    
    class User {
        -UUID userId
        -String email
        -String passwordHash
        -String salt
        -String fullName
        -UserRole role
        -Boolean isActive
        -DateTime createdAt
        +getId() UUID
        +getEmail() String
        +getRole() UserRole
        +validatePassword(password) Boolean
        +changePassword(newPassword) void
        +deactivate() void
    }
    
    class Employee {
        -UUID employeeId
        -UUID companyId
        -UUID userId
        -String employeeCode
        -String department
        -Boolean isActive
        -DateTime hireDate
        +getId() UUID
        +getCode() String
        +getCompanyId() UUID
        +getFaceData() FaceData
        +assignToShift(shift) void
    }
    
    class FaceData {
        -UUID faceId
        -UUID employeeId
        -Blob faceEmbedding
        -String imageUrl
        -DateTime createdAt
        +getId() UUID
        +getEmbedding() Blob
        +getImageUrl() String
        +updateEmbedding(embedding) void
    }
    
    class Device {
        -UUID deviceId
        -UUID companyId
        -UUID locationId
        -String name
        -String address
        -DeviceStatus status
        -String token
        -DateTime lastHeartbeat
        +getId() UUID
        +getName() String
        +getStatus() DeviceStatus
        +isOnline() Boolean
        +updateStatus(status) void
        +generateToken() String
    }
    
    class WorkShift {
        -UUID shiftId
        -UUID companyId
        -String name
        -Time startTime
        -Time endTime
        -Integer lateThresholdMinutes
        -Integer earlyLeaveMinutes
        +getId() UUID
        +getName() String
        +getStartTime() Time
        +getEndTime() Time
        +isLate(checkInTime) Boolean
        +isEarlyLeave(checkOutTime) Boolean
    }
    
    class AttendanceRecord {
        -UUID recordId
        -UUID employeeId
        -UUID deviceId
        -DateTime timestamp
        -RecordType recordType
        -Float verificationScore
        -String faceImageUrl
        -JSON metadata
        +getId() UUID
        +getEmployeeId() UUID
        +getTimestamp() DateTime
        +getType() RecordType
        +getScore() Float
    }
    
    class DailyAttendanceSummary {
        -UUID summaryId
        -UUID employeeId
        -UUID shiftId
        -Date workDate
        -DateTime checkInTime
        -DateTime checkOutTime
        -Integer totalMinutes
        -Integer overtimeMinutes
        -AttendanceStatus status
        +getId() UUID
        +getWorkDate() Date
        +getStatus() AttendanceStatus
        +calculateTotalHours() Float
        +calculateOvertime() Integer
    }
    
    %% ==================== ENUMS ====================
    
    class UserRole {
        <<enumeration>>
        SYSTEM_ADMIN
        COMPANY_ADMIN
        EMPLOYEE
    }
    
    class DeviceStatus {
        <<enumeration>>
        ONLINE
        OFFLINE
        MAINTENANCE
        ERROR
    }
    
    class RecordType {
        <<enumeration>>
        CHECK_IN
        CHECK_OUT
    }
    
    class AttendanceStatus {
        <<enumeration>>
        PRESENT
        LATE
        EARLY_LEAVE
        ABSENT
    }
    
    %% ==================== RELATIONSHIPS ====================
    
    Company "1" --o "*" Employee : has
    Company "1" --o "*" Device : owns
    Company "1" --o "*" WorkShift : defines
    Company "1" --o "1..*" CompanySecret : has
    
    User "1" --o "0..1" Employee : may be
    User "*" --> "1" UserRole : has
    
    Employee "1" --o "1" FaceData : has
    Employee "1" --o "*" AttendanceRecord : generates
    Employee "1" --o "*" DailyAttendanceSummary : has
    
    Device "1" --o "*" AttendanceRecord : captures
    Device "*" --> "1" DeviceStatus : has
    Device "1" --> "1" CompanySecret : uses
    
    WorkShift "1" --o "*" DailyAttendanceSummary : applies to
    
    AttendanceRecord "*" --> "1" RecordType : has
    DailyAttendanceSummary "*" --> "1" AttendanceStatus : has
```

---

## Figure 2.2 – Service Layer Class Diagram (Mermaid)

```mermaid
classDiagram
    direction TB
    
    %% ==================== SERVICE INTERFACES ====================
    
    class IAuthService {
        <<interface>>
        +login(email, password) AuthResponse
        +logout(token) void
        +refreshToken(refreshToken) AuthResponse
        +activateDevice(activationCode) DeviceToken
        +getCurrentUser(token) User
    }
    
    class IIdentityService {
        <<interface>>
        +listCompanies(filter) List~Company~
        +createCompany(data) Company
        +getCompany(id) Company
        +updateCompany(id, data) Company
        +deleteCompany(id) void
        +listUsers(companyId) List~User~
        +createUser(data) User
        +updateUser(id, data) User
        +deleteUser(id) void
    }
    
    class IDeviceService {
        <<interface>>
        +listDevices(companyId) List~Device~
        +createDevice(data) Device
        +getDevice(id) Device
        +updateDevice(id, data) Device
        +deleteDevice(id) void
        +updateStatus(id, status) void
    }
    
    class IWorkforceService {
        <<interface>>
        +listShifts(companyId) List~WorkShift~
        +createShift(data) WorkShift
        +updateShift(id, data) WorkShift
        +deleteShift(id) void
        +assignEmployeeToShift(employeeId, shiftId) void
    }
    
    class IAttendanceService {
        <<interface>>
        +checkIn(deviceId, faceImage) AttendanceResult
        +checkOut(deviceId, faceImage) AttendanceResult
        +getRecords(filter) List~AttendanceRecord~
        +getRecord(id) AttendanceRecord
        +getEmployeeHistory(employeeId) List~AttendanceRecord~
    }
    
    class IFaceVerificationService {
        <<interface>>
        +detectFace(image) FaceDetectionResult
        +extractEmbedding(alignedFace) Embedding
        +matchFace(embedding, companyId) MatchResult
        +checkLiveness(image) LivenessResult
    }
    
    class IAnalyticsService {
        <<interface>>
        +getDailyReport(companyId, date) DailyReport
        +getMonthlySummary(companyId, month) MonthlySummary
        +exportReport(type, params) ReportFile
    }
    
    class ISignatureService {
        <<interface>>
        +uploadSignature(userId, file) Signature
        +getSignatures(userId) List~Signature~
        +deleteSignature(id) void
    }
    
    class INotificationService {
        <<interface>>
        +pushAttendanceResult(deviceId, result) void
        +pushDeviceStatus(deviceId, status) void
        +pushAdminAlert(companyId, alert) void
    }
    
    %% ==================== SERVICE IMPLEMENTATIONS ====================
    
    class AuthService {
        -UserRepository userRepo
        -TokenService tokenService
        -PasswordEncoder encoder
        +login(email, password) AuthResponse
        +logout(token) void
        +refreshToken(refreshToken) AuthResponse
    }
    
    class IdentityService {
        -CompanyRepository companyRepo
        -UserRepository userRepo
        -EmployeeRepository employeeRepo
        -FaceDataRepository faceRepo
    }
    
    class AttendanceService {
        -AttendanceRepository attendanceRepo
        -FaceVerificationService faceService
        -WorkforceService workforceService
        -EventPublisher eventPublisher
        +checkIn(deviceId, faceImage) AttendanceResult
        +checkOut(deviceId, faceImage) AttendanceResult
    }
    
    class FaceVerificationService {
        -FaceDetector detector
        -FaceAligner aligner
        -EmbeddingExtractor extractor
        -VectorMatcher matcher
        -LivenessChecker livenessChecker
    }
    
    class AnalyticsService {
        -SummaryRepository summaryRepo
        -ReportGenerator reportGenerator
        -EventConsumer eventConsumer
    }
    
    %% ==================== RELATIONSHIPS ====================
    
    AuthService ..|> IAuthService
    IdentityService ..|> IIdentityService
    AttendanceService ..|> IAttendanceService
    FaceVerificationService ..|> IFaceVerificationService
    AnalyticsService ..|> IAnalyticsService
    
    AttendanceService --> IFaceVerificationService : uses
    AttendanceService --> IWorkforceService : uses
    AttendanceService --> INotificationService : publishes to
    
    AnalyticsService --> INotificationService : subscribes from
```

---

## Figure 2.3 – Face Verification Pipeline Classes (Mermaid)

```mermaid
classDiagram
    direction LR
    
    class FaceVerificationPipeline {
        -FaceDetector detector
        -FaceAligner aligner
        -QualityChecker qualityChecker
        -LivenessDetector livenessDetector
        -EmbeddingExtractor extractor
        -VectorMatcher matcher
        +process(image) VerificationResult
    }
    
    class FaceDetector {
        <<interface>>
        +detect(image) List~BoundingBox~
        +detectWithLandmarks(image) FaceDetectionResult
    }
    
    class RetinaFaceDetector {
        -modelPath String
        -confidenceThreshold Float
        +detect(image) List~BoundingBox~
        +detectWithLandmarks(image) FaceDetectionResult
    }
    
    class FaceAligner {
        <<interface>>
        +align(image, landmarks) AlignedFace
    }
    
    class AffineFaceAligner {
        -targetSize Size
        +align(image, landmarks) AlignedFace
    }
    
    class QualityChecker {
        <<interface>>
        +checkQuality(face) QualityScore
        +isAcceptable(score) Boolean
    }
    
    class LivenessDetector {
        <<interface>>
        +checkLiveness(image) LivenessResult
        +isReal(result) Boolean
    }
    
    class EmbeddingExtractor {
        <<interface>>
        +extract(alignedFace) Embedding
    }
    
    class ArcFaceExtractor {
        -modelPath String
        -embeddingDim Integer
        +extract(alignedFace) Embedding
    }
    
    class VectorMatcher {
        <<interface>>
        +match(embedding, gallery) MatchResult
        +calculateSimilarity(v1, v2) Float
    }
    
    class CosineMatcher {
        -threshold Float
        +match(embedding, gallery) MatchResult
        +calculateSimilarity(v1, v2) Float
    }
    
    %% Value Objects
    class BoundingBox {
        -Float x
        -Float y
        -Float width
        -Float height
        -Float confidence
    }
    
    class Embedding {
        -Float[] vector
        -Integer dimension
        +cosineSimilarity(other) Float
    }
    
    class VerificationResult {
        -UUID employeeId
        -Float confidence
        -Boolean isAccepted
        -String rejectReason
    }
    
    %% Relationships
    FaceVerificationPipeline --> FaceDetector
    FaceVerificationPipeline --> FaceAligner
    FaceVerificationPipeline --> QualityChecker
    FaceVerificationPipeline --> LivenessDetector
    FaceVerificationPipeline --> EmbeddingExtractor
    FaceVerificationPipeline --> VectorMatcher
    
    RetinaFaceDetector ..|> FaceDetector
    AffineFaceAligner ..|> FaceAligner
    ArcFaceExtractor ..|> EmbeddingExtractor
    CosineMatcher ..|> VectorMatcher
    
    FaceDetector ..> BoundingBox : produces
    EmbeddingExtractor ..> Embedding : produces
    VectorMatcher ..> VerificationResult : produces
```

---

## Figure 2.4 – Repository Layer Class Diagram (Mermaid)

```mermaid
classDiagram
    direction TB
    
    class IRepository~T~ {
        <<interface>>
        +findById(id) T
        +findAll() List~T~
        +save(entity) T
        +update(entity) T
        +delete(id) void
    }
    
    class ICompanyRepository {
        <<interface>>
        +findById(id) Company
        +findByName(name) Company
        +findAllActive() List~Company~
    }
    
    class IUserRepository {
        <<interface>>
        +findById(id) User
        +findByEmail(email) User
        +findByCompanyId(companyId) List~User~
    }
    
    class IEmployeeRepository {
        <<interface>>
        +findById(id) Employee
        +findByCode(code, companyId) Employee
        +findByCompanyId(companyId) List~Employee~
        +findWithFaceData(id) Employee
    }
    
    class IDeviceRepository {
        <<interface>>
        +findById(id) Device
        +findByCompanyId(companyId) List~Device~
        +findByToken(token) Device
        +findOnlineDevices(companyId) List~Device~
    }
    
    class IAttendanceRepository {
        <<interface>>
        +findById(id) AttendanceRecord
        +findByEmployee(employeeId, dateRange) List~AttendanceRecord~
        +findByDevice(deviceId, dateRange) List~AttendanceRecord~
        +findByCompanyAndDate(companyId, date) List~AttendanceRecord~
    }
    
    class ISummaryRepository {
        <<interface>>
        +findByEmployeeAndDate(employeeId, date) DailyAttendanceSummary
        +findByCompanyAndMonth(companyId, month) List~DailyAttendanceSummary~
        +upsert(summary) DailyAttendanceSummary
    }
    
    class IFaceDataRepository {
        <<interface>>
        +findByEmployeeId(employeeId) FaceData
        +findByCompanyId(companyId) List~FaceData~
        +saveEmbedding(employeeId, embedding) FaceData
    }
    
    %% PostgreSQL Implementations
    class PostgresCompanyRepository {
        -dataSource DataSource
    }
    
    class PostgresUserRepository {
        -dataSource DataSource
    }
    
    class PostgresEmployeeRepository {
        -dataSource DataSource
    }
    
    %% ScyllaDB Implementation
    class ScyllaAttendanceRepository {
        -session ScyllaSession
        +findByCompanyAndDate(companyId, date) List~AttendanceRecord~
    }
    
    %% Relationships
    ICompanyRepository --|> IRepository
    IUserRepository --|> IRepository
    IEmployeeRepository --|> IRepository
    IDeviceRepository --|> IRepository
    IAttendanceRepository --|> IRepository
    
    PostgresCompanyRepository ..|> ICompanyRepository
    PostgresUserRepository ..|> IUserRepository
    PostgresEmployeeRepository ..|> IEmployeeRepository
    ScyllaAttendanceRepository ..|> IAttendanceRepository
```

---

## Tổng kết

Các Class Diagram trên thể hiện:

1. **Figure 2.1 - Domain Entity Classes**: 9 lớp thực thể chính (Company, User, Employee, FaceData, Device, WorkShift, AttendanceRecord, DailyAttendanceSummary) với các thuộc tính và phương thức.

2. **Figure 2.2 - Service Layer**: Các interface và implementation của 8 service chính trong hệ thống (Auth, Identity, Device, Workforce, Attendance, FaceVerification, Analytics, Notification).

3. **Figure 2.3 - Face Verification Pipeline**: Chi tiết các lớp trong pipeline xác thực khuôn mặt (Detector, Aligner, QualityChecker, EmbeddingExtractor, Matcher).

4. **Figure 2.4 - Repository Layer**: Các repository interface theo pattern Repository cho việc truy cập dữ liệu (PostgreSQL, ScyllaDB).

**Render diagrams:**
- VS Code: Extension "Markdown Preview Mermaid Support"
- Online: [Mermaid Live Editor](https://mermaid.live)
- GitHub: Tự động render trong Markdown preview
