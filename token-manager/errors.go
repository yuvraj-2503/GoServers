package token_manager

type InvalidTokenError struct {
}

type TokenExpiryError struct {
}

type TokenDecodeError struct {
}

func (t *TokenExpiryError) Error() string {
	return "token is expired"
}

func (e *InvalidTokenError) Error() string {
	return "token is not valid"
}

func (e *TokenDecodeError) Error() string {
	return "failed to decode token"
}
