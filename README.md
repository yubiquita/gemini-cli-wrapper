# Gemini CLI Wrapper

A Go library that provides a convenient wrapper around the Gemini CLI command, enabling programmatic access to Google's Gemini AI through the command-line interface.

## Features

- **Simple API**: Easy-to-use client interface for Gemini CLI interactions
- **Model Selection**: Support for different Gemini models (flash, pro, etc.)
- **Timeout Support**: Configurable command timeouts to prevent hanging
- **Error Handling**: Comprehensive error detection including authentication failures
- **Output Filtering**: Automatically filters system messages from responses
- **Logging Integration**: Pluggable logging interface for debugging and monitoring
- **Zero Dependencies**: Only uses Go standard library (no external dependencies)

## Installation

```bash
go get github.com/yubiquita/gemini-cli-wrapper
```

## Prerequisites

- Go 1.19 or later
- `gemini` CLI command installed and available in PATH
- Valid Gemini API credentials configured

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

func main() {
    // Simple usage with default client
    response, err := geminicli.Execute("What is Go programming language?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
}
```

### Client-Based Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

func main() {
    // Create a client with custom configuration
    config := geminicli.Config{
        Timeout: 60 * time.Second,  // Custom timeout
        Logger:  geminicli.NewNoOpLogger(), // Or your custom logger
        Model:   "gemini-2.5-pro", // Custom model (default: "gemini-2.5-flash")
    }
    
    client := geminicli.NewClientWithConfig(config)
    
    // Execute commands
    response, err := client.Execute("Explain quantum computing")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
}
```

### Model-Specific Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

func main() {
    // Execute with specific model (convenience function)
    response, err := geminicli.ExecuteWithModel("What is machine learning?", "gemini-2.5-flash")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
    
    // Execute with model and custom timeout
    response2, err := geminicli.ExecuteWithModelAndTimeout(
        "Explain neural networks in detail", 
        "gemini-2.5-pro", 
        120*time.Second,
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response2)
}
```

### Custom Logger Integration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

// Implement the Logger interface
type MyLogger struct{}

func (l *MyLogger) DebugWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

func (l *MyLogger) InfoWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *MyLogger) WarnWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[WARN] %s %v", msg, keysAndValues)
}

func (l *MyLogger) ErrorWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

func main() {
    config := geminicli.Config{
        Logger: &MyLogger{},
    }
    
    client := geminicli.NewClientWithConfig(config)
    response, err := client.Execute("Hello, Gemini!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
}
```

## API Reference

### Client

#### `NewClient() *Client`

Creates a new client with default configuration.

#### `NewClientWithConfig(config Config) *Client`

Creates a new client with custom configuration.

#### `client.Execute(prompt string) (string, error)`

Executes a Gemini command with the given prompt using the client's configuration.

#### `client.ExecuteWithTimeout(prompt string, timeout time.Duration) (string, error)`

Executes a Gemini command with a custom timeout.

#### `client.ValidateAvailable() error`

Checks if the Gemini CLI command is available in the system PATH.

### Convenience Functions

#### `Execute(prompt string) (string, error)`

Executes a Gemini command using a default client.

#### `ExecuteWithTimeout(prompt string, timeout time.Duration) (string, error)`

Executes a Gemini command with custom timeout using a default client.

#### `ExecuteWithModel(prompt, model string) (string, error)`

Executes a Gemini command with a specific model using a default client.

#### `ExecuteWithModelAndTimeout(prompt, model string, timeout time.Duration) (string, error)`

Executes a Gemini command with a specific model and custom timeout using a default client.

#### `ValidateAvailable() error`

Checks if Gemini CLI is available using a default client.

#### `BuildGeminiCommand(prompt string) []string`

Builds command arguments for Gemini CLI.

#### `DetectAuthError(output []byte) bool`

Detects authentication errors in command output.

#### `ParseGeminiOutput(output []byte) (string, error)`

Parses and filters Gemini command output.

### Configuration

```go
type Config struct {
    Logger  Logger        // Custom logger implementation
    Timeout time.Duration // Command execution timeout
    Model   string        // Model name (default: "gemini-2.5-flash")
}
```

### Logger Interface

```go
type Logger interface {
    DebugWith(msg string, keysAndValues ...interface{})
    InfoWith(msg string, keysAndValues ...interface{})
    WarnWith(msg string, keysAndValues ...interface{})
    ErrorWith(msg string, keysAndValues ...interface{})
}
```

## Error Handling

The library provides comprehensive error handling:

- **Empty Prompt**: Returns error when prompt is empty
- **Command Not Found**: Returns error when Gemini CLI is not available
- **Authentication Errors**: Detects and reports API credential issues
- **Timeout Errors**: Reports when commands exceed configured timeout
- **Execution Errors**: Captures and reports command execution failures

## Output Filtering

The library automatically filters out system messages from Gemini responses:

- "Loaded cached credentials."
- "Loading cached credentials"
- "Authenticating"
- "Authentication successful"
- "Connected to Gemini API"
- "Using cached token"
- "Token refreshed"

## Testing

Run the test suite:

```bash
go test -v
```

The library includes comprehensive tests covering:

- Client creation and configuration
- Command execution and error handling
- Output parsing and filtering
- Authentication error detection
- Timeout handling
- Logger integration

## Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Architecture

This library acts as a wrapper around the Gemini CLI command, providing:

```text
Google Gemini API ← Gemini CLI ← This Library ← Your Application
```

The library handles:

- Process management and execution
- Timeout and error handling
- Output parsing and filtering
- Authentication error detection
- Logging integration

## Requirements

- Gemini CLI must be installed and available in PATH
- Valid API credentials must be configured for the Gemini CLI
- Go 1.19+ for development and usage
