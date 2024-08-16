package controllers

import (
	"backend/database"
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var medidoresLIst = make(chan []*models.Medidor)

func GetMedidores(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var medidores []models.Medidor
	database.DB.Find(&medidores)
	err := json.NewEncoder(w).Encode(&medidores)
	if err != nil {
		return
	}
}

func GetMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor
	codigoMedidor := mux.Vars(r)["codmedidor"]
	database.DB.Where("codigo = ?", codigoMedidor).First(&medidor)
	if medidor.Codigo == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(&medidor)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor
	json.NewDecoder(r.Body).Decode(&medidor)
	resul := database.DB.Create(&medidor)
	if resul.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(w).Encode(&medidor)
	if err != nil {
		return
	}
	// function socket
	var medidores []*models.Medidor
	database.DB.Find(&medidores)
	// send new list to old w/ socket
	medidoresLIst <- medidores
	w.WriteHeader(http.StatusOK)
}

var upgrader = websocket.Upgrader{}

func ObtenerMedidoresWS(w http.ResponseWriter, r *http.Request) {
	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	// get socket
	var medidores []models.Medidor
	database.DB.Find(&medidores)
	err := ws.WriteJSON(&medidores)
	if err != nil {
		return
	}
	for {
		select {
		case medidoresUpdated := <-medidoresLIst:
			err := ws.WriteJSON(&medidoresUpdated)
			if err != nil {
				return
			}
		}
	}
}
