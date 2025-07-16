package geminicli

import (
	"strings"
	"testing"
	"time"
)

// TestExecute tests the main function for executing Gemini commands
func TestExecute(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		expectError bool
		description string
	}{
		{
			name:        "EmptyPrompt",
			prompt:      "",
			expectError: true,
			description: "Should return error for empty prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			result, err := client.Execute(tt.prompt)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for test case '%s': %v", tt.name, err)
				}
				if result == "" {
					t.Errorf("Expected non-empty result for test case '%s'", tt.name)
				}
			}
		})
	}

	// Test with valid prompt but expect command not found error
	t.Run("ValidPromptCommandNotFound", func(t *testing.T) {
		client := NewClient()
		result, err := client.Execute("test prompt")

		if err == nil {
			t.Skip("Gemini command appears to be available, skipping this test")
		}

		if result != "" {
			t.Errorf("Expected empty result when command fails, got '%s'", result)
		}

		// Should get an error about command execution failure
		if !strings.Contains(err.Error(), "failed to execute Gemini command") {
			t.Errorf("Expected execution error, got: %v", err)
		}
	})
}

// TestConvenienceExecute tests the convenience function for executing Gemini commands
func TestConvenienceExecute(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		expectError bool
		description string
	}{
		{
			name:        "EmptyPrompt",
			prompt:      "",
			expectError: true,
			description: "Should return error for empty prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Execute(tt.prompt)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for test case '%s': %v", tt.name, err)
				}
				if result == "" {
					t.Errorf("Expected non-empty result for test case '%s'", tt.name)
				}
			}
		})
	}
}

// TestBuildGeminiCommand tests command construction
func TestBuildGeminiCommand(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		expectedLength int
		description    string
	}{
		{
			name:           "BasicPrompt",
			prompt:         "test prompt",
			expectedLength: 3, // ["gemini", "-p", "test prompt"]
			description:    "Should build basic Gemini command",
		},
		{
			name:           "PromptWithSpaces",
			prompt:         "test prompt with spaces",
			expectedLength: 3,
			description:    "Should handle prompts with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := BuildGeminiCommand(tt.prompt)

			if len(cmd) != tt.expectedLength {
				t.Errorf("Expected command length %d, got %d for test case '%s'",
					tt.expectedLength, len(cmd), tt.name)
			}

			if len(cmd) > 0 && cmd[0] != "gemini" {
				t.Errorf("Expected first argument to be 'gemini', got '%s'", cmd[0])
			}

			if len(cmd) > 1 && cmd[1] != "-p" {
				t.Errorf("Expected second argument to be '-p', got '%s'", cmd[1])
			}

			if len(cmd) > 2 && cmd[2] != tt.prompt {
				t.Errorf("Expected third argument to be '%s', got '%s'", tt.prompt, cmd[2])
			}
		})
	}
}

// TestValidateAvailable tests Gemini availability check
func TestValidateAvailable(t *testing.T) {
	client := NewClient()
	err := client.ValidateAvailable()

	// We expect this to fail initially since the function doesn't exist yet
	if err == nil {
		t.Log("Gemini appears to be available on this system")
	} else {
		t.Logf("Gemini not available (expected during TDD): %v", err)
	}
}

// TestConvenienceValidateAvailable tests the convenience function for Gemini availability check
func TestConvenienceValidateAvailable(t *testing.T) {
	err := ValidateAvailable()

	// We expect this to fail initially since the function doesn't exist yet
	if err == nil {
		t.Log("Gemini appears to be available on this system")
	} else {
		t.Logf("Gemini not available (expected during TDD): %v", err)
	}
}

