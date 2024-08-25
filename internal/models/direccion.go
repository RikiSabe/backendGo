package models

type Direccion struct {
	COD         uint    `gorm:"primaryKey;autoIncrement" json:"cod"`
	CoordenadaX float32 `json:"coordenadaX"`
	CoordenadaY float32 `json:"coordenadaY"`
	Medidor     Medidor `gorm:"foreignKey:CodDireccion" json:"medidor,omitempty"`
}

func (Direccion) TableName() string {
	return "direccion"
}
