package client_core

type CPlayerManager struct {
	players map[uint64]*CPlayer
	myId    uint64
}

func CreateCPlayerManager() *CPlayerManager {
	return &CPlayerManager{
		players: make(map[uint64]*CPlayer, 10),
	}
}

func (pm *CPlayerManager) Get(playerId uint64) *CPlayer {
	return pm.players[playerId]
}

func (pm *CPlayerManager) Add(player *CPlayer) {
	_, o := pm.players[player.Id()]
	if !o {
		pm.players[player.Id()] = player
	}
}

func (pm *CPlayerManager) Remove(playerId uint64) bool {
	_, o := pm.players[playerId]
	if !o {
		return false
	}
	delete(pm.players, playerId)
	return true
}

func (pm *CPlayerManager) AddMe(player *CPlayer) {
	pm.Add(player)
	pm.myId = player.Id()
}

func (pm *CPlayerManager) GetMe() *CPlayer {
	return pm.Get(pm.myId)
}
