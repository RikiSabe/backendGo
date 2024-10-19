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

func ObtenerGrupos(w http.ResponseWriter, r *http.Request) {
	var grupos []models.Grupo
	if err := services.Grupo.GetAll(&grupos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&grupos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SubirGrupo(w http.ResponseWriter, r *http.Request) {
	var grupo models.Grupo
	if err := json.NewDecoder(r.Body).Decode(&grupo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := services.Grupo.Save(&grupo); err != nil {
		http.Error(w, "Ha ocurrido un error al guardar en la BD", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&grupo); err != nil {
		http.Error(w, "Ha ocurrido un error al parsear a JSON", http.StatusInternalServerError)
		return
	}
}

func ModificarGrupo(w http.ResponseWriter, r *http.Request) {
	var grupoActualizado models.Grupo
	cod := mux.Vars(r)["cod"]

	// Buscar el grupo existente por su código
	var grupoExistente models.Grupo
	if err := services.Grupo.GetById(&grupoExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&grupoActualizado); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	grupoExistente.CodUsuario = grupoActualizado.CodUsuario
	grupoExistente.CodRuta = grupoActualizado.CodRuta
	// Guardar los cambios en el grupo existente
	if err := db.GDB.Save(&grupoExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&grupoExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}
