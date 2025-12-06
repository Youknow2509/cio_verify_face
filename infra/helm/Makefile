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

create_all_config_map_for_services: ## Create all config map for services
	@echo "Create config map for service ai"
	$(KUBECTL) create configmap service-ai-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_ai/config/.env 
	
	@echo "Create config map for service analytic"
	$(KUBECTL) create configmap service-analytic-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_analytic/config/config.yaml

	@echo "Create config map for service attendance"
	$(KUBECTL) create configmap service-attendance-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_attendance/config/config.yaml

	@echo "Create config map for service auth"
	$(KUBECTL) create configmap service-auth-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_auth/config/config.yaml

	@echo "Create config map for service device"
	$(KUBECTL) create configmap service-device-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_device/config/config.yaml

	@echo "Create config map for service identity"
	$(KUBECTL) create configmap service-identity-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_identity/config/.env
	
	@echo "Create config map for service notify"
	$(KUBECTL) create configmap service-notify-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_notify/config/config.yaml
	
	@echo "Create config map for service profile update"
	$(KUBECTL) create configmap service-profile-update-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_profile_update/config/config.yaml
	
	@echo "Create config map for service workforce"
	$(KUBECTL) create configmap service-workforce-config-map \
		--namespace=$(NAMESPACE) \
		--from-file=./charts/service_workforce/config/config.yaml

delete_all_config_map_for_services: ## Delete all config map for services
	LIST_CONFIG_MAPS="service-ai-config-map service-analytic-config-map service-attendance-config-map service-auth-config-map service-device-config-map service-identity-config-map service-notify-config-map service-profile-update-config-map service-workforce-config-map"; \
	for config_map in $$LIST_CONFIG_MAPS; do \
		echo "Deleting config map: $$config_map"; \
		$(KUBECTL) delete configmap $$config_map --namespace=$(NAMESPACE) || echo "Config map $$config_map does not exist"; \
	done

deploy_all_services: ## Deploy all services
	@echo "All services deployed"
	$(MAKE) deploy_gate_way
	$(MAKE) deploy_service_ai
	$(MAKE) deploy_service_analytic
	$(MAKE) deploy_service_attendance
	$(MAKE) deploy_service_auth
	$(MAKE) deploy_service_device
	$(MAKE) deploy_service_identity
	$(MAKE) deploy_service_notify
	$(MAKE) deploy_service_profile_update
	$(MAKE) deploy_service_workforce

deploy_gate_way: ## Deploy Gateway
	@echo "Deploy Gateway"
	$(HELM) upgrade --install gateway $(CHARTS_DIR)/gateway --namespace $(NAMESPACE) --create-namespace

deploy_service_identity: ## Deploy Service Identity
	@echo "Deploy Service Identity"
	$(HELM) upgrade --install service-identity $(CHARTS_DIR)/service_identity --namespace $(NAMESPACE) --create-namespace

deploy_service_ai: ## Deploy Service AI
	@echo "Deploy Service AI"
	$(HELM) upgrade --install service-ai $(CHARTS_DIR)/service_ai --namespace $(NAMESPACE) --create-namespace

deploy_service_workforce: ## Deploy Service Workforce
	@echo "Deploy Service Workforce"
	$(HELM) upgrade --install service-workforce $(CHARTS_DIR)/service_workforce --namespace $(NAMESPACE) --create-namespace

deploy_service_profile_update: ## Deploy Service Profile Update
	@echo "Deploy Service Profile Update"
	$(HELM) upgrade --install service-profile-update $(CHARTS_DIR)/service_profile_update --namespace $(NAMESPACE) --create-namespace

deploy_service_notify: ## Deploy Service Notify
	@echo "Deploy Service Notify"
	$(HELM) upgrade --install service-notify $(CHARTS_DIR)/service_notify --namespace $(NAMESPACE) --create-namespace

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
