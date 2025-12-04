# =============================
# Make file for Helm operations
# =============================

# Variables
NAMESPACE := cio-verify-face-dev
KUBE := kubectl

# Phony targets and default target
.PHONY: help install_dependencies

# Help target
.DEFAULT_GOAL := help

help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

check_all: ## K8s check all
	$(KUBE)  get all --namespace $(NAMESPACE)
	$(KUBE) get deployments,pods,services --namespace $(NAMESPACE)

check_status_rollout: ## Check detail status rollout
	$(KUBE) get all --namespace $(NAMESPACE)
	$(KUBE) get deployments,pods,services --namespace $(NAMESPACE)

check_logs: ## Check logs all
	$(KUBE) get events -n $(NAMESPACE) --sort-by=.metadata.creationTimestamp
