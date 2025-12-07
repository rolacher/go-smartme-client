// errors.go
package smartme

// APIError represents an error returned by the smart-me API.
// You can extend this struct to map the error details from the API.
type APIError struct {
	StatusCode int
	Message    string
	// Body []byte // Useful for debugging
}

func (e *APIError) Error() string {
	return e.Message
}
