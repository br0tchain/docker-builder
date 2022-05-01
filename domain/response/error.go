package response

import (
	"fmt"
	"net/http"
)

//ErrorResponse types
const (
	ErrBuildImage            = "BuildImageError"
	ErrSecurity              = "SecurityError"
	ErrRequestInvalidPayload = "RequestInvalidPayload"
	ErrUndefined             = "UndefinedError"
)

var (
	errCodes = map[string]int{
		ErrRequestInvalidPayload: http.StatusBadRequest,
		ErrSecurity:              http.StatusForbidden,
	}
)

//ErrorResponse : domain error (used to type error)
type ErrorResponse struct {
	ErrorType    string `json:"errorType"`    // The error type
	ErrorMessage string `json:"errorMessage"` // The error message
}

//ErrorResponse : return error formatted message
func (err ErrorResponse) Error() string {
	return fmt.Sprintf("%s : %s", err.ErrorType, err.ErrorMessage)
}

// GetHTTPErrorCode : return corresponding http code
func (err ErrorResponse) GetHTTPErrorCode() int {
	var code int
	var ok bool
	if code, ok = errCodes[err.ErrorType]; !ok {
		code = http.StatusInternalServerError
	}
	return code
}

//NewError : new domain error
func NewError(errorType string, errorMessage string) ErrorResponse {
	return ErrorResponse{
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

//Wrap : wrap error in new domain error
func Wrap(err error, errorType string, message string) ErrorResponse {
	return ErrorResponse{
		ErrorType:    errorType,
		ErrorMessage: fmt.Sprintf("%s; %s", message, err),
	}
}

//Wrapf : wrap error in new domain error formatted
func Wrapf(err error, errorType string, message string, a ...interface{}) ErrorResponse {
	mess := fmt.Sprintf(message, a...)
	return Wrap(err, errorType, mess)
}
