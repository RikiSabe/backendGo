package models

import (
	"time"
)

type Persona struct {
	COD       uint      `gorm:"primaryKey;autoIncrement" json:"cod"`
	Nombre    string    `json:"nombre"`
	Apellido  string    `json:"apellido"`
	CI        string    `json:"ci"`
	Usuario   Usuario   `gorm:"foreignKey:CodPersona" json:"-"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Persona) TableName() string {
	return "persona"
}
