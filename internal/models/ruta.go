package models

type TablerRuta interface {
	TableName() string
}

type Ruta struct {
	Codigo    uint      `gorm:"primaryKey;autoIncrement"`
	Zona      string    `json:"zona"`
	Medidores []Medidor `gorm:"foreignKey:CodigoRuta"`
}

func (Ruta) TableName() string {
	return "ruta"
}
