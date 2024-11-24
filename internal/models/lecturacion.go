package models

import (
	"time"

	"gorm.io/datatypes"
)

type Lecturacion struct {
	COD           uint `gorm:"primaryKey;autoIncrement"`
	CodRuta       uint
	CodLecturador uint
	CodMedidor    uint
	NroRegistro   uint           `gorm:"autoIncrement"`
	Medicion      *uint          `json:"consumo"`
	Hora          datatypes.Time `json:"hora"`
	Fecha         datatypes.Date `json:"fecha"`
	CreatedAt     time.Time      `gorm:"default:now()"`
	UpdatedAt     time.Time
	// Relaciones
	Medidor Medidor `gorm:"foreignKey:CodMedidor"` // Relación con Medidor
}

func (Lecturacion) TableName() string {
	return "lecturacion"
}
