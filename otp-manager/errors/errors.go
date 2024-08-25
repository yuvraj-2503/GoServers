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
	code    string
	message string
}

func NewOtpError(code string, message string) *OtpError {
	return &OtpError{
		code:    code,
		message: message,
	}
}

func (e *OtpError) Error() string {
	return fmt.Sprintf("OTP Error, Code %s, Message %s", e.code, e.message)
}

const (
	EXPIRED        = "EXPIRED"
	NOT_FOUND      = "NOT_FOUND"
	INCORRECT      = "INCORRECT"
	LIMIT_EXCEEDED = "LIMIT_EXCEEDED"
)
