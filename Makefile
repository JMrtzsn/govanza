# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOFIX=$(GOCMD) fix
GORUN=$(GOCMD) run
GOINSTALL=$(GOCMD) install
GOLANGCILINTCMD=$(GOPATH)/bin/golangci-lint
GOLANGCILINTRUN=$(GOLANGCILINTCMD) run -v
GOGENERATE=$(GOCMD) generate
GOSECCMD=$(GOPATH)/bin/gosec
GOSECRUN=$(GOSECCMD) --quiet
GOVULNCHECKCMD=$(GOPATH)/bin/govulncheck
GOVULNCHECKRUN=$(GOVULNCHECKCMD)

# Library name
LIBRARY_NAME=govanza

# Output directory
OUTPUT_DIR=bin

# Binary output file
BINARY_OUTPUT=$(OUTPUT_DIR)/$(LIBRARY_NAME).a

# Version
VERSION = $(shell git describe --abbrev=0 --tags | sed 's/[v]//g')

all: test build

build:
	@echo "Building $(LIBRARY_NAME) library..."
	$(GOBUILD) -o $(BINARY_OUTPUT) $(wildcard *.go)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: fmt
fmt:
	$(GOFMT) ./...

fix:
	@echo :"Running fix"
	$(GOFIX) ./...

clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIR)

deps:
	@echo "Downloading dependencies..."
	$(GOCMD) mod download

tidy:
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

install-golangci-lint:
	@echo "Installing golangci-lint..."
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install-golangci-lint
	$(GOLANGCILINTRUN) ./...

install-gosec:
	@echo "Installing gosec..."
	$(GOINSTALL) github.com/securego/gosec/v2/cmd/gosec@latest

sec: install-gosec
	@echo "Running gosec..."
	${GOSECRUN} ./...


install-govulncheck:
	@echo "Installing go vulncheck..."
	$(GOINSTALL) golang.org/x/vuln/cmd/govulncheck@latest

govulncheck: install-govulncheck
	@echo "Running go vulncheck..."
	$(GOVULNCHECKRUN) ./...

all: tidy build fmt fix lint test sec govulncheck
