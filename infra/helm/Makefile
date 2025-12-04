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

help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

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

install_dependencies: ## Install Helm chart dependencies
	$(HELM) dependency update $(UMBRELLA_DIR)

clean: ## Clean Helm chart dependencies
	$(HELM) dependency clean $(UMBRELLA_DIR)

lint: ## Lint Helm charts
	$(HELM) lint $(UMBRELLA_DIR)

install_platform: ## Install platform Helm chart
	$(HELM) install platform $(UMBRELLA_DIR) --namespace $(NAMESPACE) --create-namespace

upgrade_platform: ## Upgrade platform Helm chart
	$(HELM) upgrade platform $(UMBRELLA_DIR) --namespace $(NAMESPACE) --install --create-namespace

check_release_platform: ## Check release Helm in namespace
	$(HELM) list --namespace $(NAMESPACE)
	$(HELM) status platform --namespace $(NAMESPACE)

