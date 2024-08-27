package models

type Direccion struct {
	COD     uint    `gorm:"primaryKey;autoIncrement" json:"cod"`
	CoordX  float32 `json:"coordX"`
	CoordY  float32 `json:"coordY"`
	Medidor Medidor `gorm:"foreignKey:CodDireccion" json:"medidor,omitempty"`
}

func (Direccion) TableName() string {
	return "direccion"
}
