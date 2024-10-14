package models

type Medidor struct {
	COD           uint          `gorm:"primaryKey;autoIncrement" json:"cod"`
	Estado        string        `json:"estado"`
	Medicion      int           `json:"medicion"` // nova
	Nombre        string        `json:"nombre"`
	Propietario   string        `json:"propietario"`
	REC           string        `json:"rec"`      // no va
	Registro      string        `json:"registro"` // no va
	CodRuta       *uint         `json:"codRuta"`
	CodDireccion  *uint         `json:"codDireccion"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:COD" json:"lecturaciones,omitempty"`
}

func (Medidor) TableName() string {
	return "medidor"
}
