package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/phuhao00/shine/chanrpc"
	"github.com/phuhao00/shine/log"
	"google.golang.org/protobuf/proto"
)

// -------------------------
// msgLen | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      map[uint16]*MsgInfo
}

type MsgInfo struct {
	msgType    reflect.Type
	msgRouter  *chanrpc.Server
	msgHandler MsgHandler
}

type MsgHandler func([]byte)

type MsgRaw struct {
	msgID      uint16
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msgId uint16, msg proto.Message) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if uint64(len(p.msgInfo)) >= math.MaxUint64 {
		log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgId] = i
	id := uint16(len(p.msgInfo) - 1)
	return id
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msgId uint16, msgRouter *chanrpc.Server) {
	p.msgInfo[msgId].msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msgId uint16, msgHandler MsgHandler) {
	p.msgInfo[msgId].msgHandler = msgHandler
}

// goroutine safe
func (p *Processor) Route(msgId uint16, data []byte) error {

	i := p.msgInfo[msgId]
	if i.msgHandler != nil {
		i.msgHandler(data)
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgId, data)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (proto.Message, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}

	// id
	var id uint16
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}
	if id >= uint16(len(p.msgInfo)) {
		return nil, fmt.Errorf("message id %v not registered", id)
	}

	// msg
	//i := p.msgInfo[uint16(id)]

	return nil, nil

}

// goroutine safe
func (p *Processor) GetMsgID(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}
	return binary.BigEndian.Uint16(data[:2]), nil
}

// goroutine safe
func (p *Processor) Marshal(msgId uint16, msg proto.Message) ([]byte, error) {
	id := make([]byte, 8)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, msgId)
	} else {
		binary.BigEndian.PutUint16(id, msgId)
	}
	// data
	data, err := proto.Marshal(msg)
	ret := make([]byte, len(id)+len(data))
	ret = append(ret, id...)
	ret = append(ret, data...)
	return ret, err
}

// goroutine safe
func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}
