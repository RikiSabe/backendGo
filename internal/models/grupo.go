package models

import (
	"time"
)

type TablerGrupo interface {
	TableName() string
}

type Grupo struct {
	COD       uint      `gorm:"primaryKey;autoIncrement" json:"cod"`
	Nro       int       `json:"numeroGrupo"`
	Personas  []Persona `gorm:"foreignKey:COD"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Grupo) TableName() string {
	return "grupo"
}
