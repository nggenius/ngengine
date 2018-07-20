package rpc

type Error struct {
	ErrCode int32
	Err     string
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int32, err string) *Error {
	return &Error{code, err}
}
