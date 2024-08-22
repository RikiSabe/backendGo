package models

import "time"

type TablerUsuario interface {
	TableName() string
}

type Usuario struct {
	COD           uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Rol           string `json:"rol"`
	NombreUsuario string `json:"nombreUsuario"`
	Contra        string `json:"contra"`
	// Claves foráneas
	CodRuta       uint          `json:"codRuta"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:CodLecturador"` // Se corrigió la sintaxis aquí
	Persona       Persona       `gorm:"foreignKey:COD"`
	CreatedAt     time.Time     `gorm:"default:now()"`
	UpdatedAt     time.Time
}

// Implementación de la interfaz TablerUsuario
func (Usuario) TableName() string {
	return "usuario"
}
