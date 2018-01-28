package proto

import (
	"encoding/json"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"
	"ngengine/utils"
)

type JsonProto struct {
}

func (j *JsonProto) GetCodecInfo() string {
	return "json"
}

func (j *JsonProto) CreateRpcMessage(svr, method string, args interface{}) (data []byte, err error) {
	r := &s2c.Rpc{}
	r.Sender = svr
	r.Servicemethod = method
	if r.Data, err = json.Marshal(args); err != nil {
		return
	}
	data, err = json.Marshal(r)
	return
}

func (j *JsonProto) DecodeRpcMessage(msg *protocol.Message) (node, Servicemethod string, data []byte, err error) {
	request := &c2s.Rpc{}

	if err = json.Unmarshal(msg.Body, request); err != nil {
		return "", "", nil, err
	}

	return request.Node, request.ServiceMethod, request.Data, nil
}

func (j *JsonProto) DecodeMessage(msg *protocol.Message, out interface{}) error {
	r := utils.NewLoadArchiver(msg.Body)
	data, err := r.ReadData()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
