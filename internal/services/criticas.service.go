package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type critica struct {
}

var Critica critica

func (m *critica) CountActivos(total *int64) error {
	// Inicia una transacción solo para la consulta
	tx := db.GDB.Begin()

	// Cuenta los medidores activos
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
	// Buscar el medidor por ID y asegurarse de que esté activo
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *critica) Save(i *models.Critica) error {
	tx := db.GDB.Begin()

	// Intentar guardar el medidor en la base de datos
	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *critica) Update(cod string, i *models.Critica) error {
	//medidor := &i
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", cod).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
