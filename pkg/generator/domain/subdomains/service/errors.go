package service

import "fmt"

// Domain errors
var (
	ErrServiceNameRequired = fmt.Errorf("service name is required")
	ErrPortsRequired       = fmt.Errorf("at least one port is required")
)

// ErrInvalidPort represents an invalid port configuration error
type ErrInvalidPort struct {
	Index   int
	Message string
}

func (e *ErrInvalidPort) Error() string {
	return fmt.Sprintf("invalid port at index %d: %s", e.Index, e.Message)
}

// NewErrInvalidPort creates a new invalid port error
func NewErrInvalidPort(index int, message string) error {
	return &ErrInvalidPort{
		Index:   index,
		Message: message,
	}
}
