package kizcool

import "fmt"

// AuthenticationError contains an http error Code and text Message
type AuthenticationError struct {
	Code    int
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%d: %v", e.Code, e.Message)
}
