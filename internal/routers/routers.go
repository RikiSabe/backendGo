package routers

import (
	c "backend/internal/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	// r.Use(middleware.LoggingHandler)

	api := r.PathPrefix("/api").Subrouter()
	ws := r.PathPrefix("/ws").Subrouter()

	endPointsAPI(api)
	endPointsWS(ws)
}

func endPointsAPI(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	v1Medidores := v1.PathPrefix("/medidores").Subrouter()
	v1Personas := v1.PathPrefix("/personas").Subrouter()
	v1Lecturaciones := v1.PathPrefix("/lecturaciones").Subrouter()
	v1Criticas := v1.PathPrefix("/criticas").Subrouter()
	v1Rutas := v1.PathPrefix("/rutas").Subrouter()
	v1Usuarios := v1.PathPrefix("/usuarios").Subrouter()
	// v1 Personas
	v1Personas.HandleFunc("", c.ObtenerPersonas).Methods(http.MethodGet)
	v1Personas.HandleFunc("", c.SubirPersonas).Methods(http.MethodPost)

	// v1 Medidores
	v1Medidores.HandleFunc("", c.ObtenerMedidores).Methods(http.MethodGet)
	v1Medidores.HandleFunc("/{cod}", c.ObtenerMedidor).Methods(http.MethodGet)
	v1Medidores.HandleFunc("", c.PostMedidor).Methods(http.MethodPost)
	v1Medidores.HandleFunc("/{cod}", c.ModificarMedidor).Methods(http.MethodPut)
	v1Medidores.HandleFunc("/{cod}", c.EliminarMedidor).Methods(http.MethodDelete)

	// v1 Lecturaciones
	v1Lecturaciones.HandleFunc("", c.ObtenerLecturaciones).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("/{cod}", c.ObtenerLecturacion).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("", c.SubirLecturacion).Methods(http.MethodPost)
	v1Lecturaciones.HandleFunc("/{cod}", c.ModificarLecturacion).Methods(http.MethodPut)
	v1Lecturaciones.HandleFunc("/{cod}", c.EliminarLecturacion).Methods(http.MethodDelete)

	// v1 Criticas
	v1Criticas.HandleFunc("", c.ObtenerCriticas).Methods(http.MethodGet)
	v1Criticas.HandleFunc("/{cod}", c.ObtenerCritica).Methods(http.MethodGet)
	v1Criticas.HandleFunc("", c.SubirCritica).Methods(http.MethodPost)
	v1Criticas.HandleFunc("/{cod}", c.ModificarCritica).Methods(http.MethodPut)

	// v1 Rutas
	v1Rutas.HandleFunc("", c.ObtenerRutas).Methods(http.MethodGet)
	v1Rutas.HandleFunc("/{cod}", c.ObtenerRuta).Methods(http.MethodGet)
	v1Rutas.HandleFunc("", c.SubirRuta).Methods(http.MethodPost)
	v1Rutas.HandleFunc("/{cod}", c.ModificarRuta).Methods(http.MethodPut)

	//v1 Usuarios
	v1Usuarios.HandleFunc("/lecturador", c.Usuario.ObtenerLecturadores).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/lecturador", c.Usuario.ModificarLecturadorRutaGrupo).Methods(http.MethodPut)
}

func endPointsWS(ws *mux.Router) {
	v1 := ws.PathPrefix("/v1").Subrouter()
	v1Medidores := v1.PathPrefix("/medidores").Subrouter()

	// v1 Medidores
	v1Medidores.HandleFunc("/medidor", c.ObtenerMedidoresWS)
}
