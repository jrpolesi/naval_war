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
	PerformPlayerAction(connID string, actionIntent PerformPlayerActionIntent) error
}

type gameMaster struct {
	connections   connection.Pool
	game          game.Game
	round         Round
	gameIsRunning bool
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

	gm.gameIsRunning = true
	gm.round = NewRound(0, gm.getPlayersIDs())

	gm.notifyGameReadyToAll()
	return nil
}

func (gm *gameMaster) PerformPlayerAction(connID string, playerActionIntent PerformPlayerActionIntent) error {
	if !gm.gameIsRunning {
		return errors.New("game is not running")
	}

	playerActionIntent.PlayerID = connID

	err := gm.round.AddAction(playerActionIntent)
	if err != nil {
		return fmt.Errorf("failed to add action: %w", err)
	}

	if !gm.round.IsFinished() {
		return nil
	}
	err = gm.round.PerformActions(gm.game)
	if err != nil {
		return fmt.Errorf("failed to perform actions: %w", err)
	}

	gm.notifyGameUpdate()

	gm.round = NewRound(gm.round.CurrentRound+1, gm.getPlayersIDs())

	if gm.game.IsFinished() {
		gm.notifyGameOver()
		gm.gameIsRunning = false
		return nil
	}

	return nil
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

func (gm *gameMaster) notifyGameUpdate() {
	for _, connID := range gm.connections.GetConnectionsIDs() {
		gm.connections.SendMessage(
			connID,
			events.NewMessage(
				events.Server.GameUpdated,
				gm.getGameState(connID),
			),
		)
	}
}

func (gm *gameMaster) notifyGameOver() {
	for _, connID := range gm.connections.GetConnectionsIDs() {
		gm.connections.SendMessage(
			connID,
			events.NewMessage(
				events.Server.GameOver,
				gm.getGameFinishState(connID),
			),
		)
	}
}

func (gm *gameMaster) notifyGameReadyToAll() {
	for _, connID := range gm.connections.GetConnectionsIDs() {
		gm.connections.SendMessage(
			connID,
			events.NewMessage(
				events.Server.GameStarted,
				gm.getGameState(connID),
			),
		)
	}
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
	var damagedShips []Ship
	if len(gameMap.Ships) > 0 {
		damagedShips = make([]Ship, 0, len(gameMap.Ships))
		for _, ship := range gameMap.Ships {
			var damagedBy *Player
			if ship.DamagedBy != nil {
				damagedBy = &Player{
					ID:      ship.DamagedBy.ID,
					Name:    ship.DamagedBy.Name,
					IsReady: ship.DamagedBy.IsReady,
				}
			}

			ship := Ship{
				ID:        ship.ID,
				Position:  Coordinate(ship.Position),
				IsDamaged: ship.DamagedBy != nil,
				DamagedBy: damagedBy,
			}

			if ship.IsDamaged {
				damagedShips = append(damagedShips, ship)
			}
		}
	}

	if len(gameMap.Size) == 2 {
		gameMapResponse = &GameMap{
			Size: gameMap.Size,
		}
	}

	return GameResponse{
		Map:          gameMapResponse,
		Teams:        &teamsResponse,
		Ships:        &shipsResponse,
		DamagedShips: &damagedShips,
	}
}

func (gm *gameMaster) getGameFinishState(currPlayerID string) GameResponse {
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
	var damagedShips []Ship
	if len(gameMap.Ships) > 0 {
		damagedShips = make([]Ship, 0, len(gameMap.Ships))
		for _, ship := range gameMap.Ships {
			var damagedBy *Player
			if ship.DamagedBy != nil {
				damagedBy = &Player{
					ID:      ship.DamagedBy.ID,
					Name:    ship.DamagedBy.Name,
					IsReady: ship.DamagedBy.IsReady,
				}
			}

			ship := Ship{
				ID:        ship.ID,
				Position:  Coordinate(ship.Position),
				IsDamaged: ship.DamagedBy != nil,
				DamagedBy: damagedBy,
			}

			if ship.IsDamaged {
				damagedShips = append(damagedShips, ship)
			}
		}
	}

	if len(gameMap.Size) == 2 {
		gameMapResponse = &GameMap{
			Size: gameMap.Size,
		}
	}

	return GameResponse{
		Map:          gameMapResponse,
		Teams:        &teamsResponse,
		Ships:        &shipsResponse,
		DamagedShips: &damagedShips,
	}
}

func (gm *gameMaster) getPlayersIDs() []string {
	playerIDs := make([]string, 0, len(gm.game.GetPlayers()))
	for _, player := range gm.game.GetPlayers() {
		playerIDs = append(playerIDs, player.ID)
	}
	return playerIDs
}
