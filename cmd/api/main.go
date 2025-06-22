package main

import (
	"log"
	"net/http"

	"github.com/jrpolesi/naval_war/cmd/internal/websocket"
)

func main() {
	handler := websocket.NewHandler()

	http.Handle("/websocket", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
