package rck

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned from the RCK API.
type APIError struct {
	StatusCode   int
	ResponseData *UnifiedAPIResponse
}

func (e *APIError) Error() string {
	if e.ResponseData != nil && e.ResponseData.Error != "" {
		return fmt.Sprintf("API error (status %d): %s - %s", e.StatusCode, e.ResponseData.Error, e.ResponseData.Details)
	}
	return fmt.Sprintf("API error with status code %d", e.StatusCode)
}

// ValidationError represents an error in validating request parameters.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NetworkError represents a client-side network error.
type NetworkError struct {
	Message       string
	OriginalError error
}

func (e *NetworkError) Error() string {
	if e.OriginalError != nil {
		return fmt.Sprintf("network error: %s (caused by: %v)", e.Message, e.OriginalError)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NetworkError
// This is to prevent infinite recursion if we try to marshal the error itself.
func (e *NetworkError) UnmarshalJSON(data []byte) error {
	var msg string
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	e.Message = msg
	return nil
}

// Predefined error instances
var (
	ErrAuthentication = &APIError{StatusCode: 401}
	ErrAPIKeyRequired = &ValidationError{Field: "APIKey", Message: "API key is required"}
)

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}
