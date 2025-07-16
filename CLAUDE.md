# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Testing
```bash
go test -v                    # Run all tests with verbose output
go test -v -run TestExecute   # Run specific test
go test -cover               # Run tests with coverage
```

### Building
```bash
go build                     # Build the package
go mod tidy                  # Clean up dependencies
```

### Formatting
```bash
go fmt ./...                 # Format all Go files
go vet ./...                 # Static analysis
```

## Architecture

This is a Go library that provides a wrapper around the Gemini CLI command for programmatic access to Google's Gemini AI.

### Core Components

- **Client** (`client.go`): Main client structure with configuration options
  - Model selection support (gemini-2.5-flash, gemini-2.5-pro, etc.)
  - Timeout configuration
  - Logging integration
  - Authentication error detection

- **Logger Interface** (`logger.go`): Pluggable logging interface
  - NoOpLogger for silent operation
  - LoggerAdapter for integration with external loggers

- **Adapter** (`adapter.go`): Logger adapter for integrating with external logging systems

### Key Features

- **Command Execution**: Executes `gemini -m <model> -p <prompt>` with proper error handling
- **Output Filtering**: Automatically filters system messages from responses
- **Authentication Detection**: Detects and handles API credential issues
- **Timeout Support**: Configurable command timeouts
- **Model Support**: Supports different Gemini models with default gemini-2.5-flash

### Default Configuration

- Default timeout: 30 seconds
- Default model: "gemini-2.5-flash"
- Default logger: NoOpLogger (silent)

### Error Handling

The library handles:
- Empty prompts
- Command not found errors
- Authentication failures
- Timeout errors
- Output parsing errors

### Testing Strategy

The codebase uses Test-Driven Development (TDD) with comprehensive test coverage:

- Unit tests for all public functions
- Error condition testing
- Output parsing validation
- Authentication error detection
- Timeout behavior verification
- Model selection testing

Tests are designed to handle cases where the Gemini CLI is not available during development.

## Prerequisites

- Go 1.24.4 or later
- `gemini` CLI command installed and available in PATH
- Valid Gemini API credentials configured

## Package Structure

```
github.com/yubiquita/gemini-cli-wrapper
├── client.go         # Main client implementation
├── client_test.go    # Comprehensive test suite
├── logger.go         # Logger interface and NoOpLogger
├── adapter.go        # Logger adapter for external systems
├── go.mod           # Go module definition
└── README.md        # Documentation and usage examples
```