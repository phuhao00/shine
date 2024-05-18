package player

//维护player的管理器
type PlayerManager struct {
	players []*Player
}

func NewManager() *PlayerManager {
	return &PlayerManager{
		players: make([]*Player, 0),
	}
}

func (m *PlayerManager) AddPlayer(player *Player) {
	m.players = append(m.players, player)
}

func (m *PlayerManager) RemovePlayer(player *Player) {
	for i, p := range m.players {
		if p == player {
			m.players = append(m.players[:i], m.players[i+1:]...)
			break
		}
	}
}

func (m *PlayerManager) GetPlayers() []*Player {
	return m.players
}
