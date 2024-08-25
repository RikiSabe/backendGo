package models

type Critica struct {
	COD            uint   `gorm:"primaryKey;AutoIncrement" json:"cod"`
	Descripcion    string `json:"descripcion"`
	Tipo           string `json:"tipo"`
	Estado         string `json:"estado"`
	CodLecturacion uint   `json:"codLecturacion"`
}

func (Critica) TableName() string {
	return "critica"
}
