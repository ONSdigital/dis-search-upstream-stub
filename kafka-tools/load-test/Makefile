GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Environment and default values
SHELL=bash
APP=dis-search-upstream-stub
BUILD = build

GOOS=linux
GOARCH=amd64

# Default environment variables
ENV?=sandbox
SUBNET?=publishing
TIMEOUT?=30s
host_num?=publishing 3
host_bin=bin-$(APP)

DP_CONFIGS?=../../../dp-configs
SECRETS_APP?=dis-search-upstream-stub

# Derived paths
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BUILD_SCRIPT=$(BUILD_ARCH)/$(APP).sh

# Default message counts
NUM_CONTENT_UPDATED ?= 100
NUM_SEARCH_CONTENT_UPDATED ?= 100
NUM_SEARCH_CONTENT_DELETED ?= 10

.PHONY: all
all: clean build deploy clean ## Clean, build, deploy, and clean again

.PHONY: pre-build
pre-build: ## Create the build directory if it doesn't exist
	mkdir -p $(BUILD_ARCH)

.PHONY: build
build: pre-build ## Build the producer binary for the specified GOOS/GOARCH
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_ARCH)/$(APP) main.go

.PHONY: script
script: pre-build ## Generate environment-specific script for deployment
	make env-vars > $(BUILD_SCRIPT)

.PHONY: env-vars
env-vars: ## Export necessary environment variables and secrets for deployment
	@$(DP_CONFIGS)/scripts/secrets-admin $(ENV) $(SUBNET) $(SECRETS_APP) --export '*'
	@echo export TIMEOUT="$(TIMEOUT)"
	@echo export NUM_CONTENT_UPDATED="$(NUM_CONTENT_UPDATED)"
	@echo export NUM_SEARCH_CONTENT_UPDATED="$(NUM_SEARCH_CONTENT_UPDATED)"
	@echo export NUM_SEARCH_CONTENT_DELETED="$(NUM_SEARCH_CONTENT_DELETED)"

.PHONY: deploy
deploy: script build ## Deploy the binary and run remotely via dp ssh
	dp scp $(ENV) $(host_num) -r -- $(BUILD_ARCH) $(host_bin)
	dp ssh $(ENV) $(host_num) -- 'bash -c "cd $(host_bin) && source ./$(APP).sh && \
		./$(APP) \
			-num-content-updated=$$NUM_CONTENT_UPDATED \
			-num-search-content-updated=$$NUM_SEARCH_CONTENT_UPDATED \
			-num-search-content-deleted=$$NUM_SEARCH_CONTENT_DELETED"'

.PHONY: run
run: ## Run locally using current Go environment and flag defaults
	go run load.go \
		-num-content-updated=$(NUM_CONTENT_UPDATED) \
		-num-search-content-updated=$(NUM_SEARCH_CONTENT_UPDATED) \
		-num-search-content-deleted=$(NUM_SEARCH_CONTENT_DELETED)

.PHONY: run-fast
run-fast: ## Run locally with smaller test counts
	go run load.go -num-content-updated=1 -num-search-content-updated=1 -num-search-content-deleted=1


.PHONY: clean
clean: clean-deploy clean-build ## Cleans the build and deploy directories

.PHONY: clean-build
clean-build: ## Deletes the build directory if it exists
	[[ ! -d $(BUILD) ]] || rm -r $(BUILD)

.PHONY: clean-deploy
clean-deploy: ## Deletes the deploy directory on the target machine
	dp ssh $(ENV) $(host_num) -- bash -c '"[[ ! -d $(host_bin) ]] || rm -r $(host_bin)"'

.PHONY: help
help: ## Show help page for list of make targets
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)