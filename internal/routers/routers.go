package routers

import (
	c "backend/internal/controllers"
	"backend/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	r.Use(middleware.LoggingHandler)
	api := r.PathPrefix("/api").Subrouter()
	ws := r.PathPrefix("/ws").Subrouter()
	endPointsAPI(api)
	endPointsWS(ws)

}
func endPointsAPI(api *mux.Router) {
	//Subrutas
	v1 := api.PathPrefix("/v1").Subrouter()
	v1Medidores := v1.PathPrefix("/medidores").Subrouter()
	v1Personas := v1.PathPrefix("/personas").Subrouter()
	v1Lecturaciones := v1.PathPrefix("/lecturaciones").Subrouter()

	//v1 Personas
	v1Personas.HandleFunc("", c.ObtenerPersonas).Methods(http.MethodGet)
	v1Personas.HandleFunc("", c.SubirPersonas).Methods(http.MethodPost)

	//v1 Medidores
	v1Medidores.HandleFunc("", c.ObtenerMedidores).Methods(http.MethodGet)
	v1Medidores.HandleFunc("/{cod}", c.ObtenerMedidor).Methods(http.MethodGet)
	v1Medidores.HandleFunc("", c.PostMedidor).Methods(http.MethodPost)
	v1Medidores.HandleFunc("/{cod}", c.EliminarMedidor).Methods(http.MethodDelete)

	//v1 Lecturaciones
	v1Lecturaciones.HandleFunc("", c.ObtenerLecturaciones).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("/{cod}", c.ObtenerLecturacion).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("", c.SubirLecturacion).Methods(http.MethodPost)
	v1Lecturaciones.HandleFunc("/{cod}", c.EliminarLecturacion).Methods(http.MethodDelete)

}
func endPointsWS(ws *mux.Router) {
	//Subrutas
	v1 := ws.PathPrefix("/v1").Subrouter()
	v1Medidores := v1.PathPrefix("/medidores").Subrouter()

	//v1 Medidores
	v1Medidores.HandleFunc("/medidor", c.ObtenerMedidoresWS)
}
