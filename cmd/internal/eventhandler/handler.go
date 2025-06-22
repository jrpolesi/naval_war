package eventhandler

import (
	"encoding/json"
	"log"

	"github.com/jrpolesi/naval_war/cmd/internal/connection"
	"github.com/jrpolesi/naval_war/cmd/internal/events"
	"github.com/jrpolesi/naval_war/cmd/internal/gamemaster"
)

type eventHandler struct {
	poolOfConnections connection.Pool
	gameMaster        gamemaster.GameMaster
}

type EventHandler interface {
	HandleEvent(connID string, msg events.Message) error
}

func New(poolOfConnections connection.Pool, gameMaster gamemaster.GameMaster) EventHandler {
	return &eventHandler{
		poolOfConnections: poolOfConnections,
		gameMaster:        gameMaster,
	}
}

func (h eventHandler) HandleEvent(connID string, msg events.Message) error {
	log.Println("Received event:", msg.Event)
	log.Println("Payload:", msg.Payload)

	switch msg.Event {
	case events.Client.UpdatedPlayerInfo:
		return h.updatePlayerInfo(connID, msg)
	default:
		return h.unknownEvent(connID, msg)
	}
}

func (h eventHandler) unknownEvent(connID string, msg events.Message) error {
	log.Printf("Unknown event: %s", msg.Event)

	event := events.Message{
		Event: events.Server.UnknownEvent,
		Payload: map[string]string{
			"error": "Unknown event",
		},
	}

	h.poolOfConnections.SendMessage(connID, event)

	return nil
}

func (h eventHandler) updatePlayerInfo(connID string, msg events.Message) error {
	payloadBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return nil
	}

	var playerIntent gamemaster.UpdatePlayerIntent
	if err := json.Unmarshal(payloadBytes, &playerIntent); err != nil {
		log.Printf("Error unmarshalling player intent: %v", err)
		return nil
	}

	playerIntent.ID = connID

	return h.gameMaster.UpdatePlayerInfo(connID, playerIntent)
}
