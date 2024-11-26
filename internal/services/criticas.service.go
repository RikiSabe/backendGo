package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type critica struct {
}

var Critica critica

func (m *critica) CountActivos(total *int64) error {
	tx := db.GDB.Begin()

	if err := tx.Model(&models.Critica{}).Where("estado = ?", "activo").Count(total).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *critica) GetAll(l *[]models.Critica) error {
	tx := db.GDB.Begin()
	if err := tx.Order("cod asc").Find(&l).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *critica) GetById(i *models.Critica, id string) error {
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *critica) Save(i *models.Critica) error {
	tx := db.GDB.Begin()

	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *critica) Update(cod string, i *models.Critica) error {
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", cod).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
