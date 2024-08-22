package models

type TablerMedidor interface {
	TableName() string
}

type Medidor struct {
	COD           uint   `gorm:"primaryKey;AutoIncrement" json:"codigoMedidor"`
	Estado        string `json:"estado"`
	Medicion      int    `json:"medicion"`
	Nombre        string `json:"nombre"`
	Propietario   string `json:"propietario"`
	REC           string `json:"rec"`
	Registro      string `json:"registro"`
	CodRuta       uint
	CodDireccion  uint
	Lecturaciones []Lecturacion `gorm:"foreignKey:COD"`
}

func (Medidor) TableName() string {
	return "medidor"
}
