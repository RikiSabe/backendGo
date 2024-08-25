package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ObtenerRutas(w http.ResponseWriter, r *http.Request) {
	var rutas []models.Ruta
	if err := services.Ruta.GetAll(&rutas); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&rutas); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func ObtenerRuta(w http.ResponseWriter, r *http.Request) {
	var ruta models.Ruta
	cod := mux.Vars(r)["cod"]
	if err := services.Ruta.GetById(&ruta, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&ruta); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func SubirRuta(w http.ResponseWriter, r *http.Request) {
	var ruta models.Ruta
	if err := json.NewDecoder(r.Body).Decode(&ruta); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := services.Ruta.Save(&ruta); err != nil {
		http.Error(w, "Ha ocurrido un error al guardar en la BD", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&ruta); err != nil {
		http.Error(w, "Ha ocurrido un error al parsear a JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ModificarRuta(w http.ResponseWriter, r *http.Request) {
	var rutaActualizada models.Ruta
	cod := mux.Vars(r)["cod"]

	// Buscar la lecturación existente por su código
	var rutaExistente models.Ruta
	if err := services.Ruta.GetById(&rutaExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&rutaActualizada); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rutaExistente.Nombre = rutaActualizada.Nombre
	rutaExistente.Zona = rutaActualizada.Zona
	rutaExistente.Estado = rutaActualizada.Estado
	// Guardar los cambios en la lecturación existente
	if err := db.GDB.Save(&rutaExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&rutaExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}
