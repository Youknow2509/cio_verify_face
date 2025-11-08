#!/bin/bash

echo "Building service_device and push to docker hub"
docker buildx build --platform linux/amd64,linux/arm64 --push -t someone2509/cio_verify_face_service_device:latest .