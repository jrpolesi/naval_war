package game

import "github.com/google/uuid"

type Ship struct {
	ID        string
	Position  Coordinate
	DamagedBy *Player
}

func newShip(position Coordinate) Ship {
	return Ship{
		ID:       uuid.New().String(),
		Position: position,
	}
}
