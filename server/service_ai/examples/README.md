# Face Verification Service Examples

This directory contains example scripts for using the Face Verification Service.

## Files

- `usage_example.py` - Complete Python client example with all API operations

## Running Examples

### Prerequisites

```bash
pip install requests
```

### Health Check Example

```bash
python usage_example.py
```

### Custom Usage

See the `FaceVerificationClient` class in `usage_example.py` for all available methods:

- `health_check()` - Check service status
- `enroll_face()` - Enroll a new face
- `verify_face()` - Verify a face (1:1 or 1:N)
- `get_user_profiles()` - Get all profiles for a user
- `update_profile()` - Update a face profile
- `delete_profile()` - Delete a profile
- `reindex()` - Rebuild FAISS index (admin)

## Example: Enroll and Verify

```python
from usage_example import FaceVerificationClient
from uuid import uuid4

client = FaceVerificationClient("http://localhost:8080")

# Enroll a face
user_id = uuid4()
result = client.enroll_face(
    user_id=user_id,
    image_path="path/to/face.jpg",
    make_primary=True
)
print(f"Enrollment: {result}")

# Verify a face (1:N search)
result = client.verify_face(
    image_path="path/to/verify.jpg",
    search_mode="1:N"
)
print(f"Verification: {result}")
```

## JavaScript Example

```javascript
const BASE_URL = "http://localhost:8080/api/v1/face";

// Health check
fetch(`${BASE_URL}/../health`)
  .then(r => r.json())
  .then(data => console.log(data));

// Enroll face
const enrollData = {
  user_id: "550e8400-e29b-41d4-a716-446655440000",
  image_base64: "...",
  make_primary: true
};

fetch(`${BASE_URL}/enroll`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(enrollData)
})
  .then(r => r.json())
  .then(data => console.log(data));
```

## cURL Examples

### Health Check
```bash
curl http://localhost:8080/health
```

### Enroll Face
```bash
curl -X POST http://localhost:8080/api/v1/face/enroll \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "image_base64": "...",
    "make_primary": true
  }'
```

### Verify Face
```bash
curl -X POST http://localhost:8080/api/v1/face/verify \
  -H "Content-Type: application/json" \
  -d '{
    "image_base64": "...",
    "search_mode": "1:N",
    "top_k": 5
  }'
```
