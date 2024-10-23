package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

/*
var (

	mu                 sync.Mutex
	medidoresChannel   = make(chan []models.Medidor)
	wsManagerMedidores = NewWebSocketManager()

)
*/
func ObtenerMedidores(w http.ResponseWriter, r *http.Request) {
	var medidores []models.Medidor

	// Llamar al servicio para obtener todos los medidores activos
	if err := services.Medidor.GetAll(&medidores); err != nil {
		http.Error(w, "Ha ocurrido un error al obtener la lista de medidores", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&medidores); err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
		return
	}
}

// func ObtenerMedidoresByRuta(w http.ResponseWriter, r *http.Request) {
// 	var medidores []models.Medidor
// 	codigoRuta := mux.Vars(r)["cod_ruta"]

// 	if err := services.Medidor.GetByRuta(&medidores, codigoRuta); err != nil {
// 		http.Error(w, "Ha ocurrido un error al obtener la lista de medidores por ruta", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(&medidores); err != nil {
// 		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
// 		return
// 	}
// }

func ObtenerMedidoresByRuta(w http.ResponseWriter, r *http.Request) {
	// Estructura que combina los datos de Medidor y las coordenadas de Direccion
	var medidores []struct {
		CodMedidor  uint    `json:"codMedidor"`
		Estado      string  `json:"estado"`
		Medicion    int     `json:"medicion"`
		Nombre      string  `json:"nombre"`
		Propietario string  `json:"propietario"`
		CodRuta     uint    `json:"codRuta"`
		Latitud     float32 `json:"latitud,omitempty"`
		Longitud    float32 `json:"longitud,omitempty"`
	}

	// Obtén el código de la ruta desde la URL
	codigoRuta := mux.Vars(r)["cod_ruta"]

	// Consulta SQL personalizada para obtener los medidores con sus coordenadas
	query := `SELECT m.cod as cod_medidor, m.estado, m.medicion, m.nombre, m.propietario, m.cod_ruta, 
					 d.longitud, d.latitud
			  FROM medidor m
			  LEFT JOIN direccion d ON m.cod_direccion = d.cod
			  WHERE m.cod_ruta = ? AND m.estado = 'activo';`

	// Inicia la transacción
	tx := db.GDB.Begin()
	if err := tx.Raw(query, codigoRuta).Scan(&medidores).Error; err != nil {
		// Manejo de errores
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	// Respuesta con los datos de medidores
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&medidores); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ObtenerDireccion(w http.ResponseWriter, r *http.Request) {

	var direccion struct {
		Longitud float32 `json:"longitud"`
		Latitud  float32 `json:"latitud"`
	}

	codDireccion := mux.Vars(r)["cod_direccion"]
	query := `
		select d.longitud, d.latitud 
		FROM direccion as d
		left join medidor m ON m.cod_direccion = d.cod 
		where d.cod = ?
`

	tx := db.GDB.Begin()
	if err := tx.Raw(query, codDireccion).Scan(&direccion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&direccion); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ObtenerDireccionMedidores(w http.ResponseWriter, r *http.Request) {

	var direccion []struct {
		Longitud float64 `json:"longitud"`
		Latitud  float64 `json:"latitud"`
	}

	codRuta := mux.Vars(r)["cod_ruta"]
	query := `SELECT d.latitud, d.longitud
				FROM medidor as m
				left JOIN direccion as d ON d.cod = m.cod_direccion
				WHERE m.cod_ruta = ?`

	tx := db.GDB.Begin()
	if err := tx.Raw(query, codRuta).Scan(&direccion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&direccion); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ObtenerMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor
	codigoMedidor := mux.Vars(r)["cod"]

	// Llamar al servicio para obtener un medidor por su código
	if err := services.Medidor.GetByCod(&medidor, codigoMedidor); err != nil {
		http.Error(w, "Medidor no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&medidor); err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
		return
	}
}

func AgregarMedidor(w http.ResponseWriter, r *http.Request) {
	var datoMedidor struct {
		Nombre      string  `json:"nombre"`
		Propietario string  `json:"propietario"`
		Tipo        string  `json:"tipo"`
		CodRuta     uint    `json:"codRuta"`
		Latitud     float64 `json:"latitud"`
		Longitud    float64 `json:"longitud"`
	}

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&datoMedidor); err != nil {
		log.Println(err.Error())
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Iniciar la transacción
	tx := db.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verificar la existencia de la ruta
	if err := tx.Model(models.Ruta{}).Where("cod = ?", datoMedidor.CodRuta).First(&models.Ruta{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Ruta no encontrada", http.StatusNotFound)
			tx.Rollback()
			return
		}
		tx.Rollback()
		http.Error(w, "Error al verificar la ruta", http.StatusInternalServerError)
		return
	}

	// Crear el registro de la dirección
	direccion := models.Direccion{
		Latitud:  datoMedidor.Latitud,
		Longitud: datoMedidor.Longitud,
	}

	// Guardar la dirección en la base de datos
	if err := tx.Create(&direccion).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar la dirección", http.StatusInternalServerError)
		return
	}

	// Crear el objeto Medidor
	medidor := models.Medidor{
		Nombre:       datoMedidor.Nombre,
		Propietario:  datoMedidor.Propietario,
		Tipo:         datoMedidor.Tipo,
		Estado:       "activo", // Estado inicial del medidor
		CodRuta:      &datoMedidor.CodRuta,
		CodDireccion: &direccion.COD,
	}

	// Guardar el medidor en la base de datos
	if err := tx.Create(&medidor).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar el medidor", http.StatusInternalServerError)
		return
	}

	// Confirmar la transacción
	tx.Commit()

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Medidor agregado correctamente"))
}

func PostMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor

	// Decodificar el JSON recibido en el request
	log.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&medidor); err != nil {
		log.Println(err)
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	medidor.Estado = "activo"
	// Llamar al servicio para guardar el nuevo medidor
	if err := services.Medidor.Save(&medidor); err != nil {
		http.Error(w, "Ha ocurrido un error al registrar el medidor", http.StatusInternalServerError)
		return
	}

	// Obtener la lista actualizada de medidores
	var medidores []models.Medidor
	if err := services.Medidor.GetAll(&medidores); err != nil {
		http.Error(w, "Ha ocurrido un error al obtener la lista de medidores", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&medidor); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}

func EliminarMedidor(w http.ResponseWriter, r *http.Request) {
	codigoMedidor := mux.Vars(r)["cod"]

	// Llamar al servicio para marcar el medidor como inactivo
	if err := services.Medidor.Delete(codigoMedidor); err != nil {
		http.Error(w, "Ha ocurrido un error al eliminar el medidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ModificarMedidor(w http.ResponseWriter, r *http.Request) {
	var medidorActualizado models.Medidor
	codigoMedidor := mux.Vars(r)["cod"]

	// Buscar el medidor existente por su código
	var medidorExistente models.Medidor
	if err := services.Medidor.GetByCod(&medidorExistente, codigoMedidor); err != nil {
		http.Error(w, "Medidor no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&medidorActualizado); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Actualizar los campos del medidor existente con los valores del medidor actualizado
	medidorExistente.Nombre = medidorActualizado.Nombre
	medidorExistente.Propietario = medidorActualizado.Propietario
	// Guardar los cambios en el medidor existente
	if err := db.GDB.Save(&medidorExistente).Error; err != nil {
		http.Error(w, "Ha ocurrido un error al actualizar el medidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&medidorExistente); err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
		return
	}
}

/*
func ObtenerMedidoresWS(w http.ResponseWriter, r *http.Request) {
	upgrader := NewUpgrader()
	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	// get socket
	var medidores []models.Medidor
	wsManagerMedidores.AddConn(ws)
	tx := db.GDB.Begin()
	if err := tx.Find(&medidores).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return
		} else {
			tx.Rollback()
			return
		}
	} else {
		err := ws.WriteJSON(&medidores)
		if err != nil {
			return
		}
	}

	tx.Commit()
	for {
		select {
		case medidoresUpdated := <-medidoresChannel:
			wsManagerMedidores.Broadcast(&medidoresUpdated)
		}
	}
}

func ActualizarMedidor(w http.ResponseWriter, r *http.Request) {
	//cod := mux.Vars(r)["cod"]

}
*/
