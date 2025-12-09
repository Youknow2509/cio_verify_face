# CIO Verify Face — Face Attendance SaaS

## Contact

-   **Mail**: *lytranvinh.work@gmail.com*
-   **Github**: *https://github.com/Youknow2509*

[Demo Video](https://youtu.be/y7Uxjtw-5sA)

## System Overview

-   Multi-tenant SaaS for face attendance with web admin portals and device apps.
-   Core capabilities: face-based check-in and check-out, company and employee directory, device provisioning, shift scheduling, analytics and exports, signature capture, real-time events.
-   Tech layout: API gateway (NGINX) in front of Go services and a Python AI service; observability with Prometheus, Grafana, Jaeger.


## Tech Stack
-   Backend: Go (Gin, ..), Python (FastAPI, OpenCV, dlib, face_recognition, ...), Node JS (Express, ...)
-   Environment: PostgreSQL, Redis, NGINX, Kafka, MinIO, Scylladb, Mivills, Grafana, Prometheus, Jaeger
-   Client: React, TypeScript, Redux, TailwindCSS, Electron
-   DevOps: Docker, Docker Compose, GitHub Actions, Helm, Kubernetes
-   Observability: Prometheus, Grafana, Jaeger
-   ...

## Architecture and Services

-   `service-auth`: Authentication, session refresh, device activation, token issuance.
-   `service-identity` and Organization: Company, user, employee lifecycle; face data enrollment; tenant isolation.
-   `service-device`: Device registry, provisioning, health status.
-   `service-attendance`: Receive attendance events, persist raw records.
-   `service-workforce`: Shifts, schedules, assignments and calendars.
-   `service-analytic`: Aggregations, daily and monthly summaries, report exports (Excel or PDF).
-   `service-signature`: Digital signature upload and retrieval.
-   `service-ai`: Face recognition pipeline and related ML tasks (Python).
-   `service-websocket`: Real-time event broadcasting to clients via WebSocket.
-   `serice-analytics`: System usage metrics, monitoring, and alerting.
-   `service-notification`: Email and SMS sending.
-   `service-profile-update`: Employee profile self-service.
-   ...

## Client Applications

-   `apps/device-app`: On-site device UI for face capture and check in or out.
-   `apps/employee-portal`: Employee self-service history and profile view.
-   `apps/web-admin`: Company admin console.
-   `apps/system-admin`: Provider-level administration(pending).

## End-to-End Flow (high level)

1. Activate device through Auth to obtain device token.
2. Enroll employee face data through Identity; optional signature enrollment through Signature service.
3. Assign shifts and schedules through Workforce.
4. Device captures face, AI matches identity, Attendance stores raw records with metadata.
5. Analytics builds daily and monthly summaries and supports export.
6. WebSocket events (attendance_result, device_status, admin_alert) keep clients in sync in real time.

## Data Management and Isolation

-   Tenant isolation at company level for all employee, device, attendance, and report data.
-   Face data and signatures are stored per employee with controlled access paths.
-   Health checks and common utility endpoints exposed for monitoring and automation.

## Observability and Operations

-   Start observability stack: `cd server && docker-compose -f docker-compose-observability.yml up -d`.
-   Dashboards: Grafana `http://localhost:3000` (admin/admin), Prometheus `http://localhost:9091`, Jaeger `http://localhost:16686`.
-   Services expose metrics at `/metrics` (typical ports 9090 or 8000) and traces to Jaeger OTLP `http://jaeger:4318/v1/traces`. See `doc/OBSERVABILITY.md` for queries and alerts.

## Documentation Map

-   Functional and non-functional requirements: `doc/doc.md` (SRS, Vietnamese).
-   Endpoint catalog: `doc/endpoints.md` (REST and WebSocket overview).
-   UI pages and flows: `doc/client_ui_pages.md`, `doc/ui_workflows.md`.
-   Diagrams: `doc/sequence_diagrams/` and `doc/activity_diagrams/`.
-   Service responsibilities: `doc/service_usage.md`.
