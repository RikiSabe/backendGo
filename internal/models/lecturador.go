package models

type TablerLecturador interface {
	TableName() string
}

type Lecturador struct {
	Codigo           uint           `gorm:"primaryKey;autoIncrement"`
	Nombre           string         `json:"nombre"`
	Apellido         string         `json:"apellido"`
	Clave            string         `json:"clave"`
	Lecturaciones    []Lecturacion  `gorm:"foreignKey:CodigoLecturador"`
	ReportesErrorres []ReporteError `gorm:"foreignKey:Codigo"`
}

func (Lecturador) TableName() string {
	return "lecturador"
}
