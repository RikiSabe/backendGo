package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ObtenerLecturaciones(w http.ResponseWriter, r *http.Request) {
	var lecturaciones []models.Lecturacion
	if err := services.Lecturacion.GetAll(&lecturaciones); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&lecturaciones); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
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
func EliminarLecturacion(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if err := services.Lecturacion.Delete(cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	// (Actualizar otros campos según sea necesario)

	// Guardar los cambios en la lecturación existente
	if err := services.Lecturacion.Save(&lecturacionExistente); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturacionExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}
