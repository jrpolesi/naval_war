package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jrpolesi/naval_war/cmd/internal/connection"
	"github.com/jrpolesi/naval_war/cmd/internal/eventhandler"
	"github.com/jrpolesi/naval_war/cmd/internal/events"
	"github.com/jrpolesi/naval_war/cmd/internal/game"
	"github.com/jrpolesi/naval_war/cmd/internal/gamemaster"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHandler() http.Handler {
	poolOfConnections := connection.NewPool()

	newGame := game.New()
	gameMaster := gamemaster.New(poolOfConnections, newGame)
	eventHandler := eventhandler.New(poolOfConnections, gameMaster)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		connID := poolOfConnections.AddConnection(conn)

		defer func() {
			poolOfConnections.DeleteConnection(connID)
			conn.Close()
			fmt.Println("WebSocket connection closed")
		}()

		fmt.Println("WebSocket connection established")

		for {
			var msg events.Message
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("Error reading JSON:", err)
				return
			}

			eventHandler.HandleEvent(connID, msg)
		}
	})
}
