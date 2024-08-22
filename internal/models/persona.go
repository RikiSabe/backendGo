package models

import (
	"time"
)

type TablerPersona interface {
	TableName() string
}

type Persona struct {
	CodigoPersona   uint   `gorm:"primaryKey;autoIncrement" json:"codigoPersona"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	FechaNacimiento string `json:"fechaNacimiento"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Persona) TableName() string {
	return "persona"
}
