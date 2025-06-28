package game

import "github.com/google/uuid"

const max_players_per_team = 1

type Team struct {
	ID      string
	Players []*Player
}

func newTeam() Team {
	return Team{
		ID:      uuid.New().String(),
		Players: make([]*Player, 0, max_players_per_team),
	}
}

func (t *Team) AddPlayer(player *Player) {
	t.Players = append(t.Players, player)
}

func (t *Team) isFull() bool {
	return len(t.Players) >= max_players_per_team
}

func (t *Team) GetShips() []*Ship {
	ships := make([]*Ship, 0)
	for _, player := range t.Players {
		ships = append(ships, player.Ships...)
	}
	return ships
}
