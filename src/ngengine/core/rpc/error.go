package rpc

const (
	OK = int32(0)
)

type Error struct {
	Code int32
	Err  string
}

func (e Error) Error() string {
	return e.Err
}

func (e Error) ErrCode() int32 {
	return e.Code
}

func NewError(code int32, err string) *Error {
	return &Error{code, err}
}
