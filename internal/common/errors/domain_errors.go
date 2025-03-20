package errors

import (
	"errors"
	"fmt"
)

// Domain-specific error types to streamline error handling

// ErrorType represents the type of an error
type ErrorType string

// Error types
const (
	// NotFound errors
	NotFoundError ErrorType = "NOT_FOUND"

	// Validation errors
	ValidationError ErrorType = "VALIDATION"

	// Authorization errors
	UnauthorizedError ErrorType = "UNAUTHORIZED"
	ForbiddenError    ErrorType = "FORBIDDEN"

	// Conflict errors
	ConflictError ErrorType = "CONFLICT"

	// System errors
	DatabaseError    ErrorType = "DATABASE"
	NetworkError     ErrorType = "NETWORK"
	IntegrationError ErrorType = "INTEGRATION"

	// Input errors
	BadInputError ErrorType = "BAD_INPUT"

	// Business logic errors
	BusinessError ErrorType = "BUSINESS"

	// Unknown errors
	UnknownError ErrorType = "UNKNOWN"
)

// DomainError represents a structured error in the system
type DomainError struct {
	Type    ErrorType
	Message string
	Code    string
	Details map[string]interface{}
	Cause   error
}

// Implement the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (code: %s, cause: %v)", e.Type, e.Message, e.Code, e.Cause)
	}
	return fmt.Sprintf("%s: %s (code: %s)", e.Type, e.Message, e.Code)
}

// Unwrap returns the underlying cause of the error
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Generic error constructors
func NewDomainError(message string, errorType ErrorType, code string, details map[string]interface{}, cause error) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Details: details,
		Cause:   cause,
	}
}

// Specific error constructors
func NewNotFoundError(entity string, id any) *DomainError {
	return NewDomainError(
		fmt.Sprintf("%s with ID %v not found", entity, id),
		NotFoundError,
		"ENTITY_NOT_FOUND",
		map[string]any{"entity": entity, "id": id},
		nil,
	)
}

func NewValidationError(message string, details map[string]any) *DomainError {
	return NewDomainError(
		message,
		ValidationError,
		"VALIDATION_FAILED",
		details,
		nil,
	)
}

func NewUnauthorizedError(message string) *DomainError {
	if message == "" {
		message = "Authorization required"
	}
	return NewDomainError(
		message,
		UnauthorizedError,
		"UNAUTHORIZED",
		nil,
		nil,
	)
}

func NewForbiddenError(message string) *DomainError {
	if message == "" {
		message = "Access forbidden"
	}
	return NewDomainError(
		message,
		ForbiddenError,
		"FORBIDDEN",
		nil,
		nil,
	)
}

func NewConflictError(entity string, details map[string]any) *DomainError {
	return NewDomainError(
		fmt.Sprintf("Conflict with existing %s", entity),
		ConflictError,
		"ENTITY_CONFLICT",
		details,
		nil,
	)
}

func NewDatabaseError(operation string, cause error) *DomainError {
	return NewDomainError(
		fmt.Sprintf("Database error during %s", operation),
		DatabaseError,
		"DATABASE_ERROR",
		map[string]interface{}{"operation": operation},
		cause,
	)
}

func NewNetworkError(operation string, cause error) *DomainError {
	return NewDomainError(
		fmt.Sprintf("Network error during %s", operation),
		NetworkError,
		"NETWORK_ERROR",
		map[string]interface{}{"operation": operation},
		cause,
	)
}

func NewIntegrationError(integration string, operation string, cause error) *DomainError {
	return NewDomainError(
		fmt.Sprintf("Error with %s integration during %s", integration, operation),
		IntegrationError,
		"INTEGRATION_ERROR",
		map[string]interface{}{"integration": integration, "operation": operation},
		cause,
	)
}

func NewBadInputError(message string, details map[string]any) *DomainError {
	return NewDomainError(
		message,
		BadInputError,
		"BAD_INPUT",
		details,
		nil,
	)
}

func NewBusinessError(code string, message string, details map[string]any) *DomainError {
	return NewDomainError(
		message,
		BusinessError,
		code,
		details,
		nil,
	)
}

func NewUnknownError(cause error) *DomainError {
	return NewDomainError(
		"An unknown error occurred",
		UnknownError,
		"UNKNOWN_ERROR",
		nil,
		cause,
	)
}

// Utility functions
func IsDomainError(err error) bool {
	var de *DomainError
	return errors.As(err, &de)
}

func ErrorTypeOf(err error) ErrorType {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Type
	}
	return UnknownError
}

func IsNotFoundErrorDomain(err error) bool {
	return ErrorTypeOf(err) == NotFoundError
}

func IsValidationError(err error) bool {
	return ErrorTypeOf(err) == ValidationError
}

func IsAuthorizationError(err error) bool {
	errType := ErrorTypeOf(err)
	return errType == UnauthorizedError || errType == ForbiddenError
}

func IsConflictError(err error) bool {
	return ErrorTypeOf(err) == ConflictError
}
