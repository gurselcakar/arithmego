.PHONY: build run test lint clean

BINARY_NAME=arithmego
BUILD_DIR=bin

# Version info (overridable via environment)
VERSION ?= dev
COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ldflags for version injection
LDFLAGS=-ldflags "-X github.com/gurselcakar/arithmego/internal/cli.Version=$(VERSION) \
                  -X github.com/gurselcakar/arithmego/internal/cli.CommitSHA=$(COMMIT_SHA) \
                  -X github.com/gurselcakar/arithmego/internal/cli.BuildDate=$(BUILD_DATE)"

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/arithmego

# Build with version info (for releases)
build-release:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/arithmego

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean
