package gate

import (
	"github.com/phuhao00/shine/pkg/chanrpc"
	"github.com/phuhao00/shine/pkg/log"
	network2 "github.com/phuhao00/shine/pkg/network"
	"net"
	"reflect"

	"google.golang.org/protobuf/proto"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network2.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	//HTTPTimeout time.Duration
	//CertFile    string
	//KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
	//
	NewAgent func(*network2.TCPConn) network2.Agent
}

func (gate *Gate) Run(closeSig chan bool) {

	var tcpServer *network2.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network2.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = gate.NewAgent
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network2.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msgId, err := a.gate.Processor.GetMsgID(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msgId, data)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call(0, []byte(""))
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msgId uint16, msg proto.Message) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msgId, msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
