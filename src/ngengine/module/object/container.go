package object

type Container struct {
	Cap   int
	Child []Object `json:"-"`
}

func NewContainer() *Container {
	c := &Container{}
	return c
}
