package game

import (
	"errors"
)

const (
	max_teams  = 2
	min_teams  = 2
	map_size_x = 10
	map_size_y = 10
)

type Game interface {
	AddNewPlayer(connID string, playerName string) (Player, error)
	SetPlayerInfo(player Player) error
	GetTeams() []*Team
	GetPlayers() []*Player
	GetGameMap() GameMap
	IsAllUsersReady() bool
	IsFull() bool
}

type game struct {
	teams   []*Team
	players []*Player
	// currentRound CurrentRound
	gameMap GameMap
}

func New() Game {
	teams := make([]*Team, max_teams)
	for i := range max_teams {
		team := newTeam()
		teams[i] = &team
	}

	gameMap := newGameMap(map_size_x, map_size_y)

	return &game{
		teams:   teams,
		gameMap: gameMap,
	}
}

func (g *game) GetTeams() []*Team {
	return g.teams
}

func (g *game) getTeamsWithAvailableSpace() []*Team {
	teamsWithAvailableSpace := make([]*Team, 0, max_teams)

	for _, team := range g.teams {
		if !team.isFull() {
			teamsWithAvailableSpace = append(teamsWithAvailableSpace, team)
		}
	}

	return teamsWithAvailableSpace
}

func (g *game) getTeamsWithPlayers() []*Team {
	teamsWithPlayers := make([]*Team, 0, max_teams)

	for _, team := range g.teams {
		if len(team.Players) > 0 {
			teamsWithPlayers = append(teamsWithPlayers, team)
		}
	}

	return teamsWithPlayers
}

func (g *game) GetPlayers() []*Player {
	return g.players
}

func (g *game) AddNewPlayer(connID string, playerName string) (Player, error) {
	teamsWithAvailableSpace := g.getTeamsWithAvailableSpace()

	if len(teamsWithAvailableSpace) == 0 {
		return Player{}, errors.New("no available teams for new player")
	}

	newPlayer := newPlayer(connID, playerName)
	for range max_ships_per_player {
		shipPosition := newRandomCoordinate(g.gameMap.Size[0], g.gameMap.Size[1])
		ship := newShip(shipPosition)

		newPlayer.AddShip(&ship)
		g.gameMap.AddShip(&ship)
	}

	g.addPlayerToGameAndTeam(teamsWithAvailableSpace[0], &newPlayer)
	return newPlayer, nil
}

func (g *game) addPlayerToGameAndTeam(team *Team, player *Player) {
	g.players = append(g.players, player)
	team.AddPlayer(player)
}

func (g *game) SetPlayerInfo(player Player) error {
	foundPlayer, err := g.findPlayerById(player.ID)
	if err != nil {
		return err
	}

	foundPlayer.Name = player.Name
	foundPlayer.IsReady = player.IsReady

	return nil
}

func (g *game) findPlayerById(playerId string) (*Player, error) {
	for _, player := range g.players {
		if player.ID == playerId {
			return player, nil
		}
	}

	return nil, errors.New("Player not found")
}

func (g *game) GetGameMap() GameMap {
	return g.gameMap
}

func (g *game) IsAllUsersReady() bool {
	teams := g.getTeamsWithPlayers()

	if len(teams) < min_teams {
		return false
	}

	for _, team := range teams {
		for _, player := range team.Players {
			if !player.IsReady {
				return false
			}
		}
	}
	return true
}

func (g *game) IsFull() bool {
	teamsWithAvailableSpace := g.getTeamsWithAvailableSpace()
	return len(teamsWithAvailableSpace) == 0
}
