# Commands
GO=go

# Vars
BIN_DIR=$(CURDIR)/bin
BIN_PATH=$(BIN_DIR)/gofilebuildtags.so

default: help
.PHONY: default

## test: Run unit tests
test:
	$(GO) test ./... -race -count=1 -failfast
.PHONY: test

test-integration:
	$(GO) test ./... -race -count=1 -failfast -tags=integration
.PHONY: test-integration

## build: Build the source as a .so plugin
build:
	$(GO) build -o $(BIN_PATH) -buildmode=plugin ./pkg/linter
.PHONY: build

## clean-bin: Clean the generated binaries
clean-bin:
	rm -rf $(BIN_DIR)/*
.PHONY: clean-bin

help: Makefile
	@echo
	@echo "Available commands:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /' | LANG=C sort
	@echo
.PHONY: help
