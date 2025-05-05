# ----------------------
# Configuration
# ----------------------

COVERAGE_OUT=docs/coverage.out
COVERAGE_HTML=docs/coverage.html
GO=go
LINTER=golangci-lint
LINTER_REMOTE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
LINTER_OPTS=--timeout=2m

# ----------------------
# General Targets
# ----------------------

.PHONY: all check ci lint test-unit test-e2e test-full coverage example-simplest clean build-all build-macos build-windows

all: check

check: lint test-full coverage

build-all: build-macos build-windows

ci: check build-all

clean:
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML) ducto-faker*

# ----------------------
# Linting
# ----------------------

lint:
	@echo "==> Running linter"
	$(LINTER) run $(LINTER_OPTS)

lint-install:
	go install $(LINTER_REMOTE)

# ----------------------
# Testing
# ----------------------

test-unit:
	@echo "==> Running short tests"
	$(GO) test -short -coverpkg=./... -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-full:
	@echo "==> Running all tests"
	@PUBSUB_EMULATOR_HOST=localhost:8085 GOOGLE_CLOUD_PROJECT=test-project $(GO) test -coverpkg=./... -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-e2e:
	@echo "==> Running full tests"
	$(GO) test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v -run E2E ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

coverage:
	@echo "==> Generating coverage HTML report"
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

# ----------------------
# CLI
# ----------------------

build-macos:
	@echo "==> Building macOS CLI"
	GOOS=darwin GOARCH=arm64 $(GO) build -o ducto-faker ./cmd/ducto-faker

build-windows:
	@echo "==> Building Windows CLI"
	GOOS=windows GOARCH=amd64 $(GO) build -o ducto-faker.exe ./cmd/ducto-faker
