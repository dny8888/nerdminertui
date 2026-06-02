.PHONY: build test lint ci-pr clean run

# Binary output
BINARY := nerdtui
BUILD_DIR := bin

# Go settings
GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w
CGO := 0

# Linting
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.4

## build: compile the binary with CGO disabled
build:
	CGO_ENABLED=$(CGO) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/tui

## run: build and run without --mock flag
run: build
	./$(BUILD_DIR)/$(BINARY)

## mock: build and run with --mock flag
mock: build
	./$(BUILD_DIR)/$(BINARY) --mock

## test: run all tests
test:
	CGO_ENABLED=$(CGO) $(GO) test -count=1 -timeout 60s ./...

## lint: run golangci-lint
lint:
	$(GOLANGCI_LINT) run ./...

## vet: run go vet
vet:
	$(GO) vet ./...

## vuln: run govulncheck
vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## ci-pr: full CI pipeline (vet + lint + test + vuln)
ci-pr: vet lint test vuln

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR)

## tidy: tidy go modules
tidy:
	$(GO) mod tidy

## help: show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
