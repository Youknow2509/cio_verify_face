.PHONY: help
.PHONY: goose_create goose_up goose_down goose_status goose_fix goose_redo goose_reset goose_clean goose_up_by_one

.ALL: help

# List Variables Goose
GOOSE=goose
GOOSE_NAME ?= cio_verify_face
GOOSE_DRIVER ?= postgres
GOOSE_DB_DSN = "postgresql://postgres:root1234@localhost:5432/cio_verify_face?sslmode=disable"
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

help:
	@echo "Available commands:"
	@echo "\nGoose Migration SQL Commands:"
	@echo "\t goose_create \t Create a new migration"
	@echo "\t goose_up \t Run all available migrations"
	@echo "\t goose_up_by_one \t Run all available migrations by one"
	@echo "\t goose_down \t Rollback the most recent migration"
	@echo "\t goose_status \t Show the status of all migrations"
	@echo "\t goose_fix \t Fix the last migration"
	@echo "\t goose_redo \t Rollback and re-run the most recent migration"
	@echo "\t goose_reset \t Rollback all migrations"
	@echo "\t goose_clean \t Remove all migrations"

	@echo "\nGolang-migrate Migration SQL Commands:"
	@echo "\t migrate_create \t Create a new migration"
	@echo "\t migrate_up \t Run all up migrations"
	@echo "\t migrate_up_n \t Run N up migrations"
	@echo "\t migrate_down \t Rollback the most recent migration"
	@echo "\t migrate_down_n \t Rollback n recent migration"
	@echo "\t migrate_drop \t Drop all migrations"
	@echo "\t migrate_force \t Force a migration to a specific version"

	@echo "  help - Show this help message"

# Golang-migrate - Force a migration to a specific version
migrate_force:
	@echo "Forcing migration to version 1"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} force ${VERSION}
	@echo "${GREEN_COLOR_BG}Migration forced to version ${VERSION}"

# Golang-migrate - Down all migrations
migrate_down:
	@echo "Rolling back the most recent migration"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} down
	@echo "${GREEN_COLOR_BG}Migration rolled back"

# Golang-migrate - Drop all migrations
migrate_drop:
	@echo "Dropping all migrations"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} drop
	@echo "${GREEN_COLOR_BG}All migrations dropped"

# Golang-migrate - Run all up migrations
migrate_up:
	@echo "Running all up migrations"
	${MIGRATE} -path ${MIGRATE_PATH_SCHEMA} -database ${MIGRATE_DB_DSN} up
	@echo "${GREEN_COLOR_BG}Migrations completed"

# Golang-migrate - Create a new migration
migrate_create:
	@echo "Creating a new migration"
	${MIGRATE} create -ext cql -dir ${MIGRATE_PATH_SCHEMA} -seq $(NAME)
	@echo "${GREEN_COLOR_BG}Migration created"

# Goosee - Create a new migration
goose_create:
	@echo "Creating a new migration"
	${GOOSE} -dir ${GOOSE_PATH_SCHEMA} create $(NAME) sql
	@echo "${GREEN_COLOR_BG}Migration created"

# Goose - Migrate the DB up by 1 
goose_up_by_one:
	@echo "Migrate the DB up by 1"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) up-by-one
	@echo "${GREEN_COLOR_BG}Migrations completed"

# Goose - Run all available migrations
goose_up:
	@echo "Running all available migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) up
	@echo "${GREEN_COLOR_BG}Migrations completed"

# Goose - Rollback the most recent migration
goose_down:
	@echo "Rolling back the most recent migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) down
	@echo "${GREEN_COLOR_BG}Migration rolled back"

# Goose - Show the status of all migrations
goose_status:
	@echo "Showing the status of all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) status
	@echo "${GREEN_COLOR_BG}Migration status displayed"

# Goose - Fix the last migration
goose_fix:
	@echo "Fixing the last migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) fix
	@echo "${GREEN_COLOR_BG}Migration fixed"

# Goose - Rollback and re-run the most recent migration
goose_redo:
	@echo "Rolling back and re-running the most recent migration"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) redo
	@echo "${GREEN_COLOR_BG}Migration redone"

# Goose - Rollback all migrations
goose_reset:
	@echo "Rolling back all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) reset
	@echo "${GREEN_COLOR_BG}Migrations reset"

# Goose - Remove all migrations
goose_clean:
	@echo "Removing all migrations"
	goose -dir ${GOOSE_PATH_SCHEMA} ${GOOSE_DRIVER} $(GOOSE_DB_DSN) clean
	@echo "${GREEN_COLOR_BG}Migrations removed"
