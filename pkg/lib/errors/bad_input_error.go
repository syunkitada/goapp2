package errors

import "fmt"

type BadInputError struct {
	err string
}

func NewBadInputErrorf(format string, args ...interface{}) (err *BadInputError) {
	err = &BadInputError{
		err: fmt.Sprintf(format, args...),
	}
	return
}

func (self *BadInputError) Error() (err string) {
	return "BadInputError: " + self.err
}

func IsBadInputError(err error) bool {
	switch err.(type) {
	case *BadInputError:
		return true
	}
	return false
}
