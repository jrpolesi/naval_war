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
			err := closeConnection(poolOfConnections, conn, connID)
			if err != nil {
				fmt.Println("Error closing connection:", err)
			}
			fmt.Println("WebSocket connection closed")
		}()

		err = gameMaster.AddNewPlayerToGame(connID)
		if err != nil {
			fmt.Println("Error adding new player to game:", err)
			closeErr := closeConnection(poolOfConnections, conn, connID)

			if closeErr != nil {
				fmt.Println("Error closing connection after adding new player:", closeErr)
			}
			return
		}

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

func closeConnection(pool connection.Pool, conn *websocket.Conn, connID string) error {
	pool.DeleteConnection(connID)

	return conn.Close()
}
