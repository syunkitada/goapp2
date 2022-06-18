package errors

import "fmt"

type ConflictError struct {
	err string
}

func NewConflictErrorf(format string, args ...interface{}) (err *ConflictError) {
	err = &ConflictError{
		err: fmt.Sprintf(format, args...),
	}
	return
}

func (self *ConflictError) Error() (err string) {
	return "ConflictError: " + self.err
}

func IsConflictError(err error) bool {
	switch err.(type) {
	case *ConflictError:
		return true
	}
	return false
}
