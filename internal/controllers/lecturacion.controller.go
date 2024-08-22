package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
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
