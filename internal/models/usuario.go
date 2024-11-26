package models

import "time"

type Usuario struct {
	COD     uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Rol     string `json:"rol"`
	Usuario string `json:"usuario"`
	Contra  string `json:"contra"`
	Estado  string `json:"estado"`

	CodRuta       *uint         `json:"-"`
	CodPersona    uint          `json:"-"`
	Persona       *Persona      `gorm:"foreignKey:CodPersona"`
	CodGrupo      *uint         `json:"-"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:CodLecturador" json:"lecturaciones,omitempty"`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Usuario) TableName() string {
	return "usuario"
}
