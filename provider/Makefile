# Provider Makefile
# This Makefile is completely self-contained and ready to be extracted to a separate repository
# All dependencies and infrastructure setup are included

# Include shared configuration
include make/config.mk

# Paths
PROVIDER_DIR = $(shell pwd)
LOG_DIR = $(PROVIDER_DIR)/log
PACT_CLI = $(PROVIDER_DIR)/pact/bin/pact-broker

# Go Pact library version
PACT_GO_VERSION = v2.0.8

.PHONY: help
help:
	@echo "Provider Pact Makefile - $(PROVIDER_NAME)"
	@echo ""
	@echo "Infrastructure:"
	@echo "  install            - Install Pact Go library"
	@echo "  install-cli        - Install Pact CLI tools"
	@echo ""
	@echo "Testing:"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-pact          - Run Pact provider verification tests"
	@echo ""
	@echo "Deployment:"
	@echo "  can-i-deploy       - Check if safe to deploy to an environment"
	@echo "  can-i-deploy-dev   - Check if safe to deploy to dev"
	@echo "  can-i-deploy-prod  - Check if safe to deploy to production"
	@echo "  record-deploy-dev  - Record deployment to dev environment"
	@echo "  record-deploy-prod - Record deployment to production environment"
	@echo ""
	@echo "Service:"
	@echo "  run                - Run the provider service locally"
	@echo ""
	@echo "Full workflow:"
	@echo "  make install-cli && make test-unit && make test-pact"

.PHONY: test-unit
test-unit:
	@echo "--- 🔨 Running Provider Unit Tests"
	go test -tags=unit -count=1 ./... -v

.PHONY: test-pact
test-pact:
	@echo "--- 🔨 Running Provider Pact Verification Tests"
	@echo "This will verify contracts from all consumers"
	@echo "Provider: $(PROVIDER_NAME)"
	@echo "Version: $(VERSION_COMMIT)"
	@echo "Branch: $(VERSION_BRANCH)"
	@mkdir -p $(LOG_DIR)
	PACT_TEST=true \
	VERSION_COMMIT=$(VERSION_COMMIT) \
	VERSION_BRANCH=$(VERSION_BRANCH) \
	PACT_BROKER_PROTO=$(PACT_BROKER_PROTO) \
	PACT_BROKER_URL=$(PACT_BROKER_URL) \
	PACT_BROKER_USERNAME=$(PACT_BROKER_USERNAME) \
	PACT_BROKER_PASSWORD=$(PACT_BROKER_PASSWORD) \
	go test -tags=integration -count=1 ./internal/pact/provider -run 'TestPactProvider' -v

# Alias for backward compatibility
.PHONY: provider
provider: test-pact

.PHONY: can-i-deploy
can-i-deploy:
	@if [ -z "$(ENV)" ]; then \
		echo "❌ Error: ENV variable is required"; \
		echo "Usage: make can-i-deploy ENV=production"; \
		echo "       make can-i-deploy ENV=staging"; \
		echo "       make can-i-deploy ENV=dev"; \
		exit 1; \
	fi
	@echo "--- ✅ Checking if $(PROVIDER_NAME) can be deployed to $(ENV)"
	@echo "This checks if all consumers in $(ENV) are compatible with this provider version"
	$(PACT_CLI) can-i-deploy \
		--pacticipant $(PROVIDER_NAME) \
		--version $(VERSION_COMMIT) \
		--to-environment $(ENV) \
		--broker-base-url $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL) \
		--broker-username $(PACT_BROKER_USERNAME) \
		--broker-password $(PACT_BROKER_PASSWORD)

.PHONY: can-i-deploy-dev
can-i-deploy-dev:
	@$(MAKE) can-i-deploy ENV=dev

.PHONY: can-i-deploy-prod
can-i-deploy-prod:
	@$(MAKE) can-i-deploy ENV=production

.PHONY: record-deploy
record-deploy:
	@if [ -z "$(ENV)" ]; then \
		echo "❌ Error: ENV variable is required"; \
		echo "Usage: make record-deploy ENV=production"; \
		echo "       make record-deploy ENV=staging"; \
		echo "       make record-deploy ENV=dev"; \
		exit 1; \
	fi
	@echo "--- ✅ Recording deployment of $(PROVIDER_NAME) to $(ENV)"
	$(PACT_CLI) record-deployment \
		--pacticipant $(PROVIDER_NAME) \
		--version $(VERSION_COMMIT) \
		--environment $(ENV) \
		--broker-base-url $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL) \
		--broker-username $(PACT_BROKER_USERNAME) \
		--broker-password $(PACT_BROKER_PASSWORD)
	@echo "✅ Deployment recorded successfully!"

.PHONY: record-deploy-dev
record-deploy-dev:
	@$(MAKE) record-deploy ENV=dev

.PHONY: record-deploy-prod
record-deploy-prod:
	@$(MAKE) record-deploy ENV=production

.PHONY: run
run:
	@echo "--- 🚀 Starting Provider Service"
	@echo "Service will be available at http://localhost:8080"
	go run cmd/main.go

# Alias for backward compatibility
.PHONY: run-provider
run-provider: run

# ============================================================================
# Infrastructure Setup (self-contained)
# ============================================================================

.PHONY: install
install:
	@echo "--- 📦 Installing Pact Go library"
	go get github.com/pact-foundation/pact-go/v2@$(PACT_GO_VERSION)
	go mod tidy
	@echo "✅ Pact Go library installed"

.PHONY: install-cli
install-cli:
	@echo "--- 📦 Installing Pact CLI tools"
	@mkdir -p pact/bin
	@if [ ! -f pact/bin/pact-broker ]; then \
		echo "Downloading Pact CLI..."; \
		curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash -s -- -d pact -v $(PACT_STANDALONE_VERSION); \
		echo "✅ Pact CLI installed at pact/bin/"; \
	else \
		echo "✅ Pact CLI already installed"; \
	fi

# Note: Pact Broker is typically shared infrastructure
# In a real distributed environment, it would be deployed separately
# This target is kept for workshop/local development purposes only
.PHONY: broker-info
broker-info:
	@echo "--- ℹ️  Pact Broker Information"
	@echo "URL: $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL)"
	@echo "Username: $(PACT_BROKER_USERNAME)"
	@echo ""
	@echo "Note: In production, the Pact Broker is shared infrastructure"
	@echo "deployed separately from consumer/provider services."

# ============================================================================
# CI/CD Workflows
# ============================================================================

.PHONY: ci
ci: test-unit test-pact
	@echo "✅ Provider CI pipeline complete"

.PHONY: deploy-to-prod
deploy-to-prod: can-i-deploy-prod
	@echo "✅ Safe to deploy to production"
	@echo "After deploying, run: make record-deploy-prod"
