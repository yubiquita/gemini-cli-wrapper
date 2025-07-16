package geminicli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

// Constants
const (
	// Gemini command related
	GeminiCommand    = "gemini"
	GeminiPromptFlag = "-p"
	GeminiModelFlag  = "-m"
	DefaultTimeout   = 30 * time.Second
	DefaultModel     = "gemini-2.5-flash"
	MaxRetries       = 3

	// Error messages
	ErrEmptyPrompt     = "prompt cannot be empty"
	ErrCommandNotFound = "Gemini command not found in PATH"
	ErrCommandFailed   = "failed to execute Gemini command"
	ErrCommandTimeout  = "command timed out"
	ErrCommandStart    = "failed to start command"
	ErrParseOutput     = "failed to parse Gemini output"
	ErrEmptyOutput     = "empty output from Gemini command"
	ErrAuthFailed      = "authentication error: please check your Gemini API credentials"
)

// Client represents a Gemini CLI client
type Client struct {
	logger  Logger
	timeout time.Duration
	model   string // Model name to use
}

// Config represents configuration options for the client
type Config struct {
	Logger  Logger
	Timeout time.Duration
	Model   string // Model name (e.g., "gemini-2.5-flash", "gemini-2.5-pro")
}

// NewClient creates a new Gemini CLI client with default configuration
func NewClient() *Client {
	return &Client{
		logger:  NewNoOpLogger(),
		timeout: DefaultTimeout,
		model:   DefaultModel,
	}
}

// NewClientWithConfig creates a new Gemini CLI client with custom configuration
func NewClientWithConfig(config Config) *Client {
	client := &Client{
		timeout: DefaultTimeout,
		model:   DefaultModel,
	}

	if config.Logger != nil {
		client.logger = config.Logger
	} else {
		client.logger = NewNoOpLogger()
	}

	if config.Timeout > 0 {
		client.timeout = config.Timeout
	}

	if config.Model != "" {
		client.model = config.Model
	}

	return client
}

// Execute executes a Gemini command with the given prompt
func (c *Client) Execute(prompt string) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf(ErrEmptyPrompt)
	}

	// Build command
	cmdArgs := c.buildGeminiCommandWithModel(prompt)
	
	// Log command execution for debugging
	c.logger.DebugWith("Executing Gemini command", "command", cmdArgs[0], "args", cmdArgs[1:])

	// Create command with full path to avoid module resolution issues
	geminiPath, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		c.logger.ErrorWith("Failed to find gemini command", "error", err)
		return "", fmt.Errorf("gemini command not found: %w", err)
	}
	
	c.logger.DebugWith("Using gemini path", "path", geminiPath)
	cmd := exec.Command(geminiPath, cmdArgs[1:]...)
	
	// Set working directory to home directory to avoid module resolution issues
	cmd.Dir = os.Getenv("HOME")
	if cmd.Dir == "" {
		// Fallback to current user's home directory
		if user, err := user.Current(); err == nil {
			cmd.Dir = user.HomeDir
		}
	}
	c.logger.DebugWith("Set working directory", "dir", cmd.Dir)

	// Execute with timeout
	output, err := c.runCommandWithTimeout(cmd, c.timeout)
	if err != nil {
		c.logger.ErrorWith("Gemini command execution failed", "error", err)
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}

	// Parse output
	result, err := c.parseGeminiOutput(output)
	if err != nil {
		c.logger.ErrorWith("Failed to parse Gemini output", "error", err, "output_length", len(output))
		return "", fmt.Errorf("%s: %w", ErrParseOutput, err)
	}

	c.logger.DebugWith("Gemini command completed successfully", "response_length", len(result))
	return result, nil
}

// ExecuteWithTimeout executes Gemini command with custom timeout
func (c *Client) ExecuteWithTimeout(prompt string, timeout time.Duration) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf(ErrEmptyPrompt)
	}

	// Build command
	cmdArgs := c.buildGeminiCommandWithModel(prompt)

	// Create command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// Execute with custom timeout
	output, err := c.runCommandWithTimeout(cmd, timeout)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}

	// Parse output
	result, err := c.parseGeminiOutput(output)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrParseOutput, err)
	}

	return result, nil
}

// ValidateAvailable checks if Gemini command is available
func (c *Client) ValidateAvailable() error {
	_, err := exec.LookPath(GeminiCommand)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCommandNotFound, err)
	}
	return nil
}

// buildGeminiCommand builds the command arguments for Gemini
func (c *Client) buildGeminiCommand(prompt string) []string {
	return []string{GeminiCommand, GeminiPromptFlag, prompt}
}

// buildGeminiCommandWithModel builds the command arguments for Gemini with model specification
func (c *Client) buildGeminiCommandWithModel(prompt string) []string {
	return []string{GeminiCommand, GeminiModelFlag, c.model, GeminiPromptFlag, prompt}
}

