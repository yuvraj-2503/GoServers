package errors

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
