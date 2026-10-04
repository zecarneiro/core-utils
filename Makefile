SHELL := /bin/bash

# APP Info
NAME := core-utils
APP_VERSION := 2.6.6
DISPLAY_NAME := "CoreUtils"
# Make file data
GO := go
ROOT := $(CURDIR)
SCRIPTS_DIR := $(ROOT)/scripts
SO_TYPE := "linux"
CMD_NAME ?=""


BINARY_DIR := $(ROOT)/bin
PKG_DIR := $(ROOT)/pkg
INSTALLER_PREFIX_DIR := $(ROOT)/coreutils
INSTALLER_DIR := $(INSTALLER_PREFIX_DIR)-$(COREUTILS_VERSION)
INSTALLER_WIN_DIR := $(INSTALLER_DIR)/coreutils-win
INSTALLER_LINUX_DIR := $(INSTALLER_DIR)/coreutils-linux
EXTERNAL_DIR := $(ROOT)/vendor/golangutils
COVERAGE_DIR := $(ROOT)/cover

.PHONY: all build deploy check-deps clean help

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build						Build windows and linux binaries"
	@echo "  deploy						Generate installer packages"
	@echo "  clean						Remove build outputs"
	@echo "  check-deps					Verify required tools (go)"
	@echo


check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { echo "[ERROR] Please install golang!"; exit 1; }

build: check-deps
	@cd "$(ROOT)"
	@bash $(SCRIPTS_DIR)/build.sh "$(CMD_NAME)"

deploy: check-deps
	@cd "$(ROOT)"
	@bash $(SCRIPTS_DIR)/deploy.sh "$(SO_TYPE)" "$(NAME)" "$(APP_VERSION)" "$(DISPLAY_NAME)"

clean:
	@bash $(SCRIPTS_DIR)/cleaner.sh
