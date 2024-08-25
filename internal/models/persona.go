package models

import (
	"time"
)

type Persona struct {
	COD             uint      `gorm:"primaryKey;autoIncrement"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	FechaNacimiento string    `json:"fechaNacimiento"`
	Usuario         Usuario   `gorm:"foreignKey:CodPersona"`
	CreatedAt       time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (Persona) TableName() string {
	return "persona"
}
