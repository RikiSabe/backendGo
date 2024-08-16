package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"net/http"
)

func GetPersonasHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var personas []*models.Persona
	database.DB.Find(&personas)
	err := json.NewEncoder(w).Encode(&personas)
	if err != nil {
		return
	}
}

func PostPersonasHandle(w http.ResponseWriter, r *http.Request) {
	var personaDAO models.Lecturacion // dao.Personas
	json.NewDecoder(r.Body).Decode(&personaDAO)
	var persona models.Persona

	err := utils.JsonDaoToModel(personaDAO, &persona)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resul := database.DB.Create(&persona)
	if resul.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&persona)
	if err != nil {
		return
	}
}
