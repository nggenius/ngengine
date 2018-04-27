package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"

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

	fmt.Println(id, string(data))
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
	case "login.Nest":
		c.OnNestInfo(rpc.Data)
	case "login.Error":
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

func (c *Client) OnError(data []byte) {
	var errcode s2c.Error
	if err := json.Unmarshal(data, &errcode); err != nil {
		log.Println(err)
		return
	}

	log.Println(errcode)
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
	client.Connect("127.0.0.1", 2002)
	client.Login("test", "123")
	toolkit.WaitForQuit()
}
