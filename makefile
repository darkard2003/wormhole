SERVER_BINARY_NAME=wormhole-server
CLI_BINARY_NAME=wormhole-cli
OUT_DIR=out

all: build-server build-cli

build-server:
	@echo "Building $(SERVER_BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR)/$(SERVER_BINARY_NAME) ./cmd/server/main.go

build-cli:
	@echo "Building $(CLI_BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR)/$(CLI_BINARY_NAME) ./cmd/cli/main.go

build: build-server build-cli

run-server:
	@echo "Running $(SERVER_BINARY_NAME)..."
	@$(OUT_DIR)/$(SERVER_BINARY_NAME)

.PHONY: all build build-server build-cli run-server
