#!/bin/bash
# Generate gRPC code from proto files

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Create output directory if it doesn't exist
mkdir -p app/grpc_generated

# Generate Python gRPC code using current Python
python -m grpc_tools.protoc \
    -I./protos \
    --python_out=./app/grpc_generated \
    --grpc_python_out=./app/grpc_generated \
    --pyi_out=./app/grpc_generated \
    ./protos/face_service.proto ./protos/auth.proto ./protos/attendance.proto

# Fix imports in generated files to use relative imports
# This fixes the "No module named 'face_service_pb2'" error
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS sed
    sed -i '' 's/^import face_service_pb2 as face__service__pb2/from . import face_service_pb2 as face__service__pb2/' app/grpc_generated/face_service_pb2_grpc.py
    sed -i '' 's/^import auth_pb2 as auth__pb2/from . import auth_pb2 as auth__pb2/' app/grpc_generated/auth_pb2_grpc.py
    sed -i '' 's/^import attendance_pb2 as attendance__pb2/from . import attendance_pb2 as attendance__pb2/' app/grpc_generated/attendance_pb2_grpc.py

else
    # Linux sed
    sed -i 's/^import face_service_pb2 as face__service__pb2/from . import face_service_pb2 as face__service__pb2/' app/grpc_generated/face_service_pb2_grpc.py
    sed -i 's/^import auth_pb2 as auth__pb2/from . import auth_pb2 as auth__pb2/' app/grpc_generated/auth_pb2_grpc.py
    sed -i 's/^import attendance_pb2 as attendance__pb2/from . import attendance_pb2 as attendance__pb2/' app/grpc_generated/attendance_pb2_grpc.py
fi

# Create __init__.py in generated directory
cat > app/grpc_generated/__init__.py << 'EOF'
"""
Generated gRPC code for face verification service.
This package contains protobuf and gRPC generated files.
"""
from .face_service_pb2 import *
from .face_service_pb2_grpc import *
from .auth_pb2 import *
from .auth_pb2_grpc import *
from .attendance_pb2 import *
from .attendance_pb2_grpc import *

__all__ = [
    'face_service_pb2',
    'face_service_pb2_grpc',
    'auth_pb2',
    'auth_pb2_grpc',
    'attendance_pb2',
    'attendance_pb2_grpc',
]
EOF

echo "gRPC code generated successfully!"
echo "Generated files in app/grpc_generated/:"
ls -1 app/grpc_generated/*.py
