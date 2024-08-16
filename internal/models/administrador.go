package models

type TablerAdministrador interface {
	TableName() string
}

type Administrador struct {
	Codigo   uint   `gorm:"primaryKey;autoIncrement"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Clave    string `json:"clave"`
}

func (Administrador) TableName() string {
	return "administrador"
}
