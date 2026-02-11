SHELL := /bin/bash

ROOT := $(CURDIR)
SCRIPTS_DIR := $(ROOT)/scripts
BINARY_DIR := $(ROOT)/bin
PKG_DIR := $(ROOT)/pkg
EXTERNAL_DIR := $(ROOT)/vendor/golangutils
COVERAGE_DIR := $(ROOT)/cover
GO := go
CMD_NAME ?=""
TEST_NAME ?=""
TEST_CMD_PATH ?=""
TEST_GOUP_NAME ?=""

.PHONY: all build clean check-deps process-go-dependencies list-go-dependencies run-all-tests run-test cover help

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build [|CMD_NAME]                      Build windows and linux binaries. Ex: CMD_NAME=system-upgrade"
	@echo "  clean                                  Remove build outputs"
	@echo "  check-deps                             Verify required tools (go)"
	@echo "  process-go-dependencies                Processo all go dependencies for this project"
	@echo "  list-go-dependencies                   List all go dependencies for this project"
	@echo "  run-all-tests                          Run all tests"
	@echo "  run-test [|TEST_NAME] [TEST_CMD_PATH]  Run specific test. Ex: TEST_NAME=TestName and TEST_CMD_PATH=package-functions/system-upgrade"
	@echo "  run-group-tests [TEST_GOUP_NAME]       Run tests on given group. Ex: TEST_GOUP_NAME=package-functions"
	@echo "  cover                                  Run all tests with cover generated"
	@echo

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { echo "[ERROR] Please install golang!"; exit 1; }

process-go-dependencies: check-deps
	@echo ">>> Process for $(ROOT)"
	@cd $(ROOT) && $(GO) mod tidy
	@echo ">>> Process for $(EXTERNAL_DIR)"
	@cd $(EXTERNAL_DIR) && $(GO) mod tidy
	@cd $(ROOT)

list-go-dependencies: check-deps
	@echo ">>> List for $(ROOT)"
	@cd $(ROOT) && $(GO) list -m all
	@echo
	@echo ">>> List for $(EXTERNAL_DIR)"
	@cd $(EXTERNAL_DIR) && $(GO) list -m all
	@cd $(ROOT)

build: check-deps
	@bash $(SCRIPTS_DIR)/builder.sh $(CMD_NAME) $(ROOT)

run-all-tests:
	ROOT_DIR=$(ROOT) go test ./...

run-test:
	ROOT_DIR=$(ROOT) go test -v -run $(TEST_NAME) $(ROOT)/cmd/$(TEST_CMD_PATH)

run-group-tests:
	ROOT_DIR=$(ROOT) go test $(ROOT)/cmd/$(TEST_GOUP_NAME)/...

cover:
	@echo "Running all tests with cover..."
	@mkdir -p $(COVERAGE_DIR)
	@ROOT_DIR=$(ROOT) go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

	@echo "----------------------------------------------------"
	@echo "Cover generated with success!"
	@echo "File HTML: $(COVERAGE_DIR)/coverage.html"
	@echo "File de dados: $(COVERAGE_DIR)/coverage.out"
	@echo "----------------------------------------------------"

	# 3. Optional: Open in browser (Linux/WSL)
	# xdg-open $(COVERAGE_DIR)/coverage.html || open $(COVERAGE_DIR)/coverage.html

clean:
	rm -rf $(BINARY_DIR)/
	rm -rf $(PKG_DIR)/
	rm -rf $(COVERAGE_DIR)/
