BINARY_NAME := watcher
MAIN_PKG    := github.com/Mrilki/CLIServicesWatcher/cmd/CLIServicesWatcher

.PHONY: build run test test-race clean

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PKG)

run:
	@echo "Running..."
	@go run $(MAIN_PKG)

test:
	@echo "Running tests..."
	@go test -v ./...

test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

lint:
	@echo "Running linter..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run


clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe

help:
	@echo "  build      Compile binary"
	@echo "  run        Run without building"
	@echo "  test       Run tests (no race)"
	@echo "  test-race  Run tests with race detector"
	@echo "  lint       Run golangci-lint"
	@echo "  clean      Remove artifacts"
	@echo "  help       Show this message"