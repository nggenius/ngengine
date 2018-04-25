package main

import (
	"encoding/json"
	"fmt"
	"net"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"

	"github.com/mysll/toolkit"
)

type LoginResult struct {
	Result string
}

type Client struct {
	conn net.Conn
	buff []byte
}

func NewClient(addr string, port int) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic("connect server failed")
	}

	c := &Client{}
	c.conn = conn
	c.buff = make([]byte, 16*1024)
	return c
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
}

func (c *Client) Login(name, pass string) {
	l := c2s.Login{}
	l.Name = name
	l.Pass = pass
	c.SendMessage("Account.Login", &l)
	c.RecvMessage()
}

func main() {
	client := NewClient("127.0.0.1", 2002)

	client.Login("test", "123")
	toolkit.WaitForQuit()
}
