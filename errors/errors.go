package errors

type ErrUserAlreadyExists struct{}

func (e ErrUserAlreadyExists) Error() string {
	return "User already exists"
}
