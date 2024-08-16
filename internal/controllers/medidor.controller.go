package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	mu                 sync.Mutex
	medidoresChannel   = make(chan []*models.Medidor)
	wsManagerMedidores = NewWebSocketManager()
)

func ObtenerMedidores(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var medidores []models.Medidor
	tx := db.GDB.Begin()
	if err := tx.Find(&medidores).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			tx.Rollback()
			http.Error(w, "Ha ocurrido un error al obtener la lista", http.StatusInternalServerError)
			return
		}
	}

	err := json.NewEncoder(w).Encode(&medidores)
	if err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
		return
	}
}

func ObtenerMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor
	codigoMedidor := mux.Vars(r)["codmedidor"]
	tx := db.GDB.Begin()
	if err := tx.Where("codigo = ?", codigoMedidor).First(&medidor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Medidor no encontrado", http.StatusNotFound)
			return
		}
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&medidor)
	if err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusInternalServerError)
		return
	}
}

func PostMedidor(w http.ResponseWriter, r *http.Request) {
	var medidor models.Medidor
	err := json.NewDecoder(r.Body).Decode(&medidor)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	tx := db.GDB.Begin()
	if err := tx.Create(&medidor).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ha ocurrido un error en el regitro del medidor", http.StatusInternalServerError)
		return
	}
	var medidores []*models.Medidor
	if err := tx.Find(&medidores).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ha ocurrido un error al obtener la lista del medidores", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	// Enviar la nueva lista al canal
	mu.Lock()
	medidoresChannel <- medidores
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&medidor)
	if err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}

func ObtenerMedidoresWS(w http.ResponseWriter, r *http.Request) {
	upgrader := NewUpgrader()
	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	// get socket
	var medidores []models.Medidor
	wsManagerMedidores.AddConn(ws)
	tx := db.GDB.Begin()
	if err := tx.Find(&medidores).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

		} else {
			tx.Rollback()
			return
		}
	} else {
		err := ws.WriteJSON(&medidores)
		if err != nil {
			return
		}
	}

	tx.Commit()
	for {
		select {
		case medidoresUpdated := <-medidoresChannel:
			wsManagerMedidores.Broadcast(&medidoresUpdated)
		}
	}
}
