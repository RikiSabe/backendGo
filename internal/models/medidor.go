package models

import "time"

type Medidor struct {
	COD         uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Estado      string `json:"estado"`
	Nombre      string `json:"nombre"`
	Propietario string `json:"propietario"`
	Tipo        string `json:"tipo"`

	CodRuta       *uint         `json:"codRuta"`
	CodDireccion  *uint         `json:"codDireccion"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:COD" json:"lecturaciones,omitempty"`
	Ruta          Ruta          `gorm:"foreignKey:CodRuta;references:COD" json:"-"`
	NombreRuta    string        `gorm:"-" json:"nombreRuta"`
	CreatedAt     time.Time     `gorm:"default:now()"`
	UpdatedAt     time.Time
}

func (Medidor) TableName() string {
	return "medidor"
}
