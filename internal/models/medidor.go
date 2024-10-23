package models

type Medidor struct {
	COD           uint          `gorm:"primaryKey;autoIncrement" json:"cod"`
	Estado        string        `json:"estado"`
	Nombre        string        `json:"nombre"`
	Propietario   string        `json:"propietario"`
	CodRuta       *uint         `json:"codRuta"`
	CodDireccion  *uint         `json:"codDireccion"`
	Tipo          string        `json:"tipo"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:COD" json:"lecturaciones,omitempty"`
}

func (Medidor) TableName() string {
	return "medidor"
}
