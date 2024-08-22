package models

import "time"

type TablerUsuario interface {
	TableName() string
}

type Usuario struct {
	CodigoUsuario uint   `gorm:"primaryKey;autoIncrement" json:"codigoUsuario"`
	Rol           string `json:"rol"`
	NombreUsuario string `json:"nombreUsuario"`
	Contra        string `json:"contra"`

	// Foreigns Keys
	RutaAsignada Ruta    `gorm:"foreignKey:CodigoRuta"`
	Persona      Persona `gorm:"foreignKey:CodigoPersona"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Usuario) TableName() string {
	return "usuario"
}
