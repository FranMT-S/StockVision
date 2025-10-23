package models

// ApiError represents an API error type
type ApiError struct {
	Message string `json:"message"`
	Caused  error  `json:"-"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(message string, err error) *ApiError {
	return &ApiError{Message: message, Caused: err}
}

func NewResponseError(msg string) map[string]any {
	return map[string]any{
		"error": msg,
	}
}
