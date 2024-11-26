package models

import (
	"time"
)

type Grupo struct {
	COD uint `gorm:"primaryKey;autoIncrement" json:"cod"`

	CodUsuario    *uint    `json:"cod_usuario"`
	CodRuta       *uint    `json:"cod_ruta"`
	Ruta          *Ruta    `gorm:"foreignKey:CodRuta;references:COD" json:"-"`
	NombreRuta    *string  `json:"nombre_ruta"`
	Usuario       *Usuario `gorm:"foreignKey:CodUsuario;references:COD" json:"-"`
	NombreUsuario *string  `json:"nombre_usuario"`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Grupo) TableName() string {
	return "grupo"
}
