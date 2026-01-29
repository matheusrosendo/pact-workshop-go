# Consumer Makefile
# This Makefile is completely self-contained and ready to be extracted to a separate repository
# All dependencies and infrastructure setup are included

# Include shared configuration
include make/config.mk

# Paths
CONSUMER_DIR = $(shell pwd)
PACT_CLI = $(CONSUMER_DIR)/pact/bin/pact-broker

# Go Pact library version
PACT_GO_VERSION = v2.0.8

.PHONY: help
help:
	@echo "Consumer Pact Makefile - $(CONSUMER_NAME)"
	@echo ""
	@echo "Infrastructure:"
	@echo "  install            - Install Pact Go library"
	@echo "  install-cli        - Install Pact CLI tools"
	@echo ""
	@echo "Testing:"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-pact          - Run Pact consumer tests (generates contracts)"
	@echo ""
	@echo "Publishing:"
	@echo "  publish            - Publish contracts to Pact Broker"
	@echo ""
	@echo "Deployment:"
	@echo "  can-i-deploy       - Check if safe to deploy to an environment"
	@echo "  can-i-deploy-dev   - Check if safe to deploy to dev"
	@echo "  can-i-deploy-prod  - Check if safe to deploy to production"
	@echo "  record-deploy-dev  - Record deployment to dev environment"
	@echo "  record-deploy-prod - Record deployment to production environment"
	@echo ""
	@echo "Utilities:"
	@echo "  clean              - Clean generated pact files"
	@echo ""
	@echo "Full workflow:"
	@echo "  make install-cli && make test-unit && make test-pact && make publish"

.PHONY: test-unit
test-unit:
	@echo "--- 🔨 Running Consumer Unit Tests"
	go test -tags=unit -count=1 ./internal/client -run 'TestClientUnit' -v

.PHONY: test-pact
test-pact:
	@echo "--- 🔨 Running Consumer Pact Tests"
	@echo "This will generate contracts for all providers in: $(PACT_DIR)"
	@mkdir -p $(PACT_DIR)
	@mkdir -p $(LOG_DIR)
	CONSUMER_NAME=$(CONSUMER_NAME) \
	PROVIDER_NAME=GoUserService \
	PACT_DIR=$(PACT_DIR) \
	LOG_DIR=$(LOG_DIR) \
	go test -tags=integration -count=1 ./internal/client -run 'TestClientPact' -v

# Alias for backward compatibility
.PHONY: consumer
consumer: test-pact

.PHONY: publish
publish:
	@echo "--- 📝 Publishing Pact Contracts to Broker"
	@echo "Consumer: $(CONSUMER_NAME)"
	@echo "Version: $(VERSION_COMMIT)"
	@echo "Branch: $(VERSION_BRANCH)"
	$(PACT_CLI) publish $(PACT_DIR) \
		--consumer-app-version $(VERSION_COMMIT) \
		--branch $(VERSION_BRANCH) \
		--broker-base-url $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL) \
		--broker-username $(PACT_BROKER_USERNAME) \
		--broker-password $(PACT_BROKER_PASSWORD)
	@echo ""
	@echo "✅ Contracts published successfully!"
	@echo "View at: $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL)"

.PHONY: can-i-deploy
can-i-deploy:
	@if [ -z "$(ENV)" ]; then \
		echo "❌ Error: ENV variable is required"; \
		echo "Usage: make can-i-deploy ENV=production"; \
		echo "       make can-i-deploy ENV=staging"; \
		echo "       make can-i-deploy ENV=dev"; \
		exit 1; \
	fi
	@echo "--- ✅ Checking if $(CONSUMER_NAME) can be deployed to $(ENV)"
	$(PACT_CLI) can-i-deploy \
		--pacticipant $(CONSUMER_NAME) \
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
	@echo "--- ✅ Recording deployment of $(CONSUMER_NAME) to $(ENV)"
	$(PACT_CLI) record-deployment \
		--pacticipant $(CONSUMER_NAME) \
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

.PHONY: clean
clean:
	@echo "--- 🧹 Cleaning generated Pact files"
	rm -rf $(PACT_DIR)/*.json
	rm -rf $(LOG_DIR)/*
	@echo "✅ Clean complete"

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

# CI/CD workflow target
.PHONY: ci
ci: test-unit test-pact publish
	@echo "✅ Consumer CI pipeline complete"

# Full deployment workflow
.PHONY: deploy-to-prod
deploy-to-prod: can-i-deploy-prod
	@echo "✅ Safe to deploy to production"
	@echo "After deploying, run: make record-deploy-prod"
