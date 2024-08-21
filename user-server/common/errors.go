package common

type AlreadyExistsError struct {
	Message string
}

func (e *AlreadyExistsError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
