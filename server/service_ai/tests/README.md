# Tests directory

This directory contains integration tests for the Face Verification AI Service.

## Running Tests

### Prerequisites

Ensure ScyllaDB and MinIO services are running:

```bash
docker-compose up -d scylladb minio
```

### Install Test Dependencies

```bash
pip install pytest pytest-asyncio
```

### Run All Tests

```bash
pytest tests/ -v
```

### Run Specific Test File

```bash
pytest tests/test_integration.py -v
```

### Run Specific Test

```bash
pytest tests/test_integration.py::TestScyllaDBIntegration::test_save_verification_state -v
```

## Test Coverage

- **test_integration.py**: Integration tests for ScyllaDB and MinIO
  - ScyllaDB CRUD operations
  - MinIO upload/download operations
  - End-to-end workflows

## Notes

Tests will automatically skip if ScyllaDB or MinIO are not available, allowing the test suite to run in various environments.
