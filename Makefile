SHELL := /bin/bash

ROOT := $(CURDIR)
SCRIPTS_DIR := $(ROOT)/scripts
BINARY_DIR := $(ROOT)/bin
PKG_DIR := $(ROOT)/pkg
VENDOR_DIR := $(ROOT)/vendor/golangutils
GO := go

.PHONY: all build clean check-deps process-go-dependencies list-go-dependencies help

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build                    Build windows and linux binaries"
	@echo "  clean                    Remove build outputs"
	@echo "  check-deps               Verify required tools (go)"
	@echo "  process-go-dependencies  Processo all go dependencies for this project"
	@echo "  list-go-dependencies     List all go dependencies for this project"
	@echo

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { echo "[ERROR] Please install golang!"; exit 1; }

process-go-dependencies: check-deps
	@echo ">>> Process for $(ROOT)"
	@cd $(ROOT) && $(GO) mod tidy
	@echo ">>> Process for $(VENDOR_DIR)"
	@cd $(VENDOR_DIR) && $(GO) mod tidy
	@cd $(ROOT)

list-go-dependencies: check-deps
	@echo ">>> List for $(ROOT)"
	@cd $(ROOT) && $(GO) list -m all
	@echo
	@echo ">>> List for $(VENDOR_DIR)"
	@cd $(VENDOR_DIR) && $(GO) list -m all
	@cd $(ROOT)

build: check-deps
	@bash $(SCRIPTS_DIR)/builder.sh $(ROOT)

clean:
	rm -rf $(BINARY_DIR)/
	rm -rf $(PKG_DIR)/