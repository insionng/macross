package macross

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(StatusNotFound)
	ErrStatusBadRequest            = NewHTTPError(StatusBadRequest)
	ErrUnauthorized                = NewHTTPError(StatusUnauthorized)
	ErrMethodNotAllowed            = NewHTTPError(StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(StatusRequestEntityTooLarge)
	ErrRendererNotRegistered       = errors.New("renderer not registered")
	ErrInvalidRedirectCode         = errors.New("invalid redirect status code")
	ErrCookieNotFound              = errors.New("cookie not found")
)

// Error contains the error information reported by calling Context.Error().
// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Status  int    //`json:"status" xml:"status"`
	Message string //`json:"message" xml:"message"`
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(status int, message ...interface{}) *HTTPError {
	he := &HTTPError{Status: status, Message: StatusText(status)}
	if len(message) > 0 {
		he.Message = fmt.Sprint(message...)
	}
	return he
}

// Error returns the error message.
func (e *HTTPError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code.
func (e *HTTPError) StatusCode() int {
	return e.Status
}
