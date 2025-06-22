package events

type Event string

type clientEvents struct {
	UpdatedPlayerInfo     Event
	PlayerPerformedAction Event
}

var Client = clientEvents{
	UpdatedPlayerInfo:     "client_updated_player_info",
	PlayerPerformedAction: "client_player_performed_action",
}

type serverEvents struct {
	UnknownEvent       Event
	CreatedNewPlayer   Event
	UpdatedPlayersInfo Event
	GameStarted        Event
	GameUpdated        Event
	GameOver           Event
}

var Server = serverEvents{
	UnknownEvent:       "server_unknown_event",
	CreatedNewPlayer:   "server_created_new_player",
	UpdatedPlayersInfo: "server_updated_players_info",
	GameStarted:        "server_game_started",
	GameUpdated:        "server_game_updated",
	GameOver:           "server_game_over",
}

type Message struct {
	Event   Event `json:"event"`
	Payload any   `json:"payload"`
}

func NewMessage(event Event, payload any) Message {
	return Message{
		Event:   event,
		Payload: payload,
	}
}
