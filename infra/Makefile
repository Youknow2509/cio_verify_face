# =============================
# Make file for Helm operations
# =============================

# Variables
NAMESPACE := cio-verify-face-dev
KUBE := kubectl
# Database postgres
P_USERNAME ?= postgres
P_PASSWORD ?= root1234
P_URL ?= localhost:5433
P_DIR_MIGRATION_SQL := ../server/database/migration/sql/
P_DSN := postgresql://$(P_USERNAME):$(P_PASSWORD)@$(P_URL)/cio_verify_face
GOOSE := goose
# Database ScyllaDB
MIGRATE=migrate
S_USERNAME ?= cassandra
S_PASSWORD ?= root1234
S_URL ?= localhost:9042
S_MIGRATE_PATH_SCHEMA ?= ../server/database/migration/cql/
S_MIGRATE_SOURCE ?= file://${S_MIGRATE_PATH_SCHEMA}
S_MIGRATE_DRIVER ?= cassandra
S_MIGRATE_URL ?= ${S_URL}/cio_verify_face?username=${S_USERNAME}&password=${S_PASSWORD}&sslmode=disable
S_MIGRATE_DB_DSN ?= "${S_MIGRATE_DRIVER}://${S_MIGRATE_URL}"

# Phony targets and default target
.PHONY: help install_dependencies

# Help target
.DEFAULT_GOAL := help

help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

mirgation_scylladb: ## Run migration for ScyllaDB
	@echo "Running ScyllaDB migration..."
	@echo "Using USERNAME=$(S_USERNAME)"
	@echo "Using URL=$(S_URL)"
	${MIGRATE} -path ${S_MIGRATE_PATH_SCHEMA} -database ${S_MIGRATE_DB_DSN} up

migration_postgres:  ## Run migration for Postgres
	@echo "Running Postgres migration..."
	@echo "Using USERNAME=$(P_USERNAME)"
	@echo "Using URL=$(P_URL)"
	$(GOOSE) -dir $(P_DIR_MIGRATION_SQL) postgres "$(P_DSN)" up

check_all: ## K8s check all
	$(KUBE)  get all --namespace $(NAMESPACE)
	$(KUBE) get deployments,pods,services --namespace $(NAMESPACE)

check_status_rollout: ## Check detail status rollout
	$(KUBE) get all --namespace $(NAMESPACE)
	$(KUBE) get deployments,pods,services --namespace $(NAMESPACE)

check_logs: ## Check logs all
	$(KUBE) get events -n $(NAMESPACE) --sort-by=.metadata.creationTimestamp
