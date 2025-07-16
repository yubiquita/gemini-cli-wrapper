# Gemini CLI Wrapper

A Go library that provides a convenient wrapper around the Gemini CLI command, enabling programmatic access to Google's Gemini AI through the command-line interface.

## Features

- **Simple API**: Easy-to-use client interface for Gemini CLI interactions
- **Model Selection**: Support for different Gemini models (flash, pro, etc.)
- **Working Directory**: Custom working directory support for context-specific configurations
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

### Working Directory Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

func main() {
    // Execute with custom working directory
    config := geminicli.Config{
        WorkingDirectory: "/path/to/project",
        Model:           "gemini-2.5-flash",
        Timeout:         60 * time.Second,
    }
    
    client := geminicli.NewClientWithConfig(config)
    
    // This will execute gemini command in the specified directory
    response, err := client.Execute("Analyze the files in this directory")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
    
    // Convenience function for working directory
    response2, err := geminicli.ExecuteWithWorkingDirectory(
        "What programming language is used in this project?",
        "/path/to/project",
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response2)
    
    // Full configuration with all options
    response3, err := geminicli.ExecuteWithFullConfig(
        "Review the code quality",
        "gemini-2.5-pro",              // model
        "/path/to/project",            // working directory
        120*time.Second,               // timeout
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response3)
}
```

#### Relative Path Resolution

When using a `WorkingDirectory`, relative paths in prompts are automatically resolved to absolute paths based on your current working directory:

```go
// Current directory: /home/user/myproject
// Working directory: /config/gemini

config := geminicli.Config{
    WorkingDirectory: "/config/gemini",
}
client := geminicli.NewClientWithConfig(config)

// Relative paths are resolved automatically
response, err := client.Execute("Analyze ./main.go and ./config.json")
// Becomes: "Analyze /home/user/myproject/main.go and /home/user/myproject/config.json"
```

**Key behaviors:**
- `./file.txt` → `/current/directory/file.txt`
- `../parent/file.txt` → `/current/parent/file.txt`
- `subdir/file.txt` → `/current/directory/subdir/file.txt`
- `/absolute/path/file.txt` → `/absolute/path/file.txt` (unchanged)
- When `WorkingDirectory` is not set, Gemini runs in your current directory (no path resolution needed)

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

## Custom Configuration Directory

The library supports executing Gemini commands in custom directories, enabling you to use directory-specific configuration files and context files.

### Setting Up Custom Configuration

1. **Create a custom configuration directory:**
   ```bash
   mkdir -p ~/.config/my_cli_app
   ```

2. **Create a settings.json file:**
   ```bash
   cat > ~/.config/my_cli_app/settings.json << 'EOF'
   {
     "contextFileName": ["CONTEXT.md", "PROJECT.md"],
     "model": "gemini-2.5-pro"
   }
   EOF
   ```

   **To disable context files completely:**
   ```bash
   cat > ~/.config/my_cli_app/settings.json << 'EOF'
   {
     "contextFileName": [],
     "model": "gemini-2.5-flash"
   }
   EOF
   ```

3. **Create context files:**
   ```bash
   cat > ~/.config/my_cli_app/CONTEXT.md << 'EOF'
   # Project Context
   
   This is a custom CLI application that processes user queries.
   
   ## Guidelines
   - Provide concise responses
   - Focus on practical solutions
   - Include code examples when relevant
   EOF
   ```

4. **Use the custom configuration:**
   ```go
   package main
   
   import (
       "fmt"
       "log"
       "os"
       "path/filepath"
       
       "github.com/yubiquita/gemini-cli-wrapper"
   )
   
   func main() {
       // Get user's home directory
       home, err := os.UserHomeDir()
       if err != nil {
           log.Fatal(err)
       }
       
       // Custom configuration directory
       configDir := filepath.Join(home, ".config", "my_cli_app")
       
       // Execute with custom configuration directory
       response, err := geminicli.ExecuteWithWorkingDirectory(
           "How do I handle file uploads in Go?",
           configDir,
       )
       if err != nil {
           log.Fatal(err)
       }
       fmt.Println(response)
   }
   ```

### Advanced Configuration Example

For more complex setups, you can create multiple configuration directories for different contexts:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/yubiquita/gemini-cli-wrapper"
)

func main() {
    home, _ := os.UserHomeDir()
    
    // Different configurations for different contexts
    contexts := map[string]string{
        "backend":    filepath.Join(home, ".config", "backend_helper"),
        "frontend":   filepath.Join(home, ".config", "frontend_helper"),
        "devops":     filepath.Join(home, ".config", "devops_helper"),
        "no_context": filepath.Join(home, ".config", "clean_gemini"),  // No context files
    }
    
    // Use backend context
    response, err := geminicli.ExecuteWithWorkingDirectory(
        "How do I optimize database queries?",
        contexts["backend"],
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Backend advice:", response)
    
    // Use frontend context
    response2, err := geminicli.ExecuteWithWorkingDirectory(
        "How do I implement responsive design?",
        contexts["frontend"],
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Frontend advice:", response2)
    
    // Use clean context (no context files)
    response3, err := geminicli.ExecuteWithWorkingDirectory(
        "What is the capital of France?",
        contexts["no_context"],
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Clean response:", response3)
}
```

