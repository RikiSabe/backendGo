package models

import (
	"time"
)

type TablerPersona interface {
	TableName() string
}

type Persona struct {
	COD             uint      `gorm:"primaryKey;autoIncrement"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	FechaNacimiento string    `json:"fechaNacimiento"`
	CreatedAt       time.Time `gorm:"default:now()"`
	UpdatedAt       time.Time
}

func (Persona) TableName() string {
	return "persona"
}
