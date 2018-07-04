package gameobject

import "ngengine/core/rpc"

type transfer interface {
	// 发起远程调用
	Mailto(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args ...interface{}) error
	// 发起远程调用并调用回调函数
	MailtoAndCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, args ...interface{}) error
}

type Transport struct {
	mailbox  rpc.Mailbox
	transfer transfer
}

func NewTransport(tf transfer, mailbox rpc.Mailbox) *Transport {
	t := &Transport{}
	t.transfer = tf
	t.mailbox = mailbox
	return t
}

// 给自己的客户端发送消息
func (t Transport) Self(method string, args ...interface{}) error {
	return t.transfer.Mailto(&t.mailbox, &t.mailbox, method, args...)
}

// 给指定对象发送消息
func (t Transport) Send(to rpc.Mailbox, method string, args ...interface{}) error {
	return t.transfer.Mailto(&t.mailbox, &to, method, args...)
}
