package gamemaster

import (
	"errors"
	"log"

	"github.com/jrpolesi/naval_war/cmd/internal/connection"
	"github.com/jrpolesi/naval_war/cmd/internal/events"
	"github.com/jrpolesi/naval_war/cmd/internal/game"
)

type GameMaster interface {
	UpdatePlayerInfo(playerInfo UpdatePlayerIntent) error
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

func (gm *gameMaster) UpdatePlayerInfo(playerInfo UpdatePlayerIntent) error {
	gm.game.SetPlayerInfo(game.Player{
		ID:      playerInfo.ID,
		Name:    playerInfo.Name,
		IsReady: playerInfo.IsReady,
	})

	teams := gm.game.GetTeams()

	gm.notifyUsersUpdate(teams)

	if !gm.game.IsAllUsersReady() {
		return nil
	}

	return nil
}

func (gm *gameMaster) PerformPlayerAction() error {
	return errors.New("not implemented")
}

func (gm *gameMaster) notifyUsersUpdate(updatedTeams []game.Team) {
	log.Printf("Notifying users about updated teams: %+v", updatedTeams)
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
		events.Message{
			Event: events.Server.UpdatedPlayersInfo,
			Payload: GameResponse{
				Teams: &teamsResponse,
			},
		})
}
