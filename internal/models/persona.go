package models

import (
	"time"
)

type TablerPersona interface {
	TableName() string
}

type Persona struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Nombre    string `json:"nombre"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Persona) TableName() string {
	return "persona"
}
