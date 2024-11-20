package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ObtenerGrupoPorCod(w http.ResponseWriter, r *http.Request) {
	codGrupo := mux.Vars(r)["cod_grupo"]
	var grupo models.Grupo

	query := `SELECT g.cod AS cod, 
                     g.cod_usuario AS cod_usuario, 
                     g.cod_ruta AS cod_ruta, 
                     u.usuario AS nombre_usuario, 
                     r.nombre AS nombre_ruta
              FROM grupo AS g
              LEFT JOIN usuario AS u ON g.cod_usuario = u.cod
              LEFT JOIN ruta AS r ON g.cod_ruta = r.cod
              WHERE g.cod = ?;`

	if err := db.GDB.Raw(query, codGrupo).Scan(&grupo).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(grupo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerGrupos(w http.ResponseWriter, r *http.Request) {
	var grupos []models.Grupo

	query := `SELECT g.cod AS cod, 
                     g.cod_usuario AS cod_usuario, 
                     g.cod_ruta AS cod_ruta, 
                     u.usuario AS nombre_usuario, 
                     r.nombre AS nombre_ruta
              FROM grupo AS g
              LEFT JOIN usuario AS u ON g.cod_usuario = u.cod
              LEFT JOIN ruta AS r ON g.cod_ruta = r.cod;`

	// Ejecutar la consulta y escanear los resultados en `grupos`
	if err := db.GDB.Raw(query).Scan(&grupos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar la respuesta en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(grupos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SubirGrupo(w http.ResponseWriter, r *http.Request) {
	var grupo models.Grupo
	codPersona := mux.Vars(r)["cod_persona"]

	// Decodificar el JSON del cuerpo de la solicitud en el objeto grupo
	if err := json.NewDecoder(r.Body).Decode(&grupo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Asignar codPersona al campo CodUsuario en grupo
	var usuario models.Usuario
	if err := db.GDB.Where("cod_persona = ?", codPersona).First(&usuario).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	grupo.CodUsuario = &usuario.COD // Asigna el COD del usuario encontrado

	// Guardar el nuevo grupo
	if err := db.GDB.Create(&grupo).Error; err != nil {
		http.Error(w, "Ha ocurrido un error al guardar en la BD", http.StatusInternalServerError)
		return
	}

	// Actualizar el campo CodGrupo del usuario correspondiente
	usuario.CodGrupo = &grupo.COD
	if err := db.GDB.Save(&usuario).Error; err != nil {
		http.Error(w, "Ha ocurrido un error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&grupo); err != nil {
		http.Error(w, "Ha ocurrido un error al parsear a JSON", http.StatusInternalServerError)
		return
	}
}

func ModificarGrupo(w http.ResponseWriter, r *http.Request) {
	var grupoActualizado models.Grupo
	cod := mux.Vars(r)["cod"]

	// Buscar el grupo existente por su código
	var grupoExistente models.Grupo
	if err := services.Grupo.GetById(&grupoExistente, cod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON recibido en el request
	if err := json.NewDecoder(r.Body).Decode(&grupoActualizado); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	grupoExistente.CodUsuario = grupoActualizado.CodUsuario
	grupoExistente.CodRuta = grupoActualizado.CodRuta
	// Guardar los cambios en el grupo existente
	if err := db.GDB.Save(&grupoExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&grupoExistente); err != nil {
		http.Error(w, "Ha ocurrido un error al codificar a JSON", http.StatusInternalServerError)
		return
	}
}

func QuitarCodGrupo(w http.ResponseWriter, r *http.Request) {
	codUsuario := mux.Vars(r)["cod_usuario"]

	var usuario models.Usuario
	if err := db.GDB.Where("cod = ?", codUsuario).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usuario.CodGrupo = nil
	if err := db.GDB.Save(&usuario).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, "Error al codificar el usuario a JSON", http.StatusInternalServerError)
		return
	}
}

func EliminarGrupo(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]

	if err := db.GDB.Where("cod = ?", cod).Delete(&models.Grupo{}).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Grupo eliminado exitosamente")); err != nil {
		http.Error(w, "Error al enviar la respuesta", http.StatusInternalServerError)
		return
	}
}

func ObtenerDatosGenerales(w http.ResponseWriter, r *http.Request) {
	const query = `
		SELECT 
			(SELECT COUNT(m.cod) FROM medidor m) AS cantidad_medidores,
			(SELECT COUNT(u.cod) FROM usuario u WHERE u.rol = 'lecturador') AS cantidad_lecturadores,
			(SELECT COUNT(r.cod) FROM ruta r) AS cantidad_rutas,
			(SELECT COUNT(c.cod) FROM critica c) AS cantidad_criticas;`

	type CantidadesEntidades struct {
		CantidadMedidores    uint64 `json:"cantidadMedidores"`
		CantidadLecturadores uint64 `json:"cantidadLecturadores"`
		CantidadRutas        uint64 `json:"cantidadRutas"`
		CantidadCriticas     uint64 `json:"cantidadCriticas"`
	}

	var cantidades CantidadesEntidades
	// Execute query and check for errors.
	result := db.GDB.Raw(query).Scan(&cantidades)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&cantidades); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
