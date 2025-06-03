GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Environment and default values
SHELL=bash
APP=dis-search-upstream-stub

# Default environment variables
ENV?=sandbox
SUBNET?=publishing
TIMEOUT?=30s

host_num?=publishing 3
host_bin=bin-$(APP)

GOOS=linux
GOARCH=amd64

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BUILD_SCRIPT=$(BUILD_ARCH)/$(APP).sh

DP_CONFIGS?=../../../dp-configs

SECRETS_APP?=dis-search-upstream-stub

MESSAGE_COUNT_LEGACY ?= 6000
MESSAGE_COUNT_NEW ?= 100

.PHONY: all
all: clean ## Runs clean, build, deploy, and clean steps in sequence
	make GOOS=linux GOARCH=amd64 build deploy clean

.PHONY: pre-build
pre-build: ## Creates the build directory if it doesn't exist
	mkdir -p $(BUILD_ARCH)

.PHONY: build
build: pre-build ## Builds the application for the specified GOOS/GOARCH and outputs to the build directory
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_ARCH)/$(APP) load.go

.PHONY: script
script: pre-build ## Generates environment-specific scripts for deployment
	make env-vars > $(BUILD_SCRIPT)

.PHONY: env-vars
env-vars: ## Exports necessary environment variables and secrets for deployment
	@$(DP_CONFIGS)/scripts/secrets-admin $(ENV) $(SUBNET) $(SECRETS_APP) --export '*'
	@echo export TIMEOUT="$(TIMEOUT)"
	@echo export MESSAGE_COUNT_LEGACY="$(MESSAGE_COUNT_LEGACY)"
	@echo export MESSAGE_COUNT_NEW="$(MESSAGE_COUNT_NEW)"

.PHONY: deploy
deploy: script build ## Runs the script, builds the app, and deploys to the specified environment
	dp scp $(ENV) $(host_num) -r -- $(BUILD_ARCH) $(host_bin)
	dp ssh $(ENV) $(host_num) -- 'bash -c "cd $(host_bin) && source ./$(APP).sh && ./$(APP)"'

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