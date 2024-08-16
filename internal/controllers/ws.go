package controllers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WebSocketManager struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		conns: make(map[*websocket.Conn]bool),
	}
}
func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
func (manager *WebSocketManager) AddConn(ws *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.conns[ws] = true
}

func (manager *WebSocketManager) RemoveConn(ws *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.conns, ws)
}

func (manager *WebSocketManager) Broadcast(message interface{}) {
	manager.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(manager.conns))
	for ws := range manager.conns {
		conns = append(conns, ws)
	}
	manager.mu.Unlock()

	var wg sync.WaitGroup
	for _, ws := range conns {
		wg.Add(1)
		go func(conn *websocket.Conn) {
			defer wg.Done()
			err := conn.WriteJSON(message)
			if err != nil {
				manager.mu.Lock()
				conn.Close()
				delete(manager.conns, conn)
				manager.mu.Unlock()
			}
		}(ws)
	}
	wg.Wait()
}

// Handler para manejar las conexiones WebSocket
func WebSocketHandler(manager *WebSocketManager, upgrader *websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}
		defer ws.Close()

		manager.AddConn(ws)
		defer manager.RemoveConn(ws)

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				break
			}
			manager.Broadcast(message)
		}
	}
}
