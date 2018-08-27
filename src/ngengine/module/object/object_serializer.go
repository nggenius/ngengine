package object

type Container interface {
	Cap() int
	ChildAtIf(pos int) interface{}
}

func Encode(o interface{}) ([]byte, error) {

	return nil, nil
}

func Decode(b []byte) (interface{}, error) {
	return nil, nil
}
