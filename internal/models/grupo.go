package models

import (
	"time"
)

type TablerGrupo interface {
	TableName() string
}

type Grupo struct {
	CodigoGrupo uint      `gorm:"primaryKey;autoIncrement" json:"codigoGrupo"`
	NumeroGrupo int       `json:"numeroGrupo"`
	Personas    []Persona `gorm:"foreignKey:CodigoPersona"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Grupo) TableName() string {
	return "grupo"
}
