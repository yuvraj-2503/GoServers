package common

type Result struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func NewResult(errorCode int, errorMessage string) *Result {
	return &Result{ErrorCode: errorCode, ErrorMessage: errorMessage}
}
