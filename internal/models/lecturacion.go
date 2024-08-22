package models

import (
	"time"

	"gorm.io/datatypes"
)

type TablerLecturacion interface {
	TableName() string
}

type Lecturacion struct {
	COD           uint `gorm:"primaryKey;autoIncrement"`
	CodRuta       uint
	CodLecturador uint
	NroRegistro   uint `gorm:"autoIncrement"`

	Hora      datatypes.Time `json:"hora"`
	Fecha     datatypes.Date `json:"fecha"`
	CreatedAt time.Time      `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Lecturacion) TableName() string {
	return "lecturacion"
}
