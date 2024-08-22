package models

type TablerRuta interface {
	TableName() string
}

type Ruta struct {
	COD    uint   `gorm:"primaryKey;autoIncrement" json:"CodigoRuta"`
	Zona   string `json:"zona"`
	Nombre string `json:"nombre"`
	// Relación
	Usuario       Usuario       `gorm:"foreignKey:CodRuta"`
	Medidores     []Medidor     `gorm:"foreignKey:CodRuta"`
	Lecturaciones []Lecturacion `gorm:"foreignKey:CodRuta"` // Se corrigió la sintaxis aquí
}

// Implementación de la interfaz TablerRuta
func (Ruta) TableName() string {
	return "ruta"
}
