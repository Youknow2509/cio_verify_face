.PHONY: help
.PHONY: goose_create goose_up goose_down goose_status goose_fix goose_redo goose_reset goose_clean goose_up_by_one
.PHONY: docker_compose_up docker_compose_down docker_builds docker_build_service docker_run_resource docker_nginx_dev

.ALL: help

# Project Name
PROJECT_NAME ?= cio_verify_face

# List Variables Goose
GOOSE=goose
GOOSE_NAME ?= cio_verify_face
GOOSE_DRIVER ?= postgres
GOOSE_DB_DSN = "postgresql://postgres:root1234@localhost:5433/cio_verify_face"
GOOSE_MIGRATION_DSN ?= database/migration/sql/
GOOSE_PATH_SCHEMA ?= database/migration/sql/
GOOSE_PATH_QUERIES = database/queries/

# List Variables Golang-migrate
MIGRATE=migrate
MIGRATE_PATH_SCHEMA ?= database/migration/cql/
MIGRATE_SOURCE ?= file://${MIGRATE_PATH_SCHEMA}
MIGRATE_DRIVER ?= cassandra
# cassandra://host:port/keyspace?param1=value&param2=value2
MIGRATE_URL ?= localhost:9042/cio_verify_face?username=cassandra&password=root1234&sslmode=disable
MIGRATE_DB_DSN ?= "${MIGRATE_DRIVER}://${MIGRATE_URL}"

# docker
DOCKER_BAKE_FILE ?= docker-bake.hcl

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

docker_run_observabilities: ## Docker Run Observabilities Containers
	@echo "Running observabilities containers"
	docker compose -p $(PROJECT_NAME) -f docker-compose-observability.yml up -d

docker_stop_observabilities: ## Docker Stop Observabilities Containers
	@echo "Stopping observabilities containers"
	docker compose -p $(PROJECT_NAME) -f docker-compose-observability.yml down

docker_run_resource: ## Docker Run Resource Containers
	@echo "Running resource containers"
	cd environment/docker && docker compose -p $(PROJECT_NAME) -f docker-compose-dev.yml up -d

docker_stop_resource: ## Docker Stop Resource Containers
	@echo "Stopping resource containers"
	cd environment/docker && docker compose -p $(PROJECT_NAME) -f docker-compose-dev.yml down

docker_build_service: ## Docker build a specific service and push to registry
	@echo "Building Docker image for service $(SERVICE) and pushing to registry"
	docker buildx bake --file ${DOCKER_BAKE_FILE} $(SERVICE) --push

docker_builds: ## Docker Build All Services
	@echo "Building Docker images for all services and pushing to registry"
	docker buildx bake --file ${DOCKER_BAKE_FILE} --push
	
docker_stop_services: ## Docker stop all services
	@echo "Stopping all Docker services for project $(PROJECT_NAME)"
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) down

docker_run_services: ## Docker run all services
	@echo "Starting all Docker services for project $(PROJECT_NAME)"
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) up -d

cmd_create_keyspace: ## Create keyspace scylladb
	@echo "Creating keyspace"
	@echo "-- Create keyspace with proper replication"
	@echo "CREATE KEYSPACE IF NOT EXISTS face_attendance WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 3} AND durable_writes = true;"
	@echo "USE face_attendance;"

migrate_force: ## Golang-migrate - Force a migration to a specific version
	@echo "Forcing migration to version 1"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} force ${VERSION}
	@echo "${GREEN_COLOR_BG}Migration forced to version ${VERSION}"

migrate_down: ## Golang-migrate - Down all migrations
	@echo "Rolling back the most recent migration"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} down
	@echo "${GREEN_COLOR_BG}Migration rolled back"

migrate_drop: ## Golang-migrate - Drop all migrations
	@echo "Dropping all migrations"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} drop
	@echo "${GREEN_COLOR_BG}All migrations dropped"

migrate_up: ## Golang-migrate - Run all up migrations
	@echo "Running all up migrations"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} up
	@echo "${GREEN_COLOR_BG}Migrations completed"

migrate_create: ## Golang-migrate - Create a new migration
	@echo "Creating a new migration"
	${MIGRATE} create -ext cql -dir ${MIGRATE_PATH_SCHEMA} -seq $(NAME)
	@echo "${GREEN_COLOR_BG}Migration created"

goose_create: ## Goosee - Create a new migration
	@echo "Creating a new migration"
	${GOOSE} -dir ${GOOSE_PATH_SCHEMA} create $(NAME) sql
	@echo "${GREEN_COLOR_BG}Migration created"

goose_up_by_one: ## Goose - Migrate the DB up by 1 
	@echo "Migrate the DB up by 1"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) up-by-one
	@echo "${GREEN_COLOR_BG}Migrations completed"

goose_up: ## Goose - Run all available migrations
	@echo "Running all available migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) up
	@echo "${GREEN_COLOR_BG}Migrations completed"

goose_down: ## Goose - Rollback the most recent migration
	@echo "Rolling back the most recent migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) down
	@echo "${GREEN_COLOR_BG}Migration rolled back"

goose_status: ## Goose - Show the status of all migrations
	@echo "Showing the status of all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) status
	@echo "${GREEN_COLOR_BG}Migration status displayed"

goose_fix: ## Goose - Fix the last migration
	@echo "Fixing the last migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) fix
	@echo "${GREEN_COLOR_BG}Migration fixed"

goose_redo: ## Goose - Rollback and re-run the most recent migration
	@echo "Rolling back and re-running the most recent migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) redo
	@echo "${GREEN_COLOR_BG}Migration redone"

goose_reset: ## Goose - Rollback all migrations
	@echo "Rolling back all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) reset
	@echo "${GREEN_COLOR_BG}Migrations reset"

goose_clean: ## Goose - Remove all migrations
	@echo "Removing all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) clean
	@echo "${GREEN_COLOR_BG}Migrations removed"
