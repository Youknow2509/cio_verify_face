# Face Attendance SaaS - Activity Diagrams

This directory contains Mermaid activity diagrams for all API endpoints and events in the Face Attendance SaaS system. Each diagram shows the process flow and decision points for handling requests in a multi-tenant architecture.

## Directory Structure

```
activity_diagrams/
├── auth/                    # Authentication Service (5 activities)
├── identity_org/           # Identity & Organization Service (4 activities)  
├── device_mgmt/            # Device Management Service (1 activity)
├── workforce/              # Workforce Service (1 activity)
├── attendance/             # Attendance Service (3 activities)
├── analytics_reporting/    # Analytics & Reporting Service (2 activities)
├── signature/              # Signature Upload Service (1 activity)
├── system_admin/           # System Admin APIs (1 activity)
├── common/                 # Common/Utility APIs (2 activities)
└── websocket/              # WebSocket Events (2 activities)
```

## Activity Diagrams Overview

### Auth Service
- ✅ `auth/login.mmd` - Login process flow with authentication and validation
- ✅ `auth/logout.mmd` - Logout process with session invalidation  
- ✅ `auth/refresh.mmd` - Token refresh workflow with validation checks
- ✅ `auth/device_activate.mmd` - Device activation process with certificate generation
- ✅ `auth/me.mmd` - User information retrieval workflow

### Identity & Organization Service
- ✅ `identity_org/list_companies.mmd` - Company listing with filtering and pagination
- ✅ `identity_org/create_company.mmd` - Company creation with validation and setup
- ✅ `identity_org/create_employee.mmd` - Employee creation with validation workflow
- ✅ `identity_org/upload_face_data.mmd` - Face data upload with detection and validation

### Device Management Service  
- ✅ `device_mgmt/create_device.mmd` - Device creation and registration workflow

### Workforce Service
- ✅ `workforce/create_shift.mmd` - Shift creation with conflict detection

### Attendance Service
- ✅ `attendance/check_in.mmd` - Check-in process with face verification
- ✅ `attendance/check_out.mmd` - Check-out process with duration calculation
- ✅ `attendance/employee_history.mmd` - Attendance history retrieval workflow

### Analytics & Reporting Service
- ✅ `analytics_reporting/daily_report.mmd` - Daily report generation workflow
- ✅ `analytics_reporting/export_report.mmd` - Report export process with background jobs

### Signature Upload Service
- ✅ `signature/upload_signature.mmd` - Signature image upload and processing

### System Admin APIs
- ✅ `system_admin/lock_company.mmd` - Company locking administrative process

### Common/Utility APIs
- ✅ `common/ping.mmd` - Health check process with system validation
- ✅ `common/time.mmd` - Server time retrieval workflow

### WebSocket Events
- ✅ `websocket/attendance_result.mmd` - Real-time attendance result broadcasting
- ✅ `websocket/device_status.mmd` - Device status change notifications

## Activity Diagram Elements

Each activity diagram includes:

### Flow Control Elements
- **Start/End nodes**: Entry and exit points of the process
- **Decision diamonds**: Conditional branching points with Yes/No outcomes
- **Process rectangles**: Individual processing steps
- **Parallel flows**: Concurrent activities where applicable

### Color Coding
- **Light Blue** (`#e1f5fe`): Start nodes and input points
- **Light Red** (`#ffebee`): Error conditions and failure paths
- **Light Green** (`#e8f5e8`): Success outcomes and completion points
- **Light Orange** (`#fff3e0`): Warning conditions and intermediate states
- **Light Purple** (`#f3e5f5`): Skip/bypass conditions
- **Light Gray** (`#f5f5f5`): End nodes and process completion

### Key Process Patterns

#### Authentication & Authorization
- Token validation flows
- Permission checking at multiple levels
- Role-based access control decisions

#### Data Validation & Processing
- Input validation with detailed error handling
- Multi-step processing workflows
- Cache-first data retrieval patterns

#### Error Handling
- Comprehensive error condition coverage
- Appropriate HTTP status code returns
- Logging and monitoring integration points

#### Performance Optimizations
- Caching decision points
- Rate limiting checks
- Background job queuing for heavy operations

## Usage

To view these activity diagrams:

1. **GitHub**: Diagrams render automatically in GitHub's Markdown preview
2. **Mermaid Live Editor**: Copy content to https://mermaid.live/
3. **VS Code**: Use Mermaid Preview extension
4. **Documentation sites**: Most support Mermaid rendering (GitBook, Notion, etc.)

## Relationship to Sequence Diagrams

These activity diagrams complement the sequence diagrams in `../sequence_diagrams/` by:

- **Sequence diagrams**: Show *interactions* between different services and components
- **Activity diagrams**: Show the *process flow* and decision logic within a single service

Together, they provide comprehensive documentation of both the system architecture (sequence) and the business logic flow (activity).

## Architecture Context

These activity diagrams align with the system architecture and show how individual services handle requests internally, including:

- **Decision points**: Where business rules are applied
- **Error handling**: How failures are managed and reported  
- **Performance considerations**: Caching, rate limiting, and optimization points
- **Security checkpoints**: Authentication and authorization validations
- **Data flow**: How information moves through processing steps

The diagrams demonstrate the robust error handling, security controls, and performance optimizations built into each service to support enterprise-scale operations for millions of companies.