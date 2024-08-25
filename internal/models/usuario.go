package models

import "time"

type Usuario struct {
	COD     uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Rol     string `json:"rol"`
	Usuario string `json:"usuario"`
	Contra  string `json:"contra"`
	// Claves foráneas
	CodRuta       uint `json:"codRuta"`
	CodPersona    uint
	CodGrupo      *uint
	Lecturaciones []Lecturacion `gorm:"foreignKey:CodLecturador"` // Se corrigió la sintaxis aquí

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

// Implementación de la interfaz TablerUsuario
func (Usuario) TableName() string {
	return "usuario"
}
