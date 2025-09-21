package pluggedinkit

import "fmt"

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Details    interface{}
}

func (e *APIError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("API error (%d): %s - %v", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// AuthenticationError represents an authentication error
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	Message    string
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit exceeded: %s (retry after %d seconds)", e.Message, e.RetryAfter)
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s", e.Field, e.Message)
}