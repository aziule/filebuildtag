# Commands
GO=go

# Vars
BIN_DIR=$(CURDIR)/bin
BIN_PATH=$(BIN_DIR)/gofilebuildtags.so

default: help
.PHONY: default

## test: Run tests
test:
	$(GO) test ./... -race -count=1 -failfast
.PHONY: test

## build: Build the source as a .so plugin
build:
	$(GO) build -o $(BIN_PATH) -buildmode=plugin ./pkg/linter
.PHONY: build

## clean-bin: Clean the generated binaries
clean-bin:
	rm -rf $(BIN_DIR)/*
.PHONY: clean-bin

## help: Display the available targets
help: Makefile
	@echo
	@echo "Available targets:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /' | LANG=C sort
	@echo
.PHONY: help
