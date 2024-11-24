package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func ObtenerLecturaciones(w http.ResponseWriter, r *http.Request) {
	var lecturaciones []models.Lecturacion
	if err := services.Lecturacion.GetAll(&lecturaciones); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturaciones); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ObtenerLecturacion(w http.ResponseWriter, r *http.Request) {
	var lecturacion models.Lecturacion
	cod := mux.Vars(r)["cod"]
	if err := services.Lecturacion.GetById(&lecturacion, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturacion); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func SubirLecturacion(w http.ResponseWriter, r *http.Request) {
	var lecturacion models.Lecturacion
	if err := json.NewDecoder(r.Body).Decode(&lecturacion); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := services.Lecturacion.Save(&lecturacion); err != nil {
		http.Error(w, "Ha ocurrido un error al guardar en la BD", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&lecturacion); err != nil {
		http.Error(w, "Ha ocurrido un error al parsear a JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ModificarLecturacion(w http.ResponseWriter, r *http.Request) {
	var lecturacionActualizada models.Lecturacion
	cod := mux.Vars(r)["cod"]

	// Buscar la lecturación existente por su código
	var lecturacionExistente models.Lecturacion
	if err := services.Lecturacion.GetById(&lecturacionExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&lecturacionActualizada); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Actualizar los campos de la lecturación existente con los valores de la lecturación actualizada
	lecturacionExistente.Fecha = lecturacionActualizada.Fecha
	lecturacionExistente.NroRegistro = lecturacionActualizada.NroRegistro

	// Guardar los cambios en la lecturación existente
	if err := db.GDB.Save(&lecturacionExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturacionExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}

func CrearLecturacion(w http.ResponseWriter, r *http.Request) {
	// Estructura para recibir los datos del cuerpo de la solicitud
	var datosLecturacion struct {
		NombreMedidor string `json:"nombreMedidor"`
		Usuario       string `json:"usuario"`
		Medicion      uint   `json:"medicion"`
	}

	// Decodificar el JSON recibido
	if err := json.NewDecoder(r.Body).Decode(&datosLecturacion); err != nil {
		// log.Println("Error al decodificar el JSON:", err.Error())
		http.Error(w, "Error al procesar los datos", http.StatusBadRequest)
		return
	}

	// Iniciar transacción
	tx := db.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Buscar el medidor por nombre
	var medidor models.Medidor
	if err := tx.Where("nombre = ?", datosLecturacion.NombreMedidor).First(&medidor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Medidor no encontrado", http.StatusNotFound)
			tx.Rollback()
			return
		}
		http.Error(w, "Error al buscar el medidor", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	// Buscar el usuario por nombre
	var usuario models.Usuario
	if err := tx.Where("usuario = ?", datosLecturacion.Usuario).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			tx.Rollback()
			return
		}
		http.Error(w, "Error al buscar el usuario", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	// Crear una nueva lecturación
	nuevaLecturacion := models.Lecturacion{
		CodRuta:       *medidor.CodRuta,
		CodLecturador: usuario.COD,
		CodMedidor:    medidor.COD,
		Medicion:      &datosLecturacion.Medicion,
		Hora:          datatypes.Time(datatypes.NewTime(time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0)),
		Fecha:         datatypes.Date(time.Now()),
	}

	// Guardar la lecturación en la base de datos
	if err := tx.Create(&nuevaLecturacion).Error; err != nil {
		http.Error(w, "Error al guardar la lecturación", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	// Confirmar la transacción
	tx.Commit()

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Lecturación creada exitosamente"))
}
