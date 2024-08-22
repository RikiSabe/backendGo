package models

type TablerMedidor interface {
	TableName() string
}

type Medidor struct {
	CodigoMedidor uint   `gorm:"primaryKey;AutoIncrement" json:"codigoMedidor"`
	Estado        string `json:"estado"`
	Medicion      int    `json:"medicion"`
	Nombre        string `json:"nombre"`
	Propietario   string `json:"propietario"`
	REC           string `json:"rec"`
	Registro      string `json:"registro"`
	// Foreigns Keys
	Direccion Direccion `gorm:"foreignKey:CodigoDireccion" json:"codigoDireccion,omitempty"`
}

func (Medidor) TableName() string {
	return "medidor"
}
