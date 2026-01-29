.PHONY: build run test lint clean

BINARY_NAME=arithmego
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/arithmego

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean
