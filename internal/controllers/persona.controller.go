package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"net/http"
)

func ObtenerPersonas(w http.ResponseWriter, r *http.Request) {
	var personas []models.Persona
	db.GDB.Find(&personas)

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&personas)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func SubirPersonas(w http.ResponseWriter, r *http.Request) {
	var persona models.Persona // dao.Personas
	err := json.NewDecoder(r.Body).Decode(&persona)
	if err != nil {
		http.Error(w, "Error en el parseo de formulario", http.StatusBadRequest)
		return
	}
	tx := db.GDB.Begin()

	if err := db.GDB.Create(&persona).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ha ocurrido un error en el registro", http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&persona)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