// TestParseGeminiOutput tests output parsing functionality
func TestParseGeminiOutput(t *testing.T) {
	tests := []struct {
		name         string
		output       []byte
		expectedText string
		expectError  bool
		description  string
	}{
		{
			name:         "ValidOutput",
			output:       []byte("This is a valid Gemini response"),
			expectedText: "This is a valid Gemini response",
			expectError:  false,
			description:  "Should parse valid output correctly",
		},
		{
			name:         "EmptyOutput",
			output:       []byte(""),
			expectedText: "",
			expectError:  true,
			description:  "Should return error for empty output",
		},
		{
			name:         "OutputWithNewlines",
			output:       []byte("Line 1\nLine 2\nLine 3"),
			expectedText: "Line 1\nLine 2\nLine 3",
			expectError:  false,
			description:  "Should handle output with newlines",
		},
		{
			name:         "OutputWithCredentialsMessage",
			output:       []byte("Loaded cached credentials.\nHello, world!"),
			expectedText: "Hello, world!",
			expectError:  false,
			description:  "Should filter out credentials loading message",
		},
		{
			name:         "OutputWithMultipleSystemMessages",
			output:       []byte("Loading cached credentials\nAuthentication successful\nHello, world!\nThis is a response."),
			expectedText: "Hello, world!\nThis is a response.",
			expectError:  false,
			description:  "Should filter out multiple system messages",
		},
		{
			name:         "OnlySystemMessages",
			output:       []byte("Loaded cached credentials.\nAuthenticating\nConnected to Gemini API"),
			expectedText: "",
			expectError:  true,
			description:  "Should return error when only system messages remain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseGeminiOutput(tt.output)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for test case '%s': %v", tt.name, err)
				}
				if result != tt.expectedText {
					t.Errorf("Expected '%s', got '%s' for test case '%s'",
						tt.expectedText, result, tt.name)
				}
			}
		})
	}
}

// TestDetectAuthError tests authentication error detection
func TestDetectAuthError(t *testing.T) {
	tests := []struct {
		name        string
		output      []byte
		expectAuth  bool
		description string
	}{
		{
			name:        "NoAuthError",
			output:      []byte("Normal Gemini response"),
			expectAuth:  false,
			description: "Should not detect auth error in normal output",
		},
		{
			name:        "AuthenticationError",
			output:      []byte("Error: authentication failed"),
			expectAuth:  true,
			description: "Should detect authentication error",
		},
		{
			name:        "APIKeyError",
			output:      []byte("Error: invalid API key"),
			expectAuth:  true,
			description: "Should detect API key error",
		},
		{
			name:        "PermissionError",
			output:      []byte("Error: permission denied"),
			expectAuth:  true,
			description: "Should detect permission error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAuthError(tt.output)

			if result != tt.expectAuth {
				t.Errorf("Expected %v, got %v for test case '%s'",
					tt.expectAuth, result, tt.name)
			}
		})
	}
}

// TestCommandTimeout tests timeout handling
func TestCommandTimeout(t *testing.T) {
	// Test that commands respect timeout settings
	timeout := 1 * time.Second

	t.Run("TimeoutRespected", func(t *testing.T) {
		start := time.Now()

		client := NewClient()
		_, err := client.ExecuteWithTimeout("test prompt", timeout)

		elapsed := time.Since(start)

		if err == nil {
			t.Skip("Function not implemented yet, skipping timeout test")
		}

		// Should not take much longer than the timeout
		if elapsed > timeout+time.Second {
			t.Errorf("Command took too long: %v (timeout was %v)", elapsed, timeout)
		}
	})
}

// TestConvenienceExecuteWithTimeout tests the convenience function for timeout handling
func TestConvenienceExecuteWithTimeout(t *testing.T) {
	// Test that commands respect timeout settings
	timeout := 1 * time.Second

	t.Run("TimeoutRespected", func(t *testing.T) {
		start := time.Now()

		// This will fail until the function is implemented
		_, err := ExecuteWithTimeout("test prompt", timeout)

		elapsed := time.Since(start)

		if err == nil {
			t.Skip("Function not implemented yet, skipping timeout test")
		}

		// Should not take much longer than the timeout
		if elapsed > timeout+time.Second {
			t.Errorf("Command took too long: %v (timeout was %v)", elapsed, timeout)
		}
	})
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	client := NewClient()
	
	if client == nil {
		t.Error("NewClient should return a valid client")
	}
	
	if client.timeout != DefaultTimeout {
		t.Errorf("Expected default timeout %v, got %v", DefaultTimeout, client.timeout)
	}
}

// TestNewClientWithConfig tests client creation with configuration
func TestNewClientWithConfig(t *testing.T) {
	customTimeout := 60 * time.Second
	customLogger := NewNoOpLogger()
	
	config := Config{
		Logger:  customLogger,
		Timeout: customTimeout,
	}
	
	client := NewClientWithConfig(config)
	
	if client == nil {
		t.Error("NewClientWithConfig should return a valid client")
	}
	
	if client.timeout != customTimeout {
		t.Errorf("Expected custom timeout %v, got %v", customTimeout, client.timeout)
	}
}

// TestNoOpLogger tests the no-op logger implementation
func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()
	
	// These should not panic
	logger.DebugWith("test", "key", "value")
	logger.InfoWith("test", "key", "value")
	logger.WarnWith("test", "key", "value")
	logger.ErrorWith("test", "key", "value")
}