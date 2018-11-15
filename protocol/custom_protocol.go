package protocol

var (
	codec ProtoCodec
)

type ProtoCodec interface {
	GetCodecInfo() string
	CreateRpcMessage(svr, method string, args interface{}) (data []byte, err error)
	DecodeRpcMessage(msg *Message) (node, Servicemethod string, data []byte, err error)
	DecodeMessage(msg *Message, out interface{}) error
}

type Decoder interface {
	DecodeRpcMessage(msg *Message) (node, Servicemethod string, data []byte, err error)
	DecodeMessage(msg *Message, out interface{}) error
}

func RegisterProtoCodec(p ProtoCodec) {
	codec = p
}
