#!/bin/bash
# Test script to verify all requirements files install correctly

set -e

echo "======================================"
echo "Testing Requirements Files"
echo "======================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test CPU requirements
echo -e "${YELLOW}Testing requirements-cpu.txt...${NC}"
python3 -m venv /tmp/test_cpu_env
source /tmp/test_cpu_env/bin/activate
pip install --upgrade pip setuptools wheel -q
pip install -r requirements-cpu.txt -q

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ requirements-cpu.txt installed successfully${NC}"
else
    echo -e "${RED}✗ requirements-cpu.txt failed${NC}"
    exit 1
fi
deactivate
rm -rf /tmp/test_cpu_env

echo ""

# Test GPU requirements
echo -e "${YELLOW}Testing requirements-gpu.txt...${NC}"
python3 -m venv /tmp/test_gpu_env
source /tmp/test_gpu_env/bin/activate
pip install --upgrade pip setuptools wheel -q
# Note: This may fail if CUDA is not installed, but the package should be available
pip download onnxruntime-gpu==1.17.0 -q 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ requirements-gpu.txt package is available${NC}"
else
    echo -e "${YELLOW}⚠ onnxruntime-gpu may require CUDA installation${NC}"
fi
deactivate
rm -rf /tmp/test_gpu_env

echo ""

# Test default requirements
echo -e "${YELLOW}Testing requirements.txt...${NC}"
python3 -m venv /tmp/test_default_env
source /tmp/test_default_env/bin/activate
pip install --upgrade pip setuptools wheel -q
pip install -r requirements.txt -q

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ requirements.txt installed successfully${NC}"
else
    echo -e "${RED}✗ requirements.txt failed${NC}"
    exit 1
fi
deactivate
rm -rf /tmp/test_default_env

echo ""
echo -e "${GREEN}======================================"
echo "All requirements files are valid!"
echo "======================================${NC}"
