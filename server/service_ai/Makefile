.PHONY: help install install-editable install-dev install-gpu-editable install-all dev build up down logs clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Installation using requirements.txt
install: ## Install dependencies (CPU version using requirements.txt)
	pip install -r requirements.txt

install-cpu: ## Install dependencies for CPU mode (development/testing using requirements.txt)
	pip install -r requirements-cpu.txt

install-gpu: ## Install dependencies for GPU mode (production with CUDA using requirements.txt)
	pip install -r requirements-gpu.txt

# Installation using pyproject.toml (recommended)
install-editable: ## Install package in editable mode (CPU, recommended for development)
	pip install -e .

install-dev: ## Install with development tools (black, flake8, mypy, pytest)
	pip install -e ".[dev]"

install-gpu-editable: ## Install with GPU support in editable mode
	pip install -e ".[gpu]"

install-all: ## Install all optional dependencies (GPU + dev + test)
	pip install -e ".[all]"

check-deps: ## Check if all dependencies are installed
	python check_dependencies.py

generate-grpc: ## Generate gRPC code from proto files
	./generate_grpc.sh

run: ## Run production server
	python run_servers.py

dev: ## Run development server
	uvicorn app.main:app --reload --host 0.0.0.0 --port 8080

build: ## Build Docker image
	docker-compose build

up: ## Start services with Docker Compose
	docker-compose up -d

down: ## Stop services
	docker-compose down

logs: ## Show logs
	docker-compose logs -f service_ai

clean: ## Clean up generated files
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete
	find . -type f -name "*.pyo" -delete
	find . -type d -name "*.egg-info" -exec rm -rf {} + 2>/dev/null || true

test: ## Run tests
	pytest tests/ -v

lint: ## Run linter
	flake8 app/ --max-line-length=100 --exclude=__pycache__,grpc_generated

format: ## Format code with black and isort
	black app/ --line-length=100
	isort app/

type-check: ## Run mypy type checker
	mypy app/

check: lint type-check ## Run all checks
	@echo "All checks passed!"
