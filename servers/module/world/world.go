package world

// World is the top-level interface of the game server.

import (
	"github.com/phuhao00/shine/pkg/gate"
	network2 "github.com/phuhao00/shine/pkg/network"
	"github.com/phuhao00/shine/pkg/network/protobuf"
	"github.com/phuhao00/shine/servers/module/family"
	player2 "github.com/phuhao00/shine/servers/module/player"
	"sync"
)

// World is the top-level interface of the game server.

type World struct {
	family         *family.Family
	manager        *player2.PlayerManager
	clients        []*network2.TCPClient
	mutex          sync.Mutex
	SeverForClient *gate.Gate          // Server for client
	ServerForRank  *gate.Gate          // Server for rank
	ServerForGm    *gate.Gate          // Server for gm
	processor      *protobuf.Processor // Processor for message
}

func NewWorld() *World {
	return &World{
		family:    family.NewFamily(),
		manager:   player2.NewManager(),
		processor: protobuf.NewProcessor(),
	}
}

func (w *World) NewServer(newAgent func(*network2.TCPConn) network2.Agent) *gate.Gate {
	return &gate.Gate{
		MaxConnNum:      10000,
		PendingWriteNum: 1000,
		MaxMsgLen:       4096,
		Processor:       w.processor,
		AgentChanRPC:    nil,
		TCPAddr:         "",
		LenMsgLen:       2,
		LittleEndian:    false,
		NewAgent:        newAgent,
	}
}

func (w *World) Run() {
	// Create server for gm
	// Create server for rank
	// Create server for client
	w.SeverForClient = w.NewServer(player2.NewPlayerWithConn)

	gmCloseSig := make(chan bool)
	rankCloseSig := make(chan bool)
	clientCloseSig := make(chan bool)

	w.ServerForGm.Run(gmCloseSig)
	w.ServerForRank.Run(rankCloseSig)
	w.SeverForClient.Run(clientCloseSig)

}

func (w *World) SetServerForGm(server *gate.Gate) {
	w.ServerForGm = server
}

func (w *World) SetServerForRank(server *gate.Gate) {
	w.ServerForRank = server
}

func (w *World) SetServerForClient(server *gate.Gate) {
	w.SeverForClient = server
}

// AddClient adds a new TCPClient to the World.
func (w *World) AddClient(client *network2.TCPClient) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.clients = append(w.clients, client)
}

// RemoveClient removes a TCPClient from the World.
func (w *World) RemoveClient(client *network2.TCPClient) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i, c := range w.clients {
		if c == client {
			w.clients = append(w.clients[:i], w.clients[i+1:]...)
			break
		}
	}
}

// ForwardMessage forwards the message received from the client or other nodes.
func (w *World) ForwardMessage(message interface{}) {
	// TODO: Implement message forwarding logic here
}
