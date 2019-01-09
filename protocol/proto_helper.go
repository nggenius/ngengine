package protocol

import (
	"github.com/nggenius/ngengine/logger"
	"github.com/nggenius/ngengine/share"
	"github.com/nggenius/ngengine/utils"
)

type BodyWriter struct {
	*utils.StoreArchive
	msg *Message
}

func (w *BodyWriter) GetMessage() *Message {
	return w.msg
}

func (w *BodyWriter) Flush() {
	w.msg.Body = w.msg.Body[:w.Len()]
}

func (w *BodyWriter) Free() {
	w.msg.Free()
}

type HeadWriter struct {
	*utils.StoreArchive
	msg *Message
}

func (w *HeadWriter) Flush() {
	w.msg.Header = w.msg.Header[:w.Len()]
}

func NewMessageReader(msg *Message) *utils.LoadArchive {
	return utils.NewLoadArchiver(msg.Body)
}

func ParseReply(msg *Message) (int32, *utils.LoadArchive) {
	if msg == nil {
		return 0, nil
	}
	errcode := GetReplyError(msg)
	ar := utils.NewLoadArchiver(msg.Body)
	return errcode, ar
}

func ParseErrMsg(msg *Message) (int32, *Message) {
	if msg == nil {
		return 0, nil
	}

	errcode := GetReplyError(msg)

	return errcode, msg
}

func NewMessageWriter(msg *Message) *BodyWriter {
	msg.Body = msg.Body[:0]
	w := &BodyWriter{utils.NewStoreArchiver(msg.Body), msg}
	return w
}

func NewProtoMessage() *BodyWriter {
	msg := NewMessage(share.MAX_BUF_LEN)
	w := &BodyWriter{utils.NewStoreArchiver(msg.Body), msg}
	return w
}

func NewHeadReader(msg *Message) *utils.LoadArchive {
	return utils.NewLoadArchiver(msg.Header)
}

func NewHeadWriter(msg *Message) *HeadWriter {
	msg.Header = msg.Header[:0]
	w := &HeadWriter{utils.NewStoreArchiver(msg.Header), msg}
	return w
}

//获取rpc消息的错误代码，返回0没有错误
func GetReplyError(msg *Message) int32 {
	ar := utils.NewLoadArchiver(msg.Header)
	if len(msg.Header) <= 8 {
		return 0
	}
	ar.Seek(8, 0)

	haserror, err := ar.GetInt8()
	if err != nil {
		return 0
	}

	if haserror != 1 {
		return 0
	}

	errcode, err := ar.GetInt32()
	if err != nil {
		return 0
	}

	return errcode
}

// 解析消息协议，由使用的消息类型决定解码方式。用于客户端与服务器的通讯
func ParseProto(decoder Decoder, msg *Message, obj interface{}) error {
	return decoder.DecodeMessage(msg, obj)
}

// 解析参数，用于服务器之间的通讯
func ParseArgs(msg *Message, args ...interface{}) error {
	if len(args) == 0 || msg == nil {
		return nil
	}

	ar := NewMessageReader(msg)
	for i := 0; i < len(args); i++ {
		err := ar.Get(args[i])
		if err != nil {
			return err
		}
	}

	return nil
}

type errcoder interface {
	ErrCode() int32
}

// CheckRpcError 检查rpc的系统错误，不检查逻辑错误
func CheckRpcError(c errcoder) bool {
	code := int(c.ErrCode())
	if code == 0 {
		return false
	}

	if code == share.ERR_ARGS_ERROR ||
		code == share.ERR_TIME_OUT ||
		code == share.ERR_RPC_FAILED {
		return true
	}

	return false
}

// 错误处理函数，如果有错误则写入日志
func Check(l *logger.Log, err error) bool {
	if err != nil {
		l.Output(3, "[W]", err)
		return true
	}

	return false
}

// 错误处理函数，如果有错误则写入日志
func Check2(l *logger.Log, _ interface{}, err error) bool {
	if err != nil {
		l.Output(3, "[W]", err)
		return true
	}

	return false
}

// 按传入参数，序列化消息
func CreateMessage(args ...interface{}) *Message {
	if len(args) > 0 {
		msg := NewProtoMessage()
		for i := 0; i < len(args); i++ {
			err := msg.Put(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed," + err.Error())
			}
		}
		msg.Flush()
		return msg.GetMessage()
	}
	return nil
}

// 向已有的message中附加参数
func AppendArgs(msg *Message, args ...interface{}) {
	w := &BodyWriter{utils.NewStoreArchiver(msg.Body), msg}
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			err := w.Put(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed," + err.Error())
			}
		}
		w.Flush()
	}
}

// 缓冲区大小定义
const (
	DEF    = -1
	TINY   = 256
	MIDDLE = 1024
	BIG    = 4096
)

// rpc返回值, poolsize为使用的缓冲区大小，请根据实际大小进行预处理
func Reply(poolsize int, args ...interface{}) (int32, *Message) {
	if poolsize == DEF {
		poolsize = share.MAX_BUF_LEN
	}

	if poolsize > share.MAX_BUF_LEN {
		poolsize = share.MAX_BUF_LEN
	}

	msg := NewMessage(poolsize)
	if len(args) > 0 {
		mw := NewMessageWriter(msg)
		for i := 0; i < len(args); i++ {
			err := mw.Put(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed," + err.Error())
			}
		}
		mw.Flush()
	}

	return 0, msg
}

// rpc返回值, poolsize为使用的缓冲区大小，请根据实际大小进行预处理
func ReplyError(poolsize int, errcode int32, err string, args ...interface{}) (int32, *Message) {
	if poolsize == DEF {
		poolsize = share.MAX_BUF_LEN
	}

	if poolsize > share.MAX_BUF_LEN {
		poolsize = share.MAX_BUF_LEN
	}

	msg := NewMessage(poolsize)
	mw := NewMessageWriter(msg)
	if err := mw.Put(err); err != nil {
		msg.Free()
		panic("write error failed, " + err.Error())
	}

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			err := mw.Put(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed," + err.Error())
			}
		}
	}

	mw.Flush()

	return errcode, msg
}
