# Face Attendance SaaS - Client Applications

This is a monorepo containing the 4 main client applications for the Face Attendance SaaS system:

1.  **Web Admin App** (`apps/web-admin`): For company administrators (HR/Admin).
2.  **Device App** (`apps/device-app`): For IoT attendance devices (used by employees on-site).
3.  **System Admin Interface** (`apps/system-admin`): For system operators (platform level).
4.  **Employee Portal** (`apps/employee-portal`): **NEW** - For employees to view their attendance, schedule, request leave, and self-check-in/out.

## Shared Packages

*   `packages/ui-components`: Reusable UI components.
*   `packages/utils`: Shared utility functions.
*   `packages/types`: Shared TypeScript interfaces.

## Getting Started

Use `pnpm` as the package manager.

```bash
# Install dependencies
pnpm install

# Start a specific app
pnpm run dev:employee-portal
```
