package errors

// import (
// 	"fmt"
// )

// // InfrastructureError represents a low-level error
// type InfrastructureError struct {
// 	Message string
// 	Cause   error
// }

// // Implement the error interface
// func (e *InfrastructureError) Error() string {
// 	if e.Cause != nil {
// 		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
// 	}
// 	return e.Message
// }

// // Unwrap returns the underlying cause of the error
// func (e *InfrastructureError) Unwrap() error {
// 	return e.Cause
// }

// // Infrastructure-specific error constructors
// func NewDBConnectionErrorInfra(cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: "Failed to establish database connection",
// 		Cause:   cause,
// 	}
// }

// func NewDBTransactionErrorInfra(cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: "Database transaction failed",
// 		Cause:   cause,
// 	}
// }

// func NewDBQueryErrorInfra(query string, cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: fmt.Sprintf("Error executing query: %s", query),
// 		Cause:   cause,
// 	}
// }

// func NewCacheErrorInfra(operation string, cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: fmt.Sprintf("Cache error during %s", operation),
// 		Cause:   cause,
// 	}
// }

// func NewNetworkErrorInfra(endpoint string, cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: fmt.Sprintf("Network error while calling %s", endpoint),
// 		Cause:   cause,
// 	}
// }

// func NewFileErrorInfra(operation string, cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: fmt.Sprintf("File system error during %s", operation),
// 		Cause:   cause,
// 	}
// }

// func NewIntegrationErrorInfra(service string, cause error) *InfrastructureError {
// 	return &InfrastructureError{
// 		Message: fmt.Sprintf("Error with external integration: %s", service),
// 		Cause:   cause,
// 	}
// }
