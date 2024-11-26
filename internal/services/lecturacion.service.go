package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type lecturacion struct {
}

var Lecturacion lecturacion

func (m *lecturacion) GetAll(l *[]models.Lecturacion) error {
	tx := db.GDB.Begin()
	if err := tx.Find(&l).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (m *lecturacion) GetById(i *models.Lecturacion, id string) error {
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *lecturacion) Save(i *models.Lecturacion) error {
	tx := db.GDB.Begin()

	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
