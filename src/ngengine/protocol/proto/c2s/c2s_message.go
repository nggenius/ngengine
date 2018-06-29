package c2s

type Rpc struct {
	Node          string
	ServiceMethod string
	Data          []byte
}

type Login struct {
	Name string
	Pass string
}

type LoginNest struct {
	Account string
	Token   string
}

type CreateRole struct {
	Index int
	Name  string // 角色名
}
