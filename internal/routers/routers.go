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
	v1Grupos := v1.PathPrefix("/grupos").Subrouter()

	// v1 Personas
	v1Personas.HandleFunc("", c.ObtenerPersonas).Methods(http.MethodGet)
	v1Personas.HandleFunc("", c.SubirPersonas).Methods(http.MethodPost)

	// v1 Medidores
	v1Medidores.HandleFunc("/pdf", c.Reporte.MedidoresPDF).Methods(http.MethodGet)
	v1Medidores.HandleFunc("", c.ObtenerMedidores).Methods(http.MethodGet)
	v1Medidores.HandleFunc("/{cod}", c.ObtenerMedidor).Methods(http.MethodGet)
	// v1Medidores.HandleFunc("", c.PostMedidor).Methods(http.MethodPost)
	v1Medidores.HandleFunc("", c.AgregarMedidor).Methods(http.MethodPost)
	v1Medidores.HandleFunc("/{cod}", c.ModificarMedidor).Methods(http.MethodPut)
	v1Medidores.HandleFunc("/{cod}", c.EliminarMedidor).Methods(http.MethodDelete)
	v1Medidores.HandleFunc("/byruta/{cod_ruta}", c.ObtenerMedidoresByRuta).Methods(http.MethodGet)
	v1Medidores.HandleFunc("/direcciones/{cod_ruta}", c.ObtenerDireccionMedidores).Methods(http.MethodGet)
	v1Medidores.HandleFunc("/direccion/{cod_direccion}", c.ObtenerDireccion).Methods(http.MethodGet)

	// v1 Lecturaciones
	v1Lecturaciones.HandleFunc("", c.ObtenerLecturaciones).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("/{cod}", c.ObtenerLecturacion).Methods(http.MethodGet)
	v1Lecturaciones.HandleFunc("", c.SubirLecturacion).Methods(http.MethodPost)
	v1Lecturaciones.HandleFunc("/{cod}", c.ModificarLecturacion).Methods(http.MethodPut)
	v1Lecturaciones.HandleFunc("/{cod}", c.EliminarLecturacion).Methods(http.MethodDelete)

	// v1 Criticas
	v1Criticas.HandleFunc("/pdf", c.Reporte.CriticaPDF).Methods(http.MethodGet)
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
	v1Usuarios.HandleFunc("/lecturador/pdf", c.Reporte.LecturadoresPDF).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/lecturador", c.Usuario.ObtenerLecturadores).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/lecturador/ruta-grupo", c.Usuario.ModificarLecturadorRutaGrupo).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("/lecturador/{cod_usuario}/restablecer-clave", c.Usuario.RestablecerContra).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("/lecturador/{cod_lecturador}", c.Usuario.ModificarDatosLecturador).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("/lecturador", c.Usuario.AgregarLecturador).Methods(http.MethodPost)
	v1Usuarios.HandleFunc("/lecturadorbyuser/{usuario}", c.Usuario.ObtenerLecturadorPorUsuario).Methods(http.MethodGet)

	//v1 Login
	v1.HandleFunc("/login", c.Auth.AuthLogin).Methods(http.MethodPost)
	v1.HandleFunc("/loginweb", c.Auth.AuthLoginWeb).Methods(http.MethodPost)

	// Grupos
	v1Grupos.HandleFunc("", c.ObtenerGrupos).Methods(http.MethodGet)
	v1Grupos.HandleFunc("", c.SubirGrupo).Methods(http.MethodPost)

}

func endPointsWS(ws *mux.Router) {
	v1 := ws.PathPrefix("/v1").Subrouter()
	v1UbicacionLecturador := v1.PathPrefix("/ubicacion-lecturador").Subrouter()

	// v1 Medidores mobile
	v1UbicacionLecturador.HandleFunc("", c.Monitoreo.ObtenerUbicacionLecturadorWS)
	//Web
	v1UbicacionLecturador.HandleFunc("/all", c.Monitoreo.ObtenerUbicacionesLecturadorWS)
}
