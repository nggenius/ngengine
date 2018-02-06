package models

import "time"

type NxConsole struct {
	Id         int64
	Name       string    `xorm:"char(64)"`
	ServerId   int       `xorm:"int(32)"`
	ServerIp   string    `xorm:"char(64)"`
	LoginTime  time.Time `xorm:"updated"`
	DeleteTime time.Time `xorm:"deleted"`
}

func (c *NxConsole) Insert() error {
	if _, err := XormEngine.Insert(c); err != nil {
		return err
	}
	return nil
}

func (c *NxConsole) Delete() error {
	if _, err := XormEngine.Delete(c); err != nil {
		return err
	}
	return nil
}

func (c *NxConsole) Read() error {
	if _, err := XormEngine.Get(c); err != nil {
		return err
	}
	return nil
}

func (c *NxConsole) Update() error {
	if _, err := XormEngine.Update(c); err != nil {
		return err
	}
	return nil
}
