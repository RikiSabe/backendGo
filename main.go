package main

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	port := "5000"

	// Conectar a la base de datos
	err = db.Connection()
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return
	}

	// Migrar los modelos
	if err := db.GDB.AutoMigrate(
		&models.Direccion{},
		&models.Ruta{},
		&models.Medidor{},
		&models.Critica{},
		&models.Lecturacion{},
		&models.Usuario{},
		&models.Grupo{},
		&models.Persona{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}

	// Configurar el router
	r := mux.NewRouter()
	routers.InitEndPoints(r)

	// Configuración de CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Iniciar el servidor
	fmt.Printf("Servidor corriendo en puerto: %s\n", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
