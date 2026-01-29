SHELL = /bin/bash

# Pact CLI Path (relative to provider directory)
export PATH := $(PWD)/pact/bin:$(PATH)
export PATH

# Provider Configuration
export PROVIDER_NAME = GoUserService

# Pact Broker Configuration
export PACT_BROKER_PROTO = http
export PACT_BROKER_URL = localhost:8079
export PACT_BROKER_USERNAME = pact_workshop
export PACT_BROKER_PASSWORD = pact_workshop

# Versioning (from git)
export VERSION_COMMIT?=$(shell git rev-parse HEAD)
export VERSION_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)

# Pact CLI Configuration
PACT_CLI_VERSION = 0.0.2
PACT_STANDALONE_VERSION = 2.4.8