### Configuration Directory Structure

A typical configuration directory structure might look like:

```
~/.config/my_cli_app/
├── settings.json          # Gemini CLI configuration
├── CONTEXT.md            # General context file
├── PROJECT.md            # Project-specific context
├── GUIDELINES.md         # Response guidelines
└── examples/
    ├── api_examples.md   # API usage examples
    └── best_practices.md # Best practices guide

# Example: Clean Gemini configuration (no context)
~/.config/clean_gemini/
└── settings.json          # Only contains {"contextFileName": []}
```

### Key Benefits

- **Isolated Contexts**: Each configuration directory provides isolated context
- **Reusable Configurations**: Share configuration directories across projects
- **Domain-Specific Responses**: Get responses tailored to specific domains
- **Clean Responses**: Use empty context for unbiased, general responses
- **Version Control**: Configuration directories can be version controlled
- **Team Collaboration**: Share configurations with team members

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

#### `ExecuteWithWorkingDirectory(prompt, workingDirectory string) (string, error)`

Executes a Gemini command with a specified working directory using a default client.

#### `ExecuteWithWorkingDirectoryAndTimeout(prompt, workingDirectory string, timeout time.Duration) (string, error)`

Executes a Gemini command with a specified working directory and custom timeout using a default client.

#### `ExecuteWithFullConfig(prompt, model, workingDirectory string, timeout time.Duration) (string, error)`

Executes a Gemini command with all configuration options (model, working directory, and timeout) using a default client.

### Configuration

```go
type Config struct {
    Logger           Logger        // Custom logger implementation
    Timeout          time.Duration // Command execution timeout
    Model            string        // Model name (default: "gemini-2.5-flash")
    WorkingDirectory string        // Working directory for command execution
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

## Troubleshooting

### Common Issues

#### Working Directory Not Found
```
Error: failed to execute Gemini command: chdir /path/to/directory: no such file or directory
```
**Solution**: Ensure the working directory exists before execution:
```go
import "os"

if _, err := os.Stat("/path/to/directory"); os.IsNotExist(err) {
    err := os.MkdirAll("/path/to/directory", 0755)
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Permission Denied
```
Error: failed to execute Gemini command: chdir /path/to/directory: permission denied
```
**Solution**: Check directory permissions:
```bash
chmod 755 /path/to/directory
```

#### Configuration Not Loading
If your custom `settings.json` is not being loaded:

1. **Verify file location**: Ensure the file is in the working directory
2. **Check JSON syntax**: Validate your JSON with a linter
3. **Verify permissions**: Ensure the file is readable

**To disable all context files:**
```json
{
  "contextFileName": []
}
```

```bash
# Check if file exists and is readable
ls -la /path/to/directory/settings.json

# Validate JSON syntax
cat /path/to/directory/settings.json | jq .
```

#### Context Files Not Found
If your context files aren't being loaded:

1. **Check `contextFileName` in settings.json**:
   ```json
   {
     "contextFileName": ["CONTEXT.md", "PROJECT.md"]
   }
   ```

2. **Verify file paths**: Ensure context files exist in the working directory:
   ```bash
   ls -la /path/to/directory/CONTEXT.md
   ```

3. **Check file permissions**:
   ```bash
   chmod 644 /path/to/directory/*.md
   ```

#### Relative Path Issues
If relative paths in prompts aren't working as expected:

1. **Check current directory**: The library resolves relative paths based on where you run your program
   ```bash
   pwd  # Shows your current directory
   ```

2. **Use absolute paths**: When in doubt, use absolute paths in prompts
   ```go
   // Instead of: "Analyze ./main.go"
   // Use: "Analyze /full/path/to/main.go"
   ```

3. **Debug path resolution**: Enable debug logging to see how paths are resolved
   ```go
   config := geminicli.Config{
       WorkingDirectory: "/config/gemini",
       Logger:          &DebugLogger{},
   }
   ```

4. **Natural language file references**: Remember that phrases like "check the main file" won't be resolved - use explicit paths

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
import "log"

// Custom logger that prints debug info
type DebugLogger struct{}

func (l *DebugLogger) DebugWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

func (l *DebugLogger) InfoWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *DebugLogger) WarnWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[WARN] %s %v", msg, keysAndValues)
}

func (l *DebugLogger) ErrorWith(msg string, keysAndValues ...interface{}) {
    log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

func main() {
    config := geminicli.Config{
        Logger:           &DebugLogger{},
        WorkingDirectory: "/path/to/directory",
    }
    
    client := geminicli.NewClientWithConfig(config)
    response, err := client.Execute("test prompt")
    // Check debug output for detailed execution information
}
```

### Best Practices

1. **Always check errors**: Handle all returned errors appropriately
2. **Use absolute paths**: Prefer absolute paths for working directories
3. **Validate directories**: Check if directories exist before use
4. **Keep configurations simple**: Start with basic configurations and add complexity gradually
5. **Version control configurations**: Keep your configuration directories in version control

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
