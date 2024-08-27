package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		CodUsuario      uint   `json:"codUsuario"`
		Usuario         string `json:"usuario"`
		CodRuta         uint   `json:"codRuta,omitempty"`
		NombreRuta      string `json:"nombreRuta,omitempty"`
	}
	query := `select p.cod as cod_persona,p.nombre, p.apellido,p.fecha_nacimiento as fecha_nacimiento,u.cod as cod_usuario, u.usuario,r.cod as cod_ruta,r.nombre as nombre_ruta from persona p
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
}
func (usuario) RestablecerContra(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod_usuario"]
	var lecturador models.Usuario

	tx := db.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Buscar usuario por código
	if err := tx.Where("cod = ?", cod).First(&lecturador).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al buscar usuario", http.StatusInternalServerError)
		return
	}

	// Generar nueva contraseña
	newPassword, err := password.Generate(6, 2, 0, false, false)
	if err != nil {
		http.Error(w, "Error al generar la contraseña", http.StatusInternalServerError)
		return
	}

	// Cifrar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al cifrar la contraseña", http.StatusInternalServerError)
		return
	}
	lecturador.Contra = string(hashedPassword)

	// Guardar el cambio en la base de datos
	if err := tx.Save(&lecturador).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar la nueva contraseña", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// Enviar la nueva contraseña al cliente (Por ejemplo, enviarla al correo electrónico)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(newPassword)); err != nil {
		http.Error(w, "Error al enviar la contraseña", http.StatusInternalServerError)
		return
	}
}
func (usuario) ModificarDatosLecturador(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod_lecturador"]
	var personaLecturador struct {
		COD             uint   `json:"cod"`
		Nombre          string `json:"nombre"`
		Apellido        string `json:"apellido"`
		FechaNacimiento string `json:"fechaNacimiento"`
	}
	if err := json.NewDecoder(r.Body).Decode(&personaLecturador); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tx := db.GDB.Begin()
	if err := tx.Model(models.Persona{}).Where("cod = ?", cod).Updates(&personaLecturador).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}
func (usuario) AgregarLecturador(w http.ResponseWriter, r *http.Request) {
	var lecturador struct {
		Usuario         string `json:"usuario"`
		Nombre          string `json:"nombre"`
		Apellido        string `json:"apellido"`
		FechaNacimiento string `json:"fechaNacimiento"`
		CodRuta         uint   `json:"codRuta"`
		CodGrupo        uint   `json:"codGrupo"`
	}

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&lecturador); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Iniciar la transacción
	tx := db.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verificar la existencia de la ruta
	if err := tx.Model(models.Ruta{}).Where("cod = ?", lecturador.CodRuta).First(&models.Ruta{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Ruta no encontrada", http.StatusNotFound)
			return
		}
		tx.Rollback()
		http.Error(w, "Error al verificar la ruta", http.StatusInternalServerError)
		return
	}

	// Verificar la existencia del grupo
	if err := tx.Model(models.Grupo{}).Where("cod = ?", lecturador.CodGrupo).First(&models.Grupo{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Grupo no encontrado", http.StatusNotFound)
			return
		}
		tx.Rollback()
		http.Error(w, "Error al verificar el grupo", http.StatusInternalServerError)
		return
	}

	// Crear el objeto Usuario
	lecturadorR := models.Usuario{
		Usuario:  lecturador.Usuario,
		Rol:      "lecturador",
		CodGrupo: &lecturador.CodGrupo,
		CodRuta:  lecturador.CodRuta,
		Persona: &models.Persona{
			Nombre:          lecturador.Nombre,
			Apellido:        lecturador.Apellido,
			FechaNacimiento: lecturador.FechaNacimiento,
		},
	}

	// Crear la persona en la base de datos
	if err := tx.Create(*&lecturadorR.Persona).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar la persona", http.StatusInternalServerError)
		return
	}

	// Generar nueva contraseña
	newPassword, err := password.Generate(6, 2, 0, false, false)
	if err != nil {
		http.Error(w, "Error al generar la contraseña", http.StatusInternalServerError)
		return
	}

	// Cifrar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al cifrar la contraseña", http.StatusInternalServerError)
		return
	}
	lecturadorR.Contra = string(hashedPassword)

	// Asignar el código de persona al usuario y crear el usuario en la base de datos
	lecturadorR.CodPersona = lecturadorR.Persona.COD
	if err := tx.Create(&lecturadorR).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
		return
	}

	// Confirmar la transacción
	tx.Commit()
	w.Write([]byte(newPassword))
}
