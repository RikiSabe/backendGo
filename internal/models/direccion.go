package models

type TablerDireccion interface {
	TableName() string
}

type Direccion struct {
	CodigoDireccion uint    `gorm:"primaryKey;AutoIncrement" json:"codigoDireccion"`
	CoordenadaX     float32 `json:"coordenadaX"`
	CoordenadaY     float32 `json:"coordenadaY"`
}

func (Direccion) TableName() string {
	return "direccion"
}
