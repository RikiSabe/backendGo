package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type ruta struct {
}

var Ruta ruta

func (m *ruta) GetAll(l *[]models.Ruta) error {

	if err := db.GDB.Find(&l).Error; err != nil {
		return err
	}
	return nil
}

func (m *ruta) GetById(i *models.Ruta, id string) error {
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *ruta) Save(i *models.Ruta) error {
	tx := db.GDB.Begin()
	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *ruta) Update(cod string, i *models.Ruta) error {
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", cod).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
