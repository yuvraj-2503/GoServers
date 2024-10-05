package errors

import "fmt"

type NetworkConnectionError struct {
}

func (e *NetworkConnectionError) Error() string {
	return "Network connection error"
}

type OtpSendingError struct {
	Message string
}

func (e *OtpSendingError) Error() string {
	return e.Message
}

type OtpDoesNotExistError struct {
}

func (e *OtpDoesNotExistError) Error() string {
	return "otp does not exist"
}

type OtpError struct {
	Code    string
	Message string
}

func NewOtpError(code string, message string) *OtpError {
	return &OtpError{
		Code:    code,
		Message: message,
	}
}

func (e *OtpError) Error() string {
	return fmt.Sprintf("OTP Error, Code %s, Message %s", e.Code, e.Message)
}

const (
	EXPIRED        = "EXPIRED"
	NOT_FOUND      = "NOT_FOUND"
	INCORRECT      = "INCORRECT"
	LIMIT_EXCEEDED = "LIMIT_EXCEEDED"
)
