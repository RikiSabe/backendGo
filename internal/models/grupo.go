package models

import (
	"time"
)

type Grupo struct {
	COD        uint      `gorm:"primaryKey;autoIncrement" json:"cod"`
	CodUsuario *uint     `json:"cod_usuario"`
	CodRuta    *uint     `json:"cod_ruta"`
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time
}

func (Grupo) TableName() string {
	return "grupo"
}
