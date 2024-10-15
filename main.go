package main

import (
	"backend/internal/db"
	"backend/internal/routers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
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
	// &models.Direccion{},
	// &models.Medidor{},
	// &models.Ruta{},
	// &models.Critica{},
	// &models.Grupo{},
	// &models.Usuario{},
	// &models.Lecturacion{},
	// &models.Persona{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}
	// go RutinaDiaria(12, 00)

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

func RutinaDiaria(hora, minuto int) {
	for {
		horaActual := time.Now()

		// Calcula la próxima ejecución según la hora y minuto pasados por parámetro
		proximaEjecucion := time.Date(horaActual.Year(), horaActual.Month(), horaActual.Day(), hora, minuto, 0, 0, horaActual.Location())

		// Si la hora ya ha pasado hoy, programa para el día siguiente
		if horaActual.After(proximaEjecucion) {
			proximaEjecucion = proximaEjecucion.Add(24 * time.Hour)
		}

		// Calcula cuánto falta para la primera ejecución
		duracionHastaEjecucion := time.Until(proximaEjecucion)
		fmt.Printf("Esperando hasta las %s para iniciar la rutina...\n", proximaEjecucion.Format("15:04"))

		// Espera hasta la hora especificada
		time.Sleep(duracionHastaEjecucion)

		// Empieza a imprimir "Hola Mundo" cada 5 minutos
		for {
			horaActual = time.Now()
			fmt.Println("Hola Mundo - ", horaActual.Format("2006-01-02 15:04:05"))

			// Espera 5 minutos antes de volver a imprimir
			time.Sleep(1 * time.Minute)
		}
	}
}
