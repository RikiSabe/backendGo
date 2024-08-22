package models

import "time"

type TablerUsuario interface {
	TableName() string
}

type Usuario struct {
	COD           uint   `gorm:"primaryKey;autoIncrement" json:"codigoUsuario"`
	Rol           string `json:"rol"`
	NombreUsuario string `json:"nombreUsuario"`
	Contra        string `json:"contra"`

	// Foreigns Keys
	CodRuta uint
	// Persona   Persona `gorm:"foreignKey:CodigoPersona"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Usuario) TableName() string {
	return "usuario"
}
