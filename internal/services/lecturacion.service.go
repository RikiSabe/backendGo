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
	// Buscar los medidores asegurarse de que esté activo
	if err := tx.Find(&l).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (m *lecturacion) GetById(i *models.Lecturacion, id string) error {
	tx := db.GDB.Begin()
	// Buscar el medidor por ID y asegurarse de que esté activo
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *lecturacion) Save(i *models.Lecturacion) error {
	tx := db.GDB.Begin()

	// Intentar guardar el medidor en la base de datos
	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *lecturacion) Delete(id string) error {
	var lecturacion models.Lecturacion
	tx := db.GDB.Begin()

	// Marcar el medidor como inactivo en lugar de eliminarlo físicamente
	if err := tx.Where("cod = ?", id).Delete(&lecturacion).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
