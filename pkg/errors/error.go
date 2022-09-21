package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	E5xxINTERNAL      = "internal"
	E5xxUNAVAILABLE   = "service_unavailable"
	E6xxNETWORK       = "network"
	E4xxCLIENTSIDE    = "client_side"
	E4xxUNAUTHORIZED  = "unauthorized"
	E4xxNOTFOUND      = "not_found"
	E4xxUNPROCESSABLE = "unprocessable_entity"
)

var (
	status = map[string]int{
		E5xxINTERNAL:      http.StatusInternalServerError,
		E5xxUNAVAILABLE:   http.StatusServiceUnavailable,
		E4xxUNPROCESSABLE: http.StatusUnprocessableEntity,
		E4xxNOTFOUND:      http.StatusNotFound,
		E4xxCLIENTSIDE:    http.StatusBadRequest,
		E4xxUNAUTHORIZED:  http.StatusUnauthorized,
	}
)

type Error struct {
	Code    string
	Message string
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return E5xxINTERNAL.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return E5xxINTERNAL
}

// Is unwraps an application error and returns if match code.
// Non-application errors always return false.
func Is(code string, err error) bool {
	var e *Error
	if err == nil {
		return false
	} else if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "internal error".
func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "internal error"
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrorStatus(err error) int {
	code := ErrorCode(err)
	return status[code]
}
