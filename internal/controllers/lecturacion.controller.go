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

func ObtenerHistorialLecturaciones(w http.ResponseWriter, r *http.Request) {
	var lecturaciones []struct {
		NombreMedidor string `json:"nombre_medidor"`
		Hora          string `json:"hora"`
		Fecha         string `json:"fecha"`
	}

	codRuta := mux.Vars(r)["cod_ruta"]

	query := `SELECT m.nombre AS nombre_medidor, l.hora AS hora, l.fecha as fecha
				FROM lecturacion AS l
				LEFT JOIN medidor AS m ON l.cod_medidor = m.cod
				WHERE m.cod_ruta = ?`

	tx := db.GDB.Begin()
	if err := tx.Raw(query, codRuta).Scan(&lecturaciones).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturaciones); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	var lecturacionExistente models.Lecturacion
	if err := services.Lecturacion.GetById(&lecturacionExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&lecturacionActualizada); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lecturacionExistente.Fecha = lecturacionActualizada.Fecha
	lecturacionExistente.NroRegistro = lecturacionActualizada.NroRegistro

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
	var datosLecturacion struct {
		NombreMedidor string `json:"nombreMedidor"`
		Usuario       string `json:"usuario"`
		Medicion      uint   `json:"medicion"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datosLecturacion); err != nil {
		http.Error(w, "Error al procesar los datos", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	nuevaLecturacion := models.Lecturacion{
		CodRuta:       *medidor.CodRuta,
		CodLecturador: usuario.COD,
		CodMedidor:    medidor.COD,
		Medicion:      &datosLecturacion.Medicion,
		Hora:          datatypes.Time(datatypes.NewTime(time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0)),
		Fecha:         datatypes.Date(time.Now()),
	}

	if err := tx.Create(&nuevaLecturacion).Error; err != nil {
		http.Error(w, "Error al guardar la lecturación", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
}
