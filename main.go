package main

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	port := "5000"
	r := mux.NewRouter()
	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}
	// Conectar a la base de datos
	err = db.Connection()
	if err != nil {
		log.Printf("Error al cargar la base de datos: %v", err)
		return
	}
	// Migrar los modelos
	if err := db.GDB.AutoMigrate(
		&models.Lecturacion{},
		&models.Usuario{},
		&models.Ruta{},
		&models.Medidor{},
		&models.Direccion{},
		&models.Grupo{},
		&models.Critica{},
		&models.Persona{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}
	// Cargar endPoints
	routers.InitEndPoints(r)
	// Iniciar el servidor
	fmt.Printf("Servidor corriendo en puerto: %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
