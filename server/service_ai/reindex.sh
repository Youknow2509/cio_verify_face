#!/bin/bash
python -c "from app.services.face_service import FaceService; FaceService().reindex(force=True)"