// runCommandWithTimeout executes a command with the specified timeout
func (c *Client) runCommandWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	// Start the command
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCommandStart, err)
	}

	// Channel to signal command completion
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for completion or timeout
	select {
	case err := <-done:
		if err != nil {
			// Capture both stdout and stderr for detailed error reporting
			stdoutStr := strings.TrimSpace(string(stdout.Bytes()))
			stderrStr := strings.TrimSpace(string(stderr.Bytes()))
			combined := append(stdout.Bytes(), stderr.Bytes()...)
			
			// Check if it's an authentication error
			if c.detectAuthError(combined) {
				return nil, fmt.Errorf(ErrAuthFailed)
			}
			
			// Create detailed error message
			errorMsg := fmt.Sprintf("command failed: %v", err)
			if stderrStr != "" {
				errorMsg += fmt.Sprintf(" | stderr: %s", stderrStr)
			}
			if stdoutStr != "" {
				errorMsg += fmt.Sprintf(" | stdout: %s", stdoutStr)
			}
			
			return nil, fmt.Errorf("%s", errorMsg)
		}
		return stdout.Bytes(), nil
	case <-time.After(timeout):
		// Kill the process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, fmt.Errorf("%s after %v", ErrCommandTimeout, timeout)
	}
}

// parseGeminiOutput parses the output from Gemini command
func (c *Client) parseGeminiOutput(output []byte) (string, error) {
	if len(output) == 0 {
		return "", fmt.Errorf(ErrEmptyOutput)
	}

	// Convert to string and trim whitespace
	result := strings.TrimSpace(string(output))

	// Filter out authentication and system messages
	result = c.filterGeminiOutput(result)

	if result == "" {
		return "", fmt.Errorf(ErrEmptyOutput)
	}

	return result, nil
}

// detectAuthError detects authentication-related errors in command output
func (c *Client) detectAuthError(output []byte) bool {
	return c.containsAnyKeyword(string(output), c.getAuthErrorKeywords())
}

// getAuthErrorKeywords returns list of authentication error keywords
func (c *Client) getAuthErrorKeywords() []string {
	return []string{
		"authentication failed",
		"invalid api key",
		"permission denied",
		"unauthorized",
		"access denied",
	}
}

// containsAnyKeyword checks if text contains any of the specified keywords (case-insensitive)
func (c *Client) containsAnyKeyword(text string, keywords []string) bool {
	lowerText := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}
	return false
}

// filterGeminiOutput filters out authentication and system messages from Gemini output
func (c *Client) filterGeminiOutput(output string) string {
	// Split output into lines
	lines := strings.Split(output, "\n")
	var filteredLines []string
	
	// Filter patterns that should be removed
	filterPatterns := []string{
		"Loaded cached credentials.",
		"Loading cached credentials",
		"Authenticating",
		"Authentication successful",
		"Connected to Gemini API",
		"Using cached token",
		"Token refreshed",
	}
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		shouldFilter := false
		
		// Check if line matches any filter pattern
		for _, pattern := range filterPatterns {
			if strings.Contains(trimmedLine, pattern) {
				shouldFilter = true
				break
			}
		}
		
		// Keep the line if it doesn't match filter patterns and isn't empty
		if !shouldFilter && trimmedLine != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	
	// Join filtered lines and normalize whitespace
	result := strings.Join(filteredLines, "\n")
	return strings.TrimSpace(result)
}

// Convenience functions for backward compatibility

// Execute executes a Gemini command with the given prompt using default client
func Execute(prompt string) (string, error) {
	client := NewClient()
	return client.Execute(prompt)
}

// ExecuteWithTimeout executes Gemini command with custom timeout using default client
func ExecuteWithTimeout(prompt string, timeout time.Duration) (string, error) {
	client := NewClient()
	return client.ExecuteWithTimeout(prompt, timeout)
}

// ValidateAvailable checks if Gemini command is available using default client
func ValidateAvailable() error {
	client := NewClient()
	return client.ValidateAvailable()
}

// BuildGeminiCommand builds the command arguments for Gemini
func BuildGeminiCommand(prompt string) []string {
	client := NewClient()
	return client.buildGeminiCommand(prompt)
}

// DetectAuthError detects authentication-related errors in command output
func DetectAuthError(output []byte) bool {
	client := NewClient()
	return client.detectAuthError(output)
}

// ParseGeminiOutput parses the output from Gemini command
func ParseGeminiOutput(output []byte) (string, error) {
	client := NewClient()
	return client.parseGeminiOutput(output)
}

// ExecuteWithModel executes a Gemini command with the specified model
func ExecuteWithModel(prompt, model string) (string, error) {
	config := Config{Model: model}
	client := NewClientWithConfig(config)
	return client.Execute(prompt)
}

// ExecuteWithModelAndTimeout executes a Gemini command with the specified model and timeout
func ExecuteWithModelAndTimeout(prompt, model string, timeout time.Duration) (string, error) {
	config := Config{Model: model, Timeout: timeout}
	client := NewClientWithConfig(config)
	return client.Execute(prompt)
}