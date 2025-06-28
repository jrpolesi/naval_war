package connection

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jrpolesi/naval_war/cmd/internal/events"
)

type connections map[string]*websocket.Conn

type pool struct {
	connections connections
	connMutex   sync.Mutex
}

type Pool interface {
	AddConnection(conn *websocket.Conn) string
	DeleteConnection(id string)
	SendMessageToAll(message events.Message)
	SendMessage(id string, message events.Message)
	GetConnectionsIDs() []string 
}

func NewPool() Pool {
	return &pool{
		connections: make(connections),
	}
}

func (p *pool) GetConnectionsIDs() []string {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()

	ids := make([]string, 0, len(p.connections))
	for id := range p.connections {
		ids = append(ids, id)
	}
	return ids
}

func (p *pool) AddConnection(conn *websocket.Conn) string {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()

	id := uuid.New().String()
	p.connections[id] = conn
	return id
}

func (p *pool) DeleteConnection(id string) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()

	delete(p.connections, id)
}

func (p *pool) SendMessageToAll(message events.Message) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()

	for id, conn := range p.connections {
		err := conn.WriteJSON(message)
		if err != nil {
			fmt.Println("Error sending message to connection:", id, err)
			conn.Close()
			delete(p.connections, id)
		}
	}
}

func (p *pool) SendMessage(connID string, message events.Message) {

	conn, exists := p.connections[connID]
	if !exists {
		return
	}

	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Println("Error sending message to connection:", connID, err)

		conn.Close()
		delete(p.connections, connID)
	}
}
