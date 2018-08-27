package inner

type Account struct {
	Id       int64
	Account  string `xorm:"varchar(128) unique"`
	Password string `xorm:"varchar(128)"`
}

// set id
func (a *Account) SetId(val int64) {
	a.Id = val
}

// db id
func (a *Account) DBId() int64 {
	return a.Id
}
