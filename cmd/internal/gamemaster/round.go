package gamemaster

import (
	"errors"
	"log"
	"slices"

	gamepkg "github.com/jrpolesi/naval_war/cmd/internal/game"
)

type Round struct {
	CurrentRound        int
	remainingPlayersIDs []string
	Actions             []PerformPlayerActionIntent
}

func NewRound(currentRound int, remainingPlayersIDs []string) Round {
	return Round{
		CurrentRound:        currentRound,
		remainingPlayersIDs: remainingPlayersIDs,
	}
}

func (r *Round) AddAction(action PerformPlayerActionIntent) error {
	if !slices.Contains(r.remainingPlayersIDs, action.PlayerID) {
		return errors.New("this player already performed an action")
	}

	r.Actions = append(r.Actions, action)
	// remainingAmount := len(r.remainingPlayersIDs)

	for i, id := range r.remainingPlayersIDs {
		if id != action.PlayerID {
			continue
		}

		r.remainingPlayersIDs = slices.Delete(r.remainingPlayersIDs, i, i+1)
		// if remainingAmount-1 == i {
		// 	r.remainingPlayersIDs = r.remainingPlayersIDs[:i]
		// } else {
		// 	r.remainingPlayersIDs = slices.Delete(r.remainingPlayersIDs, i, i+1)
		// }

		break
	}

	return nil
}

func (r *Round) IsFinished() bool {
	return len(r.remainingPlayersIDs) == 0
}

func (r *Round) PerformActions(game gamepkg.Game) error {
	for _, action := range r.Actions {

		switch action.ActionType {
		case gamepkg.ActionTypeAttack:
			err := r.performAttack(game, action)
			if err != nil {
				log.Printf("Failed to perform attack action: %v", err)
			}
		default:
			log.Printf("Unknown action type: %s", action.ActionType)
		}
	}

	return nil
}

func (r *Round) performAttack(game gamepkg.Game, action PerformPlayerActionIntent) error {
	targetPosition := gamepkg.NewCoordinate(
		action.Payload.Position.X,
		action.Payload.Position.Y,
	)

	return game.Attack(action.PlayerID, targetPosition)
}
