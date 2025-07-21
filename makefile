BINARY_NAME=wormhole
OUT_DIR=out

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR)/$(BINARY_NAME) main.go

run:
	@echo "Running $(BINARY_NAME) with go run..."
	@go run main.go &

.PHONY: all build run
