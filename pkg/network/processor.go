package network

import "google.golang.org/protobuf/proto"

type Processor interface {
	// must goroutine safe
	Route(msgId uint16, data []byte) error
	// must goroutine safe
	GetMsgID(data []byte) (msgId uint16, err error)
	// must goroutine safe
	Unmarshal(data []byte) (proto.Message, error)
	// must goroutine safe
	Marshal(msgId uint16, msg proto.Message) ([]byte, error)
}
