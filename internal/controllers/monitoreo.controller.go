package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type monitoreo struct{}

type Localizacion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Ubicacion struct {
	Longitud float64 `json:"longitud"`
	Latitud  float64 `json:"latitud"`
}

var (
	Monitoreo monitoreo
	//Mobile
	managerUbicacionesWS = NewWebSocketManager()
	//Web
	managerAdminWS          = NewWebSocketManager()
	mu                      sync.Mutex
	channelUbicacionesUsers = make(chan map[string]Ubicacion, 10)
	ubicacionesUsers        = make(map[string]Ubicacion) // Cambiado a un mapa
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permitir todas las solicitudes, aunque deberías personalizarlo según tus necesidades.
		return true
	},
}

func (monitoreo) ObtenerUbicacionesLecturadorWS(w http.ResponseWriter, r *http.Request) {
	// Establecer la conexión WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al hacer upgrade a WebSocket:", err)
		return
	}
	defer func() {
		log.Println("Cerrando conexión WebSocket")
		ws.Close() // Asegura que la conexión se cierra
	}()

	// Añadir conexión al manager
	managerAdminWS.AddConn(ws)

	// Configurar un manejador de cierre
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("La conexión WebSocket se ha cerrado:", text)
		managerAdminWS.RemoveConn(ws) // Remover la conexión del manager al cerrar
		return nil
	})

	// Configurar ticker para enviar pings periódicos
	pingTicker := time.NewTicker(30 * time.Second) // Ajusta el intervalo de pings según sea necesario
	defer pingTicker.Stop()

	for {
		select {
		case ch := <-channelUbicacionesUsers:
			// Transmitir la ubicación al cliente
			log.Println(ch)
			managerAdminWS.Broadcast(ch)

		case <-pingTicker.C:
			// Enviar un ping para mantener la conexión viva
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error al enviar mensaje de ping:", err)
				return // Cierra la conexión en caso de error
			}

		default:
			// Introduce una pequeña pausa para evitar un bucle ocupado (busy loop)
			time.Sleep(100 * time.Millisecond) // Ajusta este valor según las necesidades de rendimiento
		}
	}
}

// Para Web
// func (monitoreo) ObtenerUbicacionesLecturadorWS(w http.ResponseWriter, r *http.Request) {
// 	var upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			// Aquí puedes permitir todas las solicitudes, aunque deberías personalizarlo según tus necesidades.
// 			return true
// 		},
// 	}
// 	// Establecer la conexión WebSocket
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Error al hacer upgrade a WebSocket:", err)
// 		return
// 	}
// 	defer func() {
// 		log.Println("Cerrando conexión WebSocket")
// 		ws.Close() // Asegura que la conexión se cierra
// 	}()

// 	managerAdminWS.AddConn(ws)

// 	// Configurar un manejador de cierre
// 	ws.SetCloseHandler(func(code int, text string) error {
// 		log.Println("La conexión WebSocket se ha cerrado:", text)
// 		managerAdminWS.RemoveConn(ws) // Opcional: remueve la conexión del manager
// 		return nil
// 	})

// 	for {
// 		select {
// 		case ch := <-channelUbicacionesUsers:
// 			// Transmitir la ubicación al cliente
// 			log.Println(ch)
// 			managerAdminWS.Broadcast(ch)

// 		default:
// 			// Verifica si hay errores en la conexión
// 			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				log.Println("Error al enviar mensaje de ping:", err)
// 				return // Cierra la conexión en caso de error
// 			}
// 		}
// 	}
// }

// ObtenerUbicacionLecturadorWS permite al usuario conectarse para enviar ubicación
func (monitoreo) ObtenerUbicacionLecturadorWS(w http.ResponseWriter, r *http.Request) {
	tokenAuth := r.Header.Get("Authorization")
	log.Println("token:", tokenAuth)

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
	token, err := verificarBearerHeader(tokenAuth)
	if err != nil {
		log.Println("Encabezado de autorización inválido:", err)
		return
	}

	jwtToken, err := verifyToken(token)
	if err != nil || !jwtToken.Valid {
		log.Println("Token inválido:", err)
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error al obtener claims del token")
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		log.Println("No se pudo obtener el nombre de usuario del token")
		return
	}
	log.Println(username)

	// Agregar conexión al manager
	mu.Lock()
	managerUbicacionesWS.AddConn(ws)
	mu.Unlock()

	log.Printf("Cantidad de conexiones: %d", len(managerUbicacionesWS.conns))

	for {
		err := ws.ReadJSON(&location)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Conexión cerrada inesperadamente:", err)
				mu.Lock()
				managerUbicacionesWS.RemoveConn(ws)
				mu.Unlock()
			} else {
				log.Println("Error al leer JSON:", err)
			}
			mu.Lock()
			delete(ubicacionesUsers, username)
			channelUbicacionesUsers <- ubicacionesUsers
			mu.Unlock()
			break
		}

		// Validar la ubicación
		if location.Longitud < -180 || location.Longitud > 180 || location.Latitud < -90 || location.Latitud > 90 {
			log.Println("Ubicación inválida recibida:", location)
			continue
		}

		// Agregar la ubicación del usuario al mapa
		mu.Lock()
		ubicacionesUsers[username] = location
		channelUbicacionesUsers <- ubicacionesUsers
		mu.Unlock()

		log.Printf("Ubicación recibida: Longitud: %f, Latitud: %f\n", location.Longitud, location.Latitud)
	}
}

func verificarBearerHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("No se proporcionó el encabezado de autorización")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("El encabezado de autorización debe comenzar con Bearer ")
	}
	return authHeader[len(bearerPrefix):], nil
}
