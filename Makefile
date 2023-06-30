# Go parameters
GO       := go
GOBUILD  := $(GO) build
GOCLEAN  := $(GO) clean
GOTEST   := $(GO) test
GOGET    := $(GO) get
GOFMT    := $(GO) fmt
GOFIX    := $(GO) fix
GORUN    := $(GO) run
GOINSTALL:= $(GO) install

# Library name
LIBRARY_NAME := govanza

# Output directory
OUTPUT_DIR := bin

# Binary output file
BINARY_OUTPUT := $(OUTPUT_DIR)/$(LIBRARY_NAME)

# Version
VERSION := $(shell git describe --abbrev=0 --tags | sed 's/[v]//g')

build:
	@echo "Building $(LIBRARY_NAME) library..."
	$(GOBUILD) -o $(BINARY_OUTPUT) $(LIBRARY_NAME).go

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

fmt:
	@echo "Running gofmt..."
	$(GOFMT) ./...

fix:
	@echo "Running fix..."
	$(GOFIX) ./...

clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIR)

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

tidy:
	@echo "Tidying up..."
	$(GO) mod tidy

install-golangci-lint:
	@echo "Installing golangci-lint..."
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install-golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run -v ./...

install-gosec:
	@echo "Installing gosec..."
	$(GOINSTALL) github.com/securego/gosec/v2/cmd/gosec@latest

sec: install-gosec
	@echo "Running gosec..."
	gosec ./...

install-govulncheck:
	@echo "Installing go vulncheck..."
	$(GOINSTALL) golang.org/x/vuln/cmd/govulncheck@latest

govulncheck: install-govulncheck
	@echo "Running go vulncheck..."
	govulncheck ./...

all: tidy build fmt fix lint test sec govulncheck
