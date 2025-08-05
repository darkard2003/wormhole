SERVER_BINARY_NAME=wormhole-server
OUT_DIR=out

all: build

build:
	@echo "Building $(SERVER_BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR)/$(SERVER_BINARY_NAME) ./cmd/server/main.go

run: build
	@echo "Running $(SERVER_BINARY_NAME)..."
	@$(OUT_DIR)/$(SERVER_BINARY_NAME)

.PHONY: all build run
