.PHONY: help
.PHONY: docker_run, docker_stop
.ALL: help

# Project info
PROJECT_NAME=cio_verify_face

help:
	@echo "Available commands:"

	@echo "\nDocker Commands:"
	@echo "\t docker_run \t Running container for dev ..."
	@echo "\t docker_stop \t Stopping container for dev ..."

	@echo "  help - Show this help message"

# For docker commands:
docker_run:
	@echo "Docker running ..."
	docker compose -p $(PROJECT_NAME) -f docker-compose-dev.yml up -d

docker_stop:
	@echo "Docker stopping ..."
	docker compose -p $(PROJECT_NAME) -f docker-compose-dev.yml down