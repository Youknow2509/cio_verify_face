# =============================
# Make file for Helm operations
# =============================

# Variables
HELM := helm
CHARTS_DIR := charts
UMBRELLA_DIR := umbrella
NAMESPACE := cio-verify-face-dev
KUBECTL := kubectl
DIR_DASHBOARDS_GRAFANA := ./environment/grafana/dashboards
DIR_DATASOURCES_GRAFANA := ./environment/grafana/provisioning/datasources
SHELL := /bin/bash

# Phony targets and default target
.PHONY: help install_dependencies

# Help target
.DEFAULT_GOAL := help

# =============================================
# Command targets
# =============================================

deploy_service_device: ## Deploy Service Device
	@echo "Deploy Service Device"
	$(HELM) upgrade --install service-device $(CHARTS_DIR)/service_device --namespace $(NAMESPACE) --create-namespace

deploy_service_analytic: ## Deploy Service Analytic
	@echo "Deploy Service Analytic"
	$(HELM) upgrade --install service-analytic $(CHARTS_DIR)/service_analytic --namespace $(NAMESPACE) --create-namespace

deploy_service_attendance: ## Deploy Service Attendance
	@echo "Deploy Service Attendance"
	$(HELM) upgrade --install service-attendance $(CHARTS_DIR)/service_attendance --namespace $(NAMESPACE) --create-namespace

deploy_service_auth: ## Deploy Service Auth
	@echo "Deploy Service Auth"
	$(HELM) upgrade --install service-auth $(CHARTS_DIR)/service_auth --namespace $(NAMESPACE) --create-namespace

create_data_src_grafana: ## Create DataSource for Grafana
	@echo "Create DataSource for Grafana"
	@$(KUBECTL) create secret generic grafana-datasource-secret \
		$(foreach f,$(wildcard $(DIR_DATASOURCES_GRAFANA)/*.yml),--from-file=$(notdir $(basename $(f)))=$(f)) \
        --namespace=$(NAMESPACE) --dry-run=client -o yaml | $(KUBECTL) apply -f -

create_config_map_grafana_dashboards: ## Create ConfigMap for Grafana dashboards
	@echo "Load all dashboards from $(DIR_DASHBOARDS_GRAFANA) into ConfigMap"
	@if [ ! -d "$(DIR_DASHBOARDS_GRAFANA)" ]; then \
        echo "Dir $(DIR_DASHBOARDS_GRAFANA) does not exist!"; \
        exit 1; \
    fi
	@$(KUBECTL) create configmap grafana-dashboards \
        $(foreach f,$(wildcard $(DIR_DASHBOARDS_GRAFANA)/*.json),--from-file=$(notdir $(basename $(f)))=$(f)) \
        --namespace=$(NAMESPACE) --dry-run=client -o yaml | $(KUBECTL) apply -f -

clean: ## Clean Helm chart dependencies
	$(HELM) dependency clean $(UMBRELLA_DIR)

lint: ## Lint Helm charts
	$(HELM) lint $(UMBRELLA_DIR)

deploy_environment: ## Deploy environment: databases, monitoring, cache, message queues, etc.
	$(HELM) install env $(UMBRELLA_DIR) --namespace $(NAMESPACE)

help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
