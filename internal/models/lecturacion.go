package models

import (
	"time"

	"gorm.io/datatypes"
)

type Lecturacion struct {
	COD         uint           `gorm:"primaryKey;autoIncrement"`
	NroRegistro uint           `gorm:"autoIncrement"`
	Medicion    *uint          `json:"consumo"`
	Hora        datatypes.Time `json:"hora"`
	Fecha       datatypes.Date `json:"fecha"`

	CodRuta       uint
	CodLecturador uint
	CodMedidor    uint

	Lecturador Usuario `gorm:"foreignKey:CodLecturador" json:"lecturador"`
	Medidor    Medidor `gorm:"foreignKey:CodMedidor" json:"medidor"`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Lecturacion) TableName() string {
	return "lecturacion"
}
