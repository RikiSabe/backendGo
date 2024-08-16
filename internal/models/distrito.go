package models

type TablerDistrito interface {
	TableName() string
}

type Distrito struct {
	Codigo uint `gorm:"primaryKey;autoIncrement"`
}

func (Distrito) TableName() string {
	return "distrito"
}
