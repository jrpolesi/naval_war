package gamemaster

import (
	"errors"
	"fmt"

	"github.com/jrpolesi/naval_war/cmd/internal/connection"
	"github.com/jrpolesi/naval_war/cmd/internal/events"
	"github.com/jrpolesi/naval_war/cmd/internal/game"
)

type GameMaster interface {
	AddNewPlayerToGame(connID string) error
	UpdatePlayerInfo(connID string, playerInfo UpdatePlayerIntent) error
	PerformPlayerAction() error
}

type gameMaster struct {
	connections connection.Pool
	game        game.Game
}

func New(connections connection.Pool, game game.Game) GameMaster {
	return &gameMaster{
		connections: connections,
		game:        game,
	}
}

func (gm *gameMaster) AddNewPlayerToGame(connID string) error {
	if gm.game.IsFull() {
		return errors.New("game is full")
	}

	player, err := gm.game.AddNewPlayer(connID, "")
	if err != nil {
		return fmt.Errorf("failed to add new player: %w", err)
	}

	newPlayerPayload := Player{
		ID:      player.ID,
		Name:    player.Name,
		IsReady: player.IsReady,
	}
	gm.connections.SendMessage(connID,
		events.NewMessage(events.Server.CreatedNewPlayer, newPlayerPayload))

	return nil
}

func (gm *gameMaster) UpdatePlayerInfo(connID string, playerInfo UpdatePlayerIntent) error {
	gm.game.SetPlayerInfo(game.Player{
		ID:      connID,
		Name:    playerInfo.Name,
		IsReady: playerInfo.IsReady,
	})

	teams := gm.game.GetTeams()
	gm.notifyUsersUpdate(teams)

	if !gm.game.IsAllUsersReady() {
		return nil
	}

	gm.notifyGameReady(connID)
	return nil
}

func (gm *gameMaster) PerformPlayerAction() error {
	return errors.New("not implemented")
}

func (gm *gameMaster) notifyUsersUpdate(updatedTeams []*game.Team) {
	teamsResponse := make([]Team, len(updatedTeams))

	for teamIndex, team := range updatedTeams {

		players := make([]Player, len(team.Players))

		for playerIndex, player := range team.Players {
			players[playerIndex] = Player{
				ID:      player.ID,
				Name:    player.Name,
				IsReady: player.IsReady,
			}
		}

		teamsResponse[teamIndex] = Team{
			ID:      team.ID,
			Players: players,
		}
	}

	gm.connections.SendMessageToAll(
		events.NewMessage(events.Server.UpdatedPlayersInfo, GameResponse{
			Teams: &teamsResponse,
		}),
	)
}

func (gm *gameMaster) notifyGameReady(connId string) {
	gm.connections.SendMessageToAll(
		events.NewMessage(
			events.Server.GameStarted,
			gm.getGameState(connId),
		),
	)
}

func (gm *gameMaster) getGameState(currPlayerID string) GameResponse {
	teams := gm.game.GetTeams()
	teamsResponse := make([]Team, len(teams))

	for teamIndex, team := range teams {
		players := make([]Player, len(team.Players))

		for playerIndex, player := range team.Players {
			players[playerIndex] = Player{
				ID:      player.ID,
				Name:    player.Name,
				IsReady: player.IsReady,
			}
		}

		teamsResponse[teamIndex] = Team{
			ID:      team.ID,
			Players: players,
		}
	}

	players := gm.game.GetPlayers()
	shipsResponse := []Ship{}
	for _, player := range players {
		if player.ID != currPlayerID {
			continue
		}

		for _, ship := range player.Ships {
			var damagedBy *Player

			if ship.DamagedBy != nil {
				damagedBy = &Player{
					ID:      ship.DamagedBy.ID,
					Name:    ship.DamagedBy.Name,
					IsReady: ship.DamagedBy.IsReady,
				}
			}

			shipsResponse = append(shipsResponse, Ship{
				ID:        ship.ID,
				Position:  Coordinate(ship.Position),
				IsDamaged: ship.DamagedBy != nil,
				DamagedBy: damagedBy,
				IsOwner:   player.ID == currPlayerID,
			})
		}
	}

	gameMap := gm.game.GetGameMap()
	var gameMapResponse *GameMap
	if len(gameMap.Size) == 2 {
		gameMapResponse = &GameMap{
			Size: gameMap.Size,
		}
	}

	return GameResponse{
		Map:   gameMapResponse,
		Teams: &teamsResponse,
		Ships: &shipsResponse,
	}
}
