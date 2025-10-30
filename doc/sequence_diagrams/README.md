# Face Attendance SaaS - Sequence Diagrams

This directory contains Mermaid sequence diagrams for all API endpoints in the Face Attendance SaaS system. Each diagram shows the optimized flow for handling millions of companies in a multi-tenant architecture.

## Directory Structure

```
sequence_diagrams/
├── auth/                    # Authentication Service (5 endpoints)
├── identity_org/           # Identity & Organization Service (12 endpoints)  
├── device_mgmt/            # Device Management Service (5 endpoints)
├── workforce/              # Workforce Service (7 endpoints)
├── attendance/             # Attendance Service (5 endpoints)
├── analytics_reporting/    # Analytics & Reporting Service (3 endpoints)
├── signature/              # Signature Upload Service (3 endpoints)
├── system_admin/           # System Admin APIs (4 endpoints)
├── common/                 # Common/Utility APIs (2 endpoints)
└── websocket/              # WebSocket Events (3 events)
```

## Completed Diagrams

### Auth Service
- ✅ `auth/login.mmd` - POST /api/v1/auth/login
- ✅ `auth/logout.mmd` - POST /api/v1/auth/logout  
- ✅ `auth/refresh.mmd` - POST /api/v1/auth/refresh
- ✅ `auth/device_activate.mmd` - POST /api/v1/auth/device/activate
- ✅ `auth/me.mmd` - GET /api/v1/auth/me

### Identity & Organization Service
- ✅ `identity_org/list_companies.mmd` - GET /api/v1/companies
- ✅ `identity_org/create_company.mmd` - POST /api/v1/companies
- ✅ `identity_org/upload_face_data.mmd` - POST /api/v1/users/{user_id}/face-data

### Device Management Service  
- ✅ `device_mgmt/create_device.mmd` - POST /api/v1/devices

### Workforce Service
- ✅ `workforce/create_shift.mmd` - POST /api/v1/shifts

### Attendance Service
- ✅ `attendance/check_in.mmd` - POST /api/v1/attendance/check-in
- ✅ `attendance/check_out.mmd` - POST /api/v1/attendance/check-out

### Analytics & Reporting Service
- ✅ `analytics_reporting/daily_report.mmd` - GET /api/v1/reports/daily

### Signature Upload Service
- ✅ `signature/upload_signature.mmd` - POST /api/v1/signatures

### System Admin APIs
- ✅ `system_admin/lock_company.mmd` - PUT /api/v1/admin/companies/{id}/lock

### Common/Utility APIs
- ✅ `common/ping.mmd` - GET /api/v1/ping

### WebSocket Events
- ✅ `websocket/attendance_result.mmd` - attendance_result event

## Key Optimizations for Scale

All sequence diagrams incorporate enterprise-scale optimizations:

### Performance Optimizations
- **Multi-level caching**: Redis clusters with intelligent TTL management
- **Database sharding**: Company-based partitioning for horizontal scaling
- **Connection pooling**: Efficient database connection management
- **Async processing**: Kafka-based event-driven architecture
- **CDN integration**: Static content delivery optimization

### Security Optimizations  
- **JWT token management**: Short-lived access tokens with refresh rotation
- **Rate limiting**: Per-tenant and per-endpoint rate controls
- **Data isolation**: Company-scoped data access with strict validation
- **Encryption**: End-to-end encryption for sensitive data (face embeddings, signatures)

### Scalability Optimizations
- **Horizontal scaling**: Load balancer distribution across service instances
- **Microservices isolation**: Independent scaling of different service components
- **Face verification pipeline**: GPU-accelerated distributed face matching
- **Time-series optimization**: Efficient attendance record storage and querying
- **Real-time capabilities**: WebSocket-based instant notifications

### Monitoring & Reliability
- **Health checks**: Comprehensive system health monitoring
- **Audit logging**: Complete trail of all administrative actions
- **Circuit breakers**: Fault tolerance for service dependencies
- **Graceful degradation**: Fallback mechanisms for service failures

## Usage

To view these diagrams:

1. **GitHub**: Diagrams will render automatically in GitHub's Markdown preview
2. **Mermaid Live Editor**: Copy content to https://mermaid.live/
3. **VS Code**: Use Mermaid Preview extension
4. **Documentation sites**: Most support Mermaid rendering (GitBook, Notion, etc.)

## Architecture Context

These sequence diagrams align with the system architecture defined in `../architecture.mmd`, showing how the microservices, databases, caches, and message queues work together to provide a scalable face attendance solution for millions of companies.

Each diagram includes:
- **Participants**: All involved services and infrastructure components
- **Flow**: Step-by-step interaction sequence
- **Error handling**: Alternative flows for error scenarios  
- **Performance notes**: Scale-specific optimizations
- **Security considerations**: Authentication and authorization flows

The diagrams demonstrate how the system can efficiently handle the load requirements of a multi-tenant SaaS platform serving enterprise customers worldwide.