package models

type TablerCritica interface {
	TableName() string
}

type Critica struct {
	COD            uint   `gorm:"primaryKey;autoIncrement" json:"cod"`
	Descripcion    string `json:"descripcion"`
	Tipo           string `json:"tipo"`
	CodLecturacion uint   `json:"codLecturacion"`
}

func (Critica) TableName() string {
	return "critica"
}
