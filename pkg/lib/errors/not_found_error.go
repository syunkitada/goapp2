package errors

import "fmt"

type NotFoundError struct {
	err string
}

func NewNotFoundErrorf(format string, args ...interface{}) (err *NotFoundError) {
	err = &NotFoundError{
		err: fmt.Sprintf(format, args...),
	}
	return
}

func (self *NotFoundError) Error() (err string) {
	return "NotFoundError: " + self.err
}

func IsNotFoundError(err error) bool {
	switch err.(type) {
	case *NotFoundError:
		return true
	}
	return false
}
