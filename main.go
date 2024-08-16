package main

import (
	"backend/controllers"
	"backend/database"

	// "backend/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{}

func main() {
	database.Conection()
	/*err1 := database.DB.AutoMigrate(
		&models.ReporteError{},
		&models.Ruta{},
		&models.Medidor{},
		&models.Lecturador{},
		&models.Administrador{},
		&models.Lecturacion{},
	)
	if err1 != nil {
		log.Println("Error al migrar los modelos de la db")
		return
	}*/
	port := "5000"
	fmt.Println("Corriendo")
	r := mux.NewRouter()
	r.Use(loggingHandler)
	r.HandleFunc("/", Home)
	r.HandleFunc("/api/persona", controllers.GetPersonasHandle).Methods(http.MethodGet)
	r.HandleFunc("/api/persona", controllers.PostPersonasHandle).Methods(http.MethodPost)
	r.HandleFunc("/api/medidor", controllers.GetMedidores).Methods(http.MethodGet)
	r.HandleFunc("/api/medidor/{codmedidor}", controllers.GetMedidor).Methods(http.MethodGet)
	r.HandleFunc("/api/medidor", controllers.PostMedidor).Methods(http.MethodPost)
	r.HandleFunc("/ws/medidor", controllers.ObtenerMedidoresWS)
	http.Handle("/", r)
	// ":" = localhost
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Principal")
}
func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
