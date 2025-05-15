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

MESSAGE_COUNT_LEGACY ?= 100
MESSAGE_COUNT_NEW ?= 10

# Default target
.PHONY: all
all: clean
	make GOOS=linux GOARCH=amd64 build deploy clean

# Pre-build steps
.PHONY: pre-build
pre-build:
	mkdir -p $(BUILD_ARCH)

# Build target
.PHONY: build
build: pre-build
	go build -o $(BUILD_ARCH)/$(APP) load.go

# Script generation
.PHONY: script
script: pre-build
	make env-vars > $(BUILD_SCRIPT)

# Environment variable export
.PHONY: env-vars
env-vars:
	@$(DP_CONFIGS)/scripts/secrets-admin $(ENV) $(SUBNET) $(SECRETS_APP) --export '*'
	@echo export TIMEOUT="$(TIMEOUT)"
	@echo export MESSAGE_COUNT_LEGACY="$(MESSAGE_COUNT_LEGACY)"
	@echo export MESSAGE_COUNT_NEW="$(MESSAGE_COUNT_NEW)"

# Deploy target
.PHONY: deploy
deploy: build script
	dp scp $(ENV) $(host_num) -r -- $(BUILD_ARCH) $(host_bin)
	dp ssh $(ENV) $(host_num) -- 'bash -c "cd $(host_bin) && source ./$(APP).sh && ./$(APP)"'

# Clean target
.PHONY: clean
clean: clean-deploy clean-build

# Clean build directory
.PHONY: clean-build
clean-build:
	[[ ! -d $(BUILD) ]] || rm -r $(BUILD)

# Clean deploy
.PHONY: clean-deploy
clean-deploy:
	dp ssh $(ENV) $(host_num) -- bash -c '"[[ ! -d $(host_bin) ]] || rm -r $(host_bin)"'
