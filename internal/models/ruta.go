package models

import "time"

type Tabler interface {
	TableName() string
}

type Ruta struct {
	COD    uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Zona   string `json:"zona"`
	Nombre string `json:"nombre"`
	Estado string `json:"estado"`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

// Implementación de la interfaz TablerRuta
func (Ruta) TableName() string {
	return "ruta"
}
