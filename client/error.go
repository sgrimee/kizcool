package client

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
