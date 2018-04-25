package inner

type Account struct {
	Id       int64
	Account  string `xorm:"varchar(128) unique"`
	Password string `xorm:"varchar(128)"`
}

type AccountCreater struct {
}

func (c *AccountCreater) Create() interface{} {
	return &Account{}
}

func (c *AccountCreater) CreateSlice() interface{} {
	return &[]*Account{}
}
