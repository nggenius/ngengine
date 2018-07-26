package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"

	"github.com/davecgh/go-spew/spew"
	"github.com/mysll/toolkit"
)

type LoginResult struct {
	Result string
}

type Client struct {
	conn    net.Conn
	account string
	buff    []byte
}

func NewClient() *Client {
	c := &Client{}
	c.buff = make([]byte, 16*1024)
	return c
}

func (c *Client) Connect(addr string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic("connect server failed")
	}
	c.conn = conn
}

func (c *Client) SendMessage(method string, msg interface{}) error {
	r := &c2s.Rpc{}
	r.Node = "."
	r.ServiceMethod = method
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	r.Data = data

	rpcdata, err := json.Marshal(r)
	if err != nil {
		return err
	}

	out, err := protocol.CreateMsg(nil, rpcdata, protocol.C2S_RPC)

	if err != nil {
		return err
	}

	c.conn.Write(out)

	return nil
}

func (c *Client) RecvMessage() {
	id, data, err := protocol.ReadPkg(c.conn, c.buff)
	if err != nil {
		panic(err)
	}

	log.Println("recv message:", id)
	switch id {
	case protocol.S2C_RPC:
		c.ParseRpcMsg(data)
	case protocol.S2C_ERROR:
		var errinfo s2c.Error
		if err := json.Unmarshal(data, &errinfo); err != nil {
			log.Println(err)
			return
		}
		log.Println(errinfo)
	}

}

func (c *Client) ParseRpcMsg(data []byte) {
	var rpc s2c.Rpc
	err := json.Unmarshal(data, &rpc)
	if err != nil {
		log.Println(err)
		return
	}

	switch rpc.Servicemethod {
	case "Login.Nest":
		c.OnNestInfo(rpc.Data)
	case "Account.Roles":
		c.OnRoleInfo(rpc.Data)
	case "system.Error":
		c.OnError(rpc.Data)
	}
}

func (c *Client) OnNestInfo(data []byte) {
	var nest s2c.NestInfo
	if err := json.Unmarshal(data, &nest); err != nil {
		log.Println(err)
		return
	}

	c.conn.Close()
	c.Connect(nest.Addr, int(nest.Port))
	login := c2s.LoginNest{}
	login.Account = c.account
	login.Token = nest.Token
	c.SendMessage("Self.Login", &login)
	c.RecvMessage()
}

func (c *Client) OnRoleInfo(data []byte) {
	var role s2c.RoleInfo
	if err := json.Unmarshal(data, &role); err != nil {
		log.Println(err)
	}

	fmt.Println(spew.Sdump(role))

re:
	for {
		fmt.Println("please input command(create/delete/choose):")
		var cmd string
		fmt.Scanln(&cmd)
	L:
		for {
			switch cmd {
			case "create":
				fmt.Println("please input index name")
				var index int
				var name string
				fmt.Scanln(&index, &name)
				create := c2s.CreateRole{}
				create.Index = index
				create.Name = name
				c.SendMessage("Self.CreateRole", &create)
				c.RecvMessage()
				break re
			case "delete":
				if len(role.Roles) == 0 {
					fmt.Println("no roles")
					continue re
				}
				fmt.Println("please input delete index")
				var index int
				fmt.Scanln(&index)
				d := c2s.DeleteRole{}
				for k := range role.Roles {
					if role.Roles[k].Index == int8(index) {
						d.RoleId = role.Roles[k].RoleId
						break
					}
				}
				if d.RoleId == 0 {
					fmt.Println("index error, please retry")
					continue L
				}

				c.SendMessage("Self.DeleteRole", &d)
				c.RecvMessage()
				break re

			case "choose":
				if len(role.Roles) == 0 {
					fmt.Println("no roles")
					continue re
				}
				fmt.Println("please input choose index")
				var index int
				fmt.Scanln(&index)
				choose := c2s.ChooseRole{}
				for k := range role.Roles {
					if role.Roles[k].Index == int8(index) {
						choose.RoleID = role.Roles[k].RoleId
						break
					}
				}
				if choose.RoleID == 0 {
					fmt.Println("index error, please retry")
					continue L
				}

				c.SendMessage("Self.ChooseRole", &choose)
				c.RecvMessage()

				break re
			default:
				fmt.Println("command error")
				continue re
			}
		}
	}
}

func (c *Client) OnError(data []byte) {
	var errcode s2c.Error
	if err := json.Unmarshal(data, &errcode); err != nil {
		log.Println(err)
		return
	}

	log.Println("error:", errcode)
}
func (c *Client) Login(name, pass string) {
	l := c2s.Login{}
	l.Name = name
	l.Pass = pass
	c.account = name
	c.SendMessage("Account.Login", &l)
	c.RecvMessage()
}

func main() {
	client := NewClient()
	client.Connect("127.0.0.1", 4000)
	client.Login("test", "123")
	toolkit.WaitForQuit()
}
