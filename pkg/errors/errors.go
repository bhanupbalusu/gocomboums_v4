package errors

import (
	"fmt"

	pkgErrors "github.com/pkg/errors"
)

// AppError represents a custom error structure with an error code
type AppError struct {
	Code    int
	Message string
}

const (
	CodeInternalServerError = 500
	CodeBadRequest          = 400
	CodeNotFound            = 404
)

func (e *AppError) Error() string {
	return e.Message
}

// New creates a new simple error with a message
func New(message string) error {
	return pkgErrors.New(message)
}

// Newf creates a new formatted error
func Newf(format string, args ...interface{}) error {
	return pkgErrors.Errorf(format, args...)
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	return pkgErrors.Wrap(err, message)
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	return pkgErrors.Wrapf(err, format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called
func WithStack(err error) error {
	return pkgErrors.WithStack(err)
}

// NewAppError creates a new application error with a code and a message
func NewAppError(code int, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorf creates a new formatted application error with a code
func NewAppErrorf(code int, format string, args ...interface{}) error {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Cause returns the underlying cause of the error, if possible
func Cause(err error) error {
	return pkgErrors.Cause(err)
}
