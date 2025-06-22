package game

import (
	"errors"
	"log"
)

type Game interface {
	SetPlayerInfo(player Player) error
	SetShipDamager(shipId, playerId string) error
	GetTeams() []Team
	IsAllUsersReady() bool
	GetShipAtPos(coord Coordinate) (Ship, error)
}

type game struct {
	teams        []Team
	Players      []Player
	currentRound CurrentRound
	Map          Map
}

func New() Game {
	return &game{}
}

func (g *game) findPlayerById(playerId string) (*Player, error) {
	for _, player := range g.Players {
		if player.ID == playerId {
			return &player, nil
		}
	}

	return nil, errors.New("Player not found")
}

func (g *game) findShipById(shipId string) (*Ship, error) {
	for _, ship := range g.Map.Ships {
		if ship.ID == shipId {
			return ship, nil
		}
	}

	return nil, errors.New("Ship not found")
}

func (g game) GetTeams() []Team {
	return g.teams
}

func (g *game) SetPlayerInfo(player Player) error {
	foundPlayer, err := g.findPlayerById(player.ID)
	if err != nil {
		return err
	}

	foundPlayer.Name = player.Name
	foundPlayer.IsReady = player.IsReady

	log.Printf("teams: %+v", g)
	return nil
}

func (g *game) SetShipDamager(shipId, playerId string) error {
	ship, err := g.findShipById(shipId)
	if err != nil {
		return err
	}

	player, err := g.findPlayerById(playerId)
	if err != nil {
		return err
	}

	ship.DamagedBy = player

	return nil
}

func (g *game) GetShipAtPos(coord Coordinate) (Ship, error) {
	// g.
	return Ship{}, nil
}

func (g *game) IsAllUsersReady() bool {
	for _, player := range g.Players {
		if !player.IsReady {
			return false
		}
	}
	return true
}
