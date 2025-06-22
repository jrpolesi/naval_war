package game

import (
	"errors"
)

type CurrentRound struct {
	Actions []Action
}

type Action struct {
	Type    string
	Payload ActionPayload
}

type ActionPayload struct {
	Position Coordinate
}

const max_ships_per_player = 3

type Player struct {
	ID      string
	Name    string
	IsReady bool
	Ships   []*Ship
}

func newPlayer(connID string, name string) Player {
	return Player{
		ID:      connID,
		Name:    name,
		IsReady: false,
		Ships:   make([]*Ship, 0, max_ships_per_player),
	}
}

func (p *Player) AddShip(ship *Ship) error {
	if len(p.Ships) >= max_ships_per_player {
		return errors.New("maximum number of ships reached for player")
	}

	p.Ships = append(p.Ships, ship)

	return nil
}
