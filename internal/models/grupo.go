package models

import (
	"time"
)

type Grupo struct {
	COD       uint      `gorm:"primaryKey;autoIncrement" json:"cod"`
	Nro       int       `json:"numeroGrupo"`
	Usuarios  []Usuario `gorm:"foreignKey:CodGrupo"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time
}

func (Grupo) TableName() string {
	return "grupo"
}
