SHELL = /bin/bash

# Pact CLI Path (relative to consumer directory)
export PATH := $(PWD)/pact/bin:$(PATH)
export PATH

# Consumer Configuration
export CONSUMER_NAME = GoAdminService
export PACT_DIR = $(PWD)/pacts
export LOG_DIR = $(PWD)/log

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
