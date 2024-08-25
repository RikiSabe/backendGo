package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"net/http"
)

type usuario struct {
}

var Usuario usuario

func (usuario) ObtenerLecturadores(w http.ResponseWriter, r *http.Request) {
	var lecturadores []struct {
		CodPersona      uint   `json:"codPersona"`
		Nombre          string `json:"nombre"`
		Apellido        string `json:"apellido"`
		FechaNacimiento string `json:"fechaNacimiento"`
		Usuario         string `json:"usuario"`
		CodRuta         uint   `json:"codRuta,omitempty"`
		NombreRuta      string `json:"nombreRuta,omitempty"`
	}
	query := `select p.cod as cod_persona,p.nombre, p.apellido,p.fecha_nacimiento as fecha_nacimiento,u.usuario,r.cod as cod_ruta,r.nombre as nombre_ruta from persona p
left join usuario u
on u.cod_persona = p.cod
left join ruta r
on u.cod_ruta = r.cod
where u.rol ='lecturador';`
	tx := db.GDB.Begin()
	if err := tx.Raw(query).Find(&lecturadores).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&lecturadores); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
func (usuario) ModificarLecturadorRutaGrupo(w http.ResponseWriter, r *http.Request) {
	var lecturador struct {
		COD      uint  `json:"cod"`
		CodRuta  uint  `json:"codRuta"`
		CodGrupo *uint `json:"codGrupo"`
	}
	tx := db.GDB.Begin()
	tx.Model(models.Usuario{}).Select("cod_ruta", "cod_grupo").Where("cod = ?", lecturador.COD).Updates(lecturador)
	tx.Commit()
	return
}
