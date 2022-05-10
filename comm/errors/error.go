package errors

import (
	"fmt"

	"go.uber.org/multierr"
)

type Error struct {
	Status int
	Err    string
}

func (r *Error) Error() string {
	return fmt.Sprintf("status %d: err %v", r.Status, r.Err)
}

func New(status int, format string, args ...interface{}) error {
	return &Error{
		Status: status,
		Err:    fmt.Sprintf(format, args...),
	}
}

func Is(err error, status int) bool {
	switch er := err.(type) {
	case *Error:
		return er.Status == status
	}
	return false
}

func Combine(errs ...error) error {
	return multierr.Combine(errs...)
}
