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

func ObtenerCriticas(w http.ResponseWriter, r *http.Request) {
	var criticas []models.Critica
	if err := services.Critica.GetAll(&criticas); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&criticas); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func ObtenerCritica(w http.ResponseWriter, r *http.Request) {
	var critica models.Critica
	cod := mux.Vars(r)["cod"]
	if err := services.Critica.GetById(&critica, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&critica); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SubirCritica(w http.ResponseWriter, r *http.Request) {
	var critica models.Critica
	if err := json.NewDecoder(r.Body).Decode(&critica); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := services.Critica.Save(&critica); err != nil {
		http.Error(w, "Ha ocurrido un error al guardar en la BD", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&critica); err != nil {
		http.Error(w, "Ha ocurrido un error al parsear a JSON", http.StatusInternalServerError)
		return
	}

}

func ModificarCritica(w http.ResponseWriter, r *http.Request) {
	var criticaActualizada models.Critica
	cod := mux.Vars(r)["cod"]

	// Buscar la lecturación existente por su código
	var criticaExistente models.Critica
	if err := services.Critica.GetById(&criticaExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&criticaActualizada); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Actualizar los campos de la lecturación existente con los valores de la lecturación actualizada
	// lecturacionExistente.Fecha = lecturacionActualizada.Fecha
	// lecturacionExistente.NroRegistro = lecturacionActualizada.NroRegistro
	// (Actualizar otros campos según sea necesario)
	criticaExistente.Descripcion = criticaActualizada.Descripcion
	criticaExistente.Tipo = criticaActualizada.Tipo
	criticaExistente.Estado = criticaActualizada.Estado
	// Guardar los cambios en la lecturación existente
	if err := db.GDB.Save(&criticaExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&criticaExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}
