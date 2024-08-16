package models

import "gorm.io/datatypes"

type TablerLecturacion interface {
	TableName() string
}

type Lecturacion struct {
	Codigo           uint `gorm:"primaryKey;autoIncrement"`
	CodigoLecturador uint
	Hora             datatypes.Time `json:"hora"`
	Fecha            datatypes.Date `json:"fecha"`
	Registro         uint
}

func (Lecturacion) TableName() string {
	return "lecturacion"
}
