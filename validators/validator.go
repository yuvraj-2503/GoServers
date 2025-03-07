package validators

type Validator[T any] interface {
	Validate(T) (bool, error)
}
