package player

import (
	network2 "github.com/phuhao00/shine/pkg/network"
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
	conn *network2.TCPConn
}

func NewPlayer() *Player {
	return &Player{}
}

func NewPlayerWithConn(conn *network2.TCPConn) network2.Agent {
	return &Player{conn: conn}
}

func (p *Player) SetConn(conn *network2.TCPConn) {
	p.conn = conn
}

func (p *Player) GetConn() *network2.TCPConn {
	return p.conn
}

func (p *Player) Run() {
	// implementation for Run method
}

func (p *Player) OnClose() {
	// implementation for OnClose method
}
