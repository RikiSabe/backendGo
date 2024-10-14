package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
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

func (monitoreo) ObtenerUbicacionLecturadoresWS(w http.ResponseWriter, r *http.Request) {

}

// ObtenerUbicacionLecturadorWS permite al usuario conectarse para enviar ubicación
func (monitoreo) ObtenerUbicacionLecturadorWS(w http.ResponseWriter, r *http.Request) {
	type Ubicacion struct {
		Longitud float64 `json:"longitud"`
		Latitud  float64 `json:"latitud"`
	}

	var location Ubicacion
	var upgrader = websocket.Upgrader{}

	// Establecer la conexión WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al hacer upgrade a WebSocket:", err)
		return
	}
	defer ws.Close()

	// Verificar el header de autorización
	authorizationHeader := r.Header.Get("Authorization")
	if err := verificarBearerHeader(authorizationHeader); err != nil {
		log.Println("Encabezado de autorización inválido:", err)
		return
	}

	// Agregar conexión al manager
	mu.Lock()
	managerUbicacionesWS.AddConn(ws)
	mu.Unlock()
	log.Printf("Cantidad de conexiones: %d", len(managerUbicacionesWS.conns))
	// Loop para leer ubicaciones enviadas por el cliente
	for {
		// Leer ubicación en formato JSON
		err := ws.ReadJSON(&location)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Conexión cerrada inesperadamente:", err)
				managerUbicacionesWS.RemoveConn(ws)
			} else {
				log.Println("Error al leer JSON:", err)
			}
			break // Salir del loop en caso de error
		}

		// Loguear la ubicación recibida
		log.Printf("Ubicación recibida: Longitud: %f, Latitud: %f\n", location.Longitud, location.Latitud)
	}

	// Lógica para manejar el canal (comentado)
	// for {
	// 	select {
	// 	case ubicacion := <-channelUbicaciones: // Recibir ubicación del canal
	// 		managerUbicacionesWS.Broadcast(ubicacion) // Enviar la ubicación a través de WebSocket
	// 	}
	// }
}

func verificarBearerHeader(authHeader string) error {
	// Verificar el encabezado de autorización
	if authHeader == "" {
		return fmt.Errorf("No se proporcionó el encabezado de autorización")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return fmt.Errorf("El encabezado de autorización debe comenzar con Bearer ")
	}
	return nil
}

func verificarToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
