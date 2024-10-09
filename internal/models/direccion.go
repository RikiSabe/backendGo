package models

type Direccion struct {
	COD      uint    `gorm:"primaryKey;autoIncrement" json:"cod"`
	Longitud float64 `json:"longitud"`
	Latitud  float64 `json:"latitud"`
	Medidor  Medidor `gorm:"foreignKey:CodDireccion" json:"medidor,omitempty"`
}

func (Direccion) TableName() string {
	return "direccion"
}
