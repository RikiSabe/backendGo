package models

import "time"

type TablerReporteError interface {
	TableName() string
}

type ReporteError struct {
	Codigo      uint   `gorm:"primaryKey;autoIncrement"`
	Descripcion string `json:"descripcion"`
	CreatedAt   time.Time
}

func (ReporteError) TableName() string {
	return "reporte_error"
}
