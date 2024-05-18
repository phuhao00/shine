package player

import (
	"github.com/phuhao00/shine/network"
)

type Player struct {
	// ID
	ID uint64
	// Name
	Name string
	// Level
	Level uint32
	// Experience
	Experience uint32
	// Gold
	Gold uint32
	// Diamond
	Diamond uint32
	// Energy
	Energy uint32
	// Stamina
	Stamina uint32
	// LastLoginTime
	LastLoginTime uint32
	// LastLogoutTime
	LastLogoutTime uint32
	// LastSaveTime
	LastSaveTime uint32
	// LastLoadTime
	LastLoadTime uint32
	// LastLogoutTime
	conn *network.TCPConn
}

func NewPlayer() *Player {
	return &Player{}
}

func NewPlayerWithConn(conn *network.TCPConn) network.Agent {
	return &Player{conn: conn}
}

func (p *Player) SetConn(conn *network.TCPConn) {
	p.conn = conn
}

func (p *Player) GetConn() *network.TCPConn {
	return p.conn
}

func (p *Player) Run() {
	// implementation for Run method
}

func (p *Player) OnClose() {
	// implementation for OnClose method
}
