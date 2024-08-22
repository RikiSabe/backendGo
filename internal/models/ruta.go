package models

type TablerRuta interface {
	TableName() string
}

type Ruta struct {
	CodigoRuta uint   `gorm:"primaryKey;autoIncrement" json:"CodigoRuta"`
	Zona       string `json:"zona"`
	Nombre     string `json:"nombre"`
	// Relation
	Medidores []Medidor `gorm:"foreignKey:CodigoMedidor"`
}

func (Ruta) TableName() string {
	return "ruta"
}
