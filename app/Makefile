SHELL := /bin/bash
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=pagination-app
BUILD_DIR=bin
INTERNAL_CODE=./internal

# Directories
SRC_DIR=./main.go
DOMAIN_DIR=$(INTERNAL_CODE)/domain...
API_DIR=$(INTERNAL_CODE)/api...
PKG_DIR=./pkg/...

# Default target executed when no arguments are given to make
all: test build

# Build the project
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)
	chmod +x $(BUILD_DIR)/$(BINARY_NAME)

# Run the project
run:
	$(BUILD_DIR)/$(BINARY_NAME)

# Test all ./internal packages
test-domain:
	$(GOTEST) -v $(DOMAIN_DIR)
test-api:
	$(GOTEST) -v $(API_DIR)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -u ./...

# Format the code
fmt:
	$(GOCMD) fmt ./...

# Lint the code
lint:
	@golangci-lint run

# Pre-commit hook
pre-commit: fmt lint test


.PHONY: all build run test-domain test-api clean deps fmt lint pre-commit
