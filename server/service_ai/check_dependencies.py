"""
Check if all required dependencies are installed
"""
import sys

def check_dependencies():
    """Check if all required dependencies are installed"""
    missing = []
    
    # Core dependencies
    required_modules = {
        'fastapi': 'fastapi',
        'uvicorn': 'uvicorn',
        'pydantic': 'pydantic',
        'pydantic_settings': 'pydantic-settings',
        'insightface': 'insightface',
        'cv2': 'opencv-python',
        'numpy': 'numpy',
        'PIL': 'Pillow',
        'faiss': 'faiss-cpu',
        'onnxruntime': 'onnxruntime',
        'asyncpg': 'asyncpg',
        'sqlalchemy': 'sqlalchemy',
        'alembic': 'alembic',
        'httpx': 'httpx',
        'aiohttp': 'aiohttp',
        'prometheus_client': 'prometheus-client',
        'pythonjsonlogger': 'python-json-logger',
        'cryptography': 'cryptography',
        'multipart': 'python-multipart',
        'grpc': 'grpcio',
        'dotenv': 'python-dotenv',
        'yaml': 'PyYAML',
    }
    
    for module, package in required_modules.items():
        try:
            __import__(module)
        except ImportError:
            missing.append(package)
    
    if missing:
        print("❌ Missing required dependencies:")
        for package in missing:
            print(f"   - {package}")
        print("\n📦 Install dependencies using one of these methods:")
        print("\n   Option 1 (Modern - pyproject.toml):")
        print("   pip install -e .")
        print("\n   Option 2 (Traditional - requirements.txt):")
        print("   pip install -r requirements-cpu.txt")
        print("\n   For GPU support:")
        print("   pip install -e \".[gpu]\"")
        print("   # or")
        print("   pip install -r requirements-gpu.txt")
        return False
    else:
        print("✅ All required dependencies are installed!")
        return True

if __name__ == "__main__":
    success = check_dependencies()
    sys.exit(0 if success else 1)
