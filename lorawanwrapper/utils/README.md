# LoRaWAN Wrapper Utils

This directory contains the Go implementation of the LoRaWAN wrapper with CGO fallback support.

## Build System

The build system automatically detects CGO availability and chooses the appropriate build method:

- **CGO Available**: Builds `lorawanWrapper.so` (shared library)
- **CGO Unavailable**: Builds `lorawanWrapper` (command-line tool)

## Build Commands

```bash
# Auto-detect and build (recommended)
make build

# Test CGO availability
make test-cgo

# Force CGO build (requires C compiler)
make cgo

# Force non-CGO build (portable)
make nocgo

# Clean build artifacts
make clean

# Install Go dependencies
make deps

# Show help
make help
```

## Build Outputs

- `lorawanWrapper.so` + `lorawanWrapper.h` - CGO shared library
- `lorawanWrapper` - Command-line tool (non-CGO)

## Environment Variables

- `CGO_ENABLED=0` - Disables CGO, forces command-line build
- `ENVIRONMENT=PROD` - Sets production log level

## File Structure

### CGO Build Files
- `lorawanWrapper.go` - Main CGO wrapper functions
- `sessionKeysGenerator.go` - Session key generation (CGO)
- `jsonUnmarshaler.go` - JSON parsing (CGO)
- `micGenerator.go` - MIC generation (CGO)
- `main_cgo.go` - CGO build main function

### Non-CGO Build Files
- `cmd_main.go` - Command-line interface
- `cmd_functions.go` - Command-line function implementations
- `sessionKeysGenerator_common.go` - Session key generation (non-CGO)
- `jsonUnmarshaler_common.go` - JSON parsing (non-CGO)
- `micGenerator_common.go` - MIC generation (non-CGO)
- `main_nocgo.go` - Non-CGO build main function

### Common Files
- `common.go` - Shared utility functions
- `hashGenerator.go` - Hash generation (works with both builds)
- `Makefile` - Build system

## Python Integration

The Python wrapper (`../LorawanWrapperNew.py`) automatically detects which build is available and uses the appropriate interface:

- CGO mode: Direct C function calls via ctypes
- Non-CGO mode: Subprocess calls to command-line tool

Both modes provide identical functionality and API.