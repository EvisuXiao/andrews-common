package exception

import (
	"fmt"
)

type DbError struct {
	errorString string
}

type CustomError struct {
	errorString string
}

var (
	FAILURE_ERR       = CustomErrWrapper(FAILURE_MSG)
	INVALID_PARAM_ERR = CustomErrWrapper(INVALID_PARAM_MSG)
	DB_ERROR_ERR      = CustomErrWrapper(DB_ERROR_MSG)
	SERVER_ERROR_ERR  = CustomErrWrapper(SERVER_ERROR_MSG)
)

func DbErrWrapper(err error) *DbError {
	if err == nil {
		return nil
	}
	return &DbError{err.Error()}
}

func (e *DbError) Error() string {
	if e == nil {
		return ""
	}
	return e.errorString
}

func CustomErrWrapper(err string, args ...interface{}) *CustomError {
	if err == "" {
		return nil
	}
	return &CustomError{fmt.Sprintf(err, args...)}
}

func (e *CustomError) Error() string {
	if e == nil {
		return ""
	}
	return e.errorString
}
