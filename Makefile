APP_NAME := schedulord
BIN_DIR := bin
GO := go

.DEFAULT_GOAL := help
.PHONY: help build run test lint vet fmt tidy upgrade clean containerize

help:
	@echo "Available targets:"
	@echo "  build        - compile the project into binary"
	@echo "  run          - build and execute the binary"
	@echo "  test         - run unit tests"
	@echo "  lint         - run linter"
	@echo "  vet          - check for static errors with go vet"
	@echo "  fmt          - format code"
	@echo "  tidy         - tidy go.mod and go.sum"
	@echo "  upgrade      - upgrade all dependencies"
	@echo "  clean        - remove build artifacts"
	@echo "  containerize - build Podman image from Containerfile"
	@echo "  help         - show this help message"

build: clean
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

run: build
	@echo "Running $(APP_NAME)"
	./$(BIN_DIR)/$(APP_NAME)

test:
	$(GO) test ./...

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

fmt:
	gofumpt -w .
	oxfmt

tidy:
	$(GO) mod tidy

upgrade:
	$(GO) get -u ./...
	$(GO) mod tidy

clean:
	@rm -rf $(BIN_DIR)

containerize:
	podman build --file Containerfile --ignorefile .containerignore -t $(APP_NAME) .
