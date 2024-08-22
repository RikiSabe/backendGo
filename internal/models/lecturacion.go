package models

import (
	"time"
)

type TablerLecturacion interface {
	TableName() string
}

type Lecturacion struct {
	CodigoLecturacion uint      `gorm:"primaryKey;autoIncrement" json:"codigoLecturacion"`
	CodigoLecturador  Usuario   `gorm:"foreignKey:CodigoUsuario"`
	Criticas          []Critica `gorm:"foreignKey:CodigoCritica"`
	CodigoRuta        Ruta      `gorm:"foreignKey:CodigoRuta"`
	Hora              time.Time
	Fecha             time.Time
	NumeroRegistro    int `gorm:"autoIncrement"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Lecturacion) TableName() string {
	return "lecturacion"
}
