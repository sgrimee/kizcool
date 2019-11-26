package api

import "fmt"

// AuthenticationError contains an http error Code and text Message
type AuthenticationError struct {
	Code    int
	Message string
}

// NewAuthenticationError creates a new AuthenticationError
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Code:    401,
		Message: message,
	}
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%d: %v", e.Code, e.Message)
}

// TooManyRequestsError is returned when the Login call gets a reject from
// the API server because too many authentication attempts too place in a
// short period. In this case a pause needs to be observed before trying again.
type TooManyRequestsError struct {
	Code    int
	Message string
}

// NewTooManyRequestsError creates a new TooManyRequestsError
func NewTooManyRequestsError(message string) *TooManyRequestsError {
	return &TooManyRequestsError{
		Code:    401,
		Message: message,
	}
}

func (e *TooManyRequestsError) Error() string {
	return fmt.Sprintf("%d: %v", e.Code, e.Message)
}

// NoRegisteredEventListenerError happens when an event poll is done on an invalid event listener id
type NoRegisteredEventListenerError struct {
	Code    int
	Message string
}

// NewNoRegisteredEventListenerError creates a new TooManyRequestsError
func NewNoRegisteredEventListenerError(message string) *NoRegisteredEventListenerError {
	return &NoRegisteredEventListenerError{
		Code:    400,
		Message: message,
	}
}

func (e *NoRegisteredEventListenerError) Error() string {
	return fmt.Sprintf("%d: %v", e.Code, e.Message)
}

// JSONError happens when unmarshalling json fails
type JSONError struct {
	Data []byte
	Err  error
}

// Unwrap returns the wrapped error
func (e *JSONError) Unwrap() error { return e.Err }

func (e *JSONError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
