package apperror

import (
	"fmt"
	"maps"
	"net/http"
	"runtime"
	"strings"
)

// AppError represents a structured application error
type AppError struct {
	Code       string         `json:"code"`
	Message    string         `json:"-"`
	ApiMessage string         `json:"api_message,omitempty"`
	Origin     string         `json:"-"`
	StatusCode int            `json:"status_code"`
	Values     map[string]any `json:"values,omitempty"`
	Cause      error          `json:"-"`
}

func Define(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		ApiMessage: message,
		StatusCode: http.StatusInternalServerError,
	}
}

func (e *AppError) Error() string {
	var parts []string

	if e.Origin != "" {
		parts = append(parts, fmt.Sprintf("[ %s ]", e.Origin))
	}

	parts = append(parts, fmt.Sprintf("(%s) %s", e.Code, e.Message))

	if e.ApiMessage != "" && e.ApiMessage != e.Message {
		parts = append(parts, fmt.Sprintf("- %s", e.ApiMessage))
	}

	if len(e.Values) > 0 {
		valuesStr := make([]string, 0, len(e.Values))
		for k, v := range e.Values {
			valuesStr = append(valuesStr, fmt.Sprintf("%s: %v", k, v))
		}
		parts = append(parts, fmt.Sprintf("values: {%s}", strings.Join(valuesStr, ", ")))
	}

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}

	return strings.Join(parts, " ")
}

func (e *AppError) Is(target error) bool {
	if other, ok := target.(*AppError); ok {
		return e.Code == other.Code
	}
	return false
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) copy() *AppError {
	newErr := &AppError{
		Code:       e.Code,
		Message:    e.Message,
		ApiMessage: e.ApiMessage,
		Origin:     e.Origin,
		StatusCode: e.StatusCode,
		Cause:      e.Cause,
	}

	if e.Values != nil {
		newErr.Values = make(map[string]any)
		maps.Copy(newErr.Values, e.Values)
	}

	return newErr
}

// WithOrigin adds caller information to the error
func (e *AppError) WithOrigin() *AppError {
	newErr := e.copy()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		newErr.Origin = "unknown origin"
		return newErr
	}

	newErr.Origin = fmt.Sprintf("%s:%d", file, line)
	return newErr
}

// WithMessage sets the main error message
func (e *AppError) WithMessage(msg string) *AppError {
	newErr := e.copy()
	newErr.Message = msg
	return newErr
}

// WithMessagef sets the main error message with formatting
func (e *AppError) WithMessagef(format string, args ...any) *AppError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}

// WithApiMessage adds an API message for the API response
func (e *AppError) WithApiMessage(msg string) *AppError {
	newErr := e.copy()
	newErr.ApiMessage = msg
	return newErr
}

// WithApiMessagef sets the API message with formatting
func (e *AppError) WithApiMessagef(format string, args ...any) *AppError {
	return e.WithApiMessage(fmt.Sprintf(format, args...))
}

// WithValues adds arbitrary data for the API response
func (e *AppError) WithValues(values map[string]any) *AppError {
	newErr := e.copy()
	if newErr.Values == nil {
		newErr.Values = make(map[string]any)
	}
	maps.Copy(newErr.Values, values)
	return newErr
}

// WithValue adds a single key-value pair (convenience method)
func (e *AppError) WithValue(key string, value any) *AppError {
	return e.WithValues(map[string]any{key: value})
}

// WithStatus sets the HTTP status code
func (e *AppError) WithStatus(statusCode int) *AppError {
	newErr := e.copy()
	newErr.StatusCode = statusCode
	return newErr
}

// WithCause wraps another error as the cause
func (e *AppError) WithCause(err error) *AppError {
	newErr := e.copy()
	newErr.Cause = err
	return newErr
}
