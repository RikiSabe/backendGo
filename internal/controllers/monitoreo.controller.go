package controllers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type monitoreo struct {
}

type Localizacion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	Monitoreo            monitoreo
	channelUbicaciones   = make(chan Localizacion)
	managerUbicacionesWS = NewWebSocketManager()
	mu                   sync.Mutex
)

func (monitoreo) ObtenerUbicacionLecturadorWS(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	mu.Lock()
	managerUbicacionesWS.AddConn(ws)
	mu.Unlock()
	for {
		select {
		case ubicacion := <-channelUbicaciones: // Recibir ubicación del canal
			managerUbicacionesWS.Broadcast(ubicacion) // Enviar la ubicación a través de WebSocket
		}
	}

}
