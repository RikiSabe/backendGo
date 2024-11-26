package models

import "time"

type Direccion struct {
	COD      uint    `gorm:"primaryKey;autoIncrement" json:"cod"`
	Longitud float64 `json:"longitud"`
	Latitud  float64 `json:"latitud"`

	Medidor Medidor `gorm:"foreignKey:CodDireccion" json:"medidor,omitempty"`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Direccion) TableName() string {
	return "direccion"
}
