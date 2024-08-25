package models

type TablerRuta interface {
	TableName() string
}

type Ruta struct {
	COD    uint   `gorm:"primaryKey;autoIncrement" json:"CodigoRuta"`
	Zona   string `json:"zona"`
	Nombre string `json:"nombre"`
	// Relación
	Estado string `json:"estado"`
}

// Implementación de la interfaz TablerRuta
func (Ruta) TableName() string {
	return "ruta"
}
