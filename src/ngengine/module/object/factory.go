package object

type factory struct {
	serial uint32
}

func newFactory() *factory {
	f := &factory{}
	return f
}

func (f *factory) create(typ string) *Object {
	return nil
}
