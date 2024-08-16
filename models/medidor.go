package models

type TablerMedidor interface {
	TableName() string
}

type Medidor struct {
	Codigo          uint   `gorm:"primaryKey;AutoIncrement" json:"id"`
	CodigoUbicacion string `json:"codigo_ubicacion"`
	Direccion       string
	Estado          string
	Medicion        uint
	Nombre          string `json:"nombre"`
	Propietario     string
	REC             uint
	CodigoRuta      *string
	Lecturaciones   []Lecturacion `gorm:"foreignKey:Registro"`
}

func (Medidor) TableName() string {
	return "medidor"
}
