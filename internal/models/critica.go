package models

type TablerCritica interface {
	TableName() string
}

type Critica struct {
	CodigoCritica uint   `gorm:"primaryKey;autoIncrement" json:"codigoCritica"`
	Descripcion   string `json:"descripcion"`
	Tipo          string `json:"tipo"`
}

func (Critica) TableName() string {
	return "critica"
}
