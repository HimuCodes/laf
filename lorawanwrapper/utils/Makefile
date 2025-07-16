# Makefile for LoRaWAN Wrapper with CGO fallback support

BINARY_NAME=lorawanWrapper
SHARED_LIB=lorawanWrapper.so
CMD_BINARY=lorawanWrapper

# Go build variables
GO_BUILD=go build
CGO_FILES=lorawanWrapper.go sessionKeysGenerator.go jsonUnmarshaler.go micGenerator.go hashGenerator.go common.go main_cgo.go
NOCGO_FILES=cmd_main.go cmd_functions.go sessionKeysGenerator_common.go jsonUnmarshaler_common.go micGenerator_common.go hashGenerator.go common.go main_nocgo.go

# CGO detection
CGO_ENABLED ?= 1
export CGO_ENABLED

# Default target
.PHONY: all
all: build

# Build target - tries CGO first, falls back to command-line tool
.PHONY: build
build:
	@echo "Building LoRaWAN Wrapper..."
	@if [ "$(CGO_ENABLED)" = "0" ]; then \
		echo "CGO disabled by environment variable - building command-line tool"; \
		$(MAKE) build-nocgo; \
	elif $(MAKE) build-cgo 2>/dev/null; then \
		echo "✓ CGO build successful - created shared library"; \
	else \
		echo "✗ CGO build failed - falling back to command-line tool"; \
		$(MAKE) build-nocgo; \
	fi

# CGO build (shared library)
.PHONY: build-cgo
build-cgo:
	@echo "Attempting CGO build..."
	@CGO_ENABLED=1 $(GO_BUILD) -tags cgo -o $(SHARED_LIB) -buildmode=c-shared $(CGO_FILES)

# Non-CGO build (command-line tool)
.PHONY: build-nocgo
build-nocgo:
	@echo "Building command-line tool (no CGO)..."
	@CGO_ENABLED=0 $(GO_BUILD) -tags "!cgo" -o $(CMD_BINARY) $(NOCGO_FILES)

# Force CGO build
.PHONY: cgo
cgo:
	@echo "Forcing CGO build..."
	@CGO_ENABLED=1 $(GO_BUILD) -tags cgo -o $(SHARED_LIB) -buildmode=c-shared $(CGO_FILES)

# Force non-CGO build
.PHONY: nocgo
nocgo:
	@echo "Forcing non-CGO build..."
	@CGO_ENABLED=0 $(GO_BUILD) -tags "!cgo" -o $(CMD_BINARY) $(NOCGO_FILES)

# Test CGO availability
.PHONY: test-cgo
test-cgo:
	@echo "Testing CGO availability..."
	@if CGO_ENABLED=1 go build -tags cgo -o /tmp/test_cgo -buildmode=c-shared $(CGO_FILES) 2>/dev/null; then \
		echo "✓ CGO is available"; \
		rm -f /tmp/test_cgo /tmp/test_cgo.h; \
	else \
		echo "✗ CGO is not available"; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(SHARED_LIB) $(CMD_BINARY) *.h

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing Go dependencies..."
	@go mod tidy
	@go mod download

# Help target
.PHONY: help
help:
	@echo "LoRaWAN Wrapper Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build      - Build with CGO if available, fallback to command-line tool"
	@echo "  cgo        - Force CGO build (shared library)"
	@echo "  nocgo      - Force non-CGO build (command-line tool)"
	@echo "  test-cgo   - Test if CGO is available"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install Go dependencies"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  CGO_ENABLED - Set to 0 to disable CGO (default: 1)"
	@echo ""
	@echo "Output Files:"
	@echo "  $(SHARED_LIB) - Shared library (CGO build)"
	@echo "  $(CMD_BINARY) - Command-line tool (non-CGO build)"

.DEFAULT_GOAL := help