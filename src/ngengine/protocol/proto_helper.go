package protocol

import (
	"ngengine/logger"
	"ngengine/share"
	"ngengine/utils"
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

	haserror, err := ar.ReadInt8()
	if err != nil {
		return 0
	}

	if haserror != 1 {
		return 0
	}

	errcode, err := ar.ReadInt32()
	if err != nil {
		return 0
	}

	return errcode
}

func ParseProto(decoder Decoder, msg *Message, obj interface{}) error {
	return decoder.DecodeMessage(msg, obj)
}

func ParseArgs(msg *Message, args ...interface{}) error {
	if len(args) == 0 || msg == nil {
		return nil
	}

	ar := NewMessageReader(msg)
	for i := 0; i < len(args); i++ {
		err := ar.Read(args[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func Check(l *logger.Log, err error) bool {
	if err != nil {
		l.Output(3, "[W]", err)
		return true
	}

	return false
}

func Check2(l *logger.Log, _ interface{}, err error) bool {
	if err != nil {
		l.Output(3, "[W]", err)
		return true
	}

	return false
}

func CreateMessage(args ...interface{}) *Message {
	if len(args) > 0 {
		msg := NewProtoMessage()
		for i := 0; i < len(args); i++ {
			err := msg.Write(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed")
			}
		}
		msg.Flush()
		return msg.GetMessage()
	}
	return nil
}

func AppendArgs(msg *Message, args ...interface{}) {
	w := &BodyWriter{utils.NewStoreArchiver(msg.Body), msg}
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			err := w.Write(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed")
			}
		}
		w.Flush()
	}
}

//错误消息
func ReplyErrorMessage(errcode int32) *Message {
	msg := NewMessage(1)
	if errcode == 0 {
		return msg
	}
	sr := utils.NewStoreArchiver(msg.Header)
	sr.Write(int8(1))
	sr.Write(errcode)
	msg.Header = msg.Header[:sr.Len()]
	return msg
}

func ReplyErrorAndArgs(errcode int32, args ...interface{}) *Message {
	msg := NewMessage(share.MAX_BUF_LEN)

	if errcode > 0 {
		sr := utils.NewStoreArchiver(msg.Header)
		sr.Write(int8(1))
		sr.Write(errcode)
		msg.Header = msg.Header[:sr.Len()]
	}

	if len(args) > 0 {
		mw := NewMessageWriter(msg)
		for i := 0; i < len(args); i++ {
			err := mw.Write(args[i])
			if err != nil {
				msg.Free()
				panic("write args failed")
				return nil
			}
		}
		mw.Flush()
	}

	return msg
}
