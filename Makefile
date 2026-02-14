.PHONY: build run test lint clean

BINARY_NAME=arithmego
BUILD_DIR=bin

# Version info (overridable via environment)
VERSION ?= dev

# ldflags for version injection
LDFLAGS=-ldflags "-X github.com/gurselcakar/arithmego/internal/cli.Version=$(VERSION)"

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
