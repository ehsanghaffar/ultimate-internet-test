package utils

import "fmt"

// TestError represents an error that occurred during testing
type TestError struct {
	TestType string // The type of test that failed (e.g., "HTTP", "Ping", "Speed", "VPN")
	Message  string
	Err      error // The underlying error
}

// Error implements the error interface
func (e *TestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s test failed: %s (underlying: %v)", e.TestType, e.Message, e.Err)
	}
	return fmt.Sprintf("%s test failed: %s", e.TestType, e.Message)
}

// Unwrap returns the underlying error for error chain inspection
func (e *TestError) Unwrap() error {
	return e.Err
}

// NewTestError creates a new TestError
func NewTestError(testType, message string, err error) *TestError {
	return &TestError{
		TestType: testType,
		Message:  message,
		Err:      err,
	}
}

// NetworkError represents a network-related error
type NetworkError struct {
	*TestError
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(testType, message string, err error) *NetworkError {
	return &NetworkError{
		TestError: NewTestError(testType, message, err),
	}
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	*TestError
}

// NewTimeoutError creates a new TimeoutError
func NewTimeoutError(testType, message string) *TimeoutError {
	return &TimeoutError{
		TestError: NewTestError(testType, message, nil),
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	*TestError
}

// NewValidationError creates a new ValidationError
func NewValidationError(testType, message string) *ValidationError {
	return &ValidationError{
		TestError: NewTestError(testType, message, nil),
	}
}

// ParseError represents an error parsing data
type ParseError struct {
	*TestError
}

// NewParseError creates a new ParseError
func NewParseError(testType, message string, err error) *ParseError {
	return &ParseError{
		TestError: NewTestError(testType, message, err),
	}
}